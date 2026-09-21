# Factory Maintenance Request Log Portal
**PT. HIROSE ELECTRIC INDONESIA — Technical Take-Home Test**

An internal web application built for factory floor operations: Operators raise maintenance requests when machinery encounters issues, Supervisors review and approve/reject requests, and Admins manage personnel and records.

---

## Quick Start (Zero-Friction Setup)

The entire application (PostgreSQL + Go Backend + Nuxt 3 Frontend) runs with a single command from a clean clone:

```bash
# 1. Copy the example environment configuration
cp .env.example .env

# 2. Build and start all services via Docker Compose
docker compose up --build
```

- **Frontend Portal**: [http://localhost:3000](http://localhost:3000)
- **Backend API**: [http://localhost:8080/api/v1](http://localhost:8080/api/v1)
- **Health Check Endpoint**: [http://localhost:8080/health](http://localhost:8080/health)

---

## Seed Login Credentials

The database automatically migrates schema and seeds default users on initial startup. You can log in immediately using any of the following accounts:

| Role | Username | Password | Permissions Summary |
| :--- | :--- | :--- | :--- |
| **Admin** | `admin` | `admin123` | Full access: CRUD all requests, review, delete, user management |
| **Supervisor** | `supervisor` | `supervisor123` | View all requests, Approve / Reject requests, raise requests |
| **Operator 1** | `operator1` | `operator123` | Raise requests, view & edit own requests while `Submitted` |
| **Operator 2** | `operator2` | `operator123` | Raise requests, view & edit own requests while `Submitted` |

*(Tip: The login page includes one-click "Quick Fill" helper buttons for rapid evaluation).*

---

## Roles & Permission Matrix (Server-Enforced)

All permission rules are **strictly enforced on the Go backend server** and cannot be bypassed via direct API calls:

| Action | Operator | Supervisor | Admin | Enforcement Mechanism |
| :--- | :---: | :---: | :---: | :--- |
| **Create request** | ✅ | ✅ | ✅ | Open to all authenticated sessions |
| **View own requests** | ✅ | ✅ | ✅ | Filtered via SQL parameterized query |
| **View all requests** | ❌ | ✅ | ✅ | If `role == 'operator'`, SQL appends `AND created_by = $user_id` |
| **Edit own request (status = Submitted)** | ✅ | ✅ | ✅ | Verified by `Service.UpdateRequest`; rejects if status != `Submitted` |
| **Edit any request** | ❌ | ❌ | ✅ | Rejects non-admin with `403 Forbidden` |
| **Approve or reject request** | ❌ | ✅ | ✅ | Guarded by `RequireRoles("supervisor", "admin")` |
| **Delete a request** | ❌ | ❌ | ✅ | Guarded by `RequireRoles("admin")` |
| **Create / Edit / Deactivate users** | ❌ | ❌ | ✅ | Guarded by `RequireRoles("admin")` |

---

## Architecture & Key Technical Decisions

### Backend: Go 1.22 + Chi Router
- **Why Go?** Fast startup time, minimal resource footprint, single static binary, and strong concurrency model.
- **Why Chi Router (`go-chi/chi/v5`)?** Minimalist and 100% compatible with standard `net/http`. Allows idiomatic subrouting and clean middleware chaining for role guards.
- **Layered Design**:
  - `cmd/server/main.go`: Entry point, config loader, auto-migration & seeder runner, graceful shutdown.
  - `internal/middleware`: JWT authentication, RBAC role guard, CORS, structured latency logger.
  - `internal/auth`: Session verification, bcrypt hashing, JWT issuance.
  - `internal/requests`: Core maintenance requests, parameterized query filtering, and audit log transitions.
  - `internal/users`: Personnel management and account status toggles.
  - `internal/database`: PostgreSQL connection pool with startup retry loop (handles container readiness).

### Frontend: Nuxt 3 + Vue 3 + Tailwind CSS
- **Single Page Application (SPA) Mode (`ssr: false`)**: Chosen for internal operational portals to guarantee consistent client-side JWT handling without SSR hydration mismatches.
- **Tailwind CSS**: Rapid, clean styling suited for industrial dashboard environments.
- **Composables**:
  - `useAuth()`: Centralized reactive auth state, role helper computed properties, and session persistence.
  - `useApi()`: `$fetch` interceptor injecting Bearer tokens and automatically handling 401 unauthenticated redirects.
  - `useRequests()`: Encapsulated request CRUD methods.

### Authentication & Token Justification
- The backend issues signed **HMAC-SHA256 JWT tokens** containing `user_id`, `username`, and `role`.
- Tokens are transmitted via `Authorization: Bearer <token>` and also stored in an **`HttpOnly` cookie** upon login.
- **Justification**: Bearer headers allow stateless decoupling for headless API testing (e.g. Postman, curl, automated test suites), while `HttpOnly` cookies provide defense-in-depth against client-side Cross-Site Scripting (XSS).

---

## Jenkinsfile CI Pipeline Stages

The repository includes a declarative `Jenkinsfile` at the root designed for continuous integration:

1. **Stage 1: Checkout**: Pulls code from source control.
2. **Stage 2: Backend Tests**: Runs `go vet ./...` (static analysis) and `go test -v -race` (unit tests and RBAC verification).
3. **Stage 3: Frontend Build & Check**: Executes `npm ci` and `npm run build` to verify type safety and bundle generation.
4. **Stage 4: Docker Images Build**: Verifies that `docker compose build` succeeds without layer caching errors.
5. **Stage 5: Container Smoke Test**: Boots up the containers in detached mode, tests `http://localhost:8080/health`, and tears down containers cleanly in the `always` post-block.

---

## Optional Tasks Completed (Bonus)

1. **Audit Trail**: Every request status transition (creation, supervisor approval, admin override) is immutably logged in `request_status_logs` with the actor's ID, previous status, new status, note, and timestamp. Displayed as a visual timeline on the request detail page.
2. **Multi-stage Dockerfiles**:
   - Backend: Builder (`golang:1.22-alpine`) $\rightarrow$ Runtime (`alpine:3.19`) with a non-root user. Image size reduced from ~850MB to **~24MB**.
   - Frontend: 3-stage Node build generating a standalone Nitro server output.
3. **Health Check Endpoint**: Available at `GET /health`, reporting database connectivity and server uptime.
4. **Automated RBAC Tests**: Located at `backend/internal/requests/rbac_test.go`, verifying that permission boundaries hold on the server side.

---

## Known Limitations

- Real-time updates: Currently updates require page reload or refetch; WebSockets or Server-Sent Events (SSE) could be introduced for instant notification when a supervisor approves a request.
- File attachments: Machinery inspection photos are currently referenced via text description rather than binary object uploads (e.g., MinIO/S3).

---

## AI Disclosure (Mandatory Section)

- **Tools Used**: Visual Studio Code With AntiGravity
- **Assisted Areas**:
  - Drafting initial boilerplate templates for `docker-compose.yml`, multi-stage `Dockerfile`s, and declarative `Jenkinsfile` syntax.
  - Generating utility Tailwind CSS layout classes for table components and timeline markers.
- **Authored by Hand / Carefully Verified**:
  - Core RBAC authorization logic: Ensuring all database queries enforce `created_by` filtering on the server when accessed by an Operator.
  - State machine transition rules: Enforcing that only requests with `Submitted` status can be edited by Operators.
  - Connection retry logic in `internal/database/database.go` to handle Docker database startup race conditions.
- **Case of Rejection / Rewrite**:
  - *Scenario*: The AI assistant initially suggested checking roles only on the frontend Nuxt middleware and passing a `can_view_all=true` query parameter from the client.
  - *Resolution*: This was **firmly rejected**. In accordance with the test requirements (*"all rules must be enforced on the server, not only hidden in the UI"*), the client parameter was removed, and the Go service was rewritten to extract the user's role directly from the verified JWT context and append database query constraints server-side.

