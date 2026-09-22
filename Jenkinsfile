// ==============================================================
// Jenkinsfile — CI/CD Pipeline for Maintenance Request Log
// ==============================================================

pipeline {
    agent any

    environment {
        DOCKER_BUILDKIT   = '1'
        APP_ENV           = 'test'
        POSTGRES_DB       = 'maintenance_test_db'
        POSTGRES_USER     = 'test_user'
        POSTGRES_PASSWORD = 'test_password'

        BACKEND_TEST_IMAGE = 'maintenance-backend-test'
    }

    options {
        timeout(time: 30, unit: 'MINUTES')
        disableConcurrentBuilds()
        ansiColor('xterm')
    }

    stages {

        // =====================================================
        // Stage 0: Checkout Source Code
        // =====================================================
        stage('Checkout') {
            steps {
                echo 'Checking out source repository...'
                script {
                    // Cek apakah workspace sudah memiliki source code (misal dari SCM atau sync lokal)
                    if (!fileExists('backend/Dockerfile')) {
                        echo 'Workspace empty, cloning from Git repository...'
                        git branch: 'master', url: 'https://github.com/Grandvill/maintenance-request-log.git'
                    } else {
                        echo 'Source code already present in workspace.'
                    }
                }
            }
        }

        // =====================================================
        // Stage 1: Backend Docker Build & Unit Tests
        // =====================================================
        stage('Backend Tests') {
            steps {
                echo 'Building Go backend builder image...'
                sh '''
                    docker build \
                        --target builder \
                        -t ${BACKEND_TEST_IMAGE} \
                        ./backend
                '''

                echo 'Running Go static checks inside Docker...'
                sh '''
                    docker run --rm \
                        ${BACKEND_TEST_IMAGE} \
                        go vet ./...
                '''

                echo 'Running Go unit tests inside Docker...'
                sh '''
                    docker run --rm \
                        ${BACKEND_TEST_IMAGE} \
                        go test -v -race ./...
                '''
            }
            post {
                always {
                    echo 'Backend test stage completed.'
                }
            }
        }

        // =====================================================
        // Stage 2: Frontend Build Check (Isolated inside Docker)
        // =====================================================
        stage('Frontend Build & Check') {
            steps {
                echo 'Testing frontend build inside Node.js container...'
                // Menjalankan build di container node:24-alpine agar tidak memerlukan npm terpasang di host Jenkins
                sh '''
                    docker run --rm \
                        -v "$(pwd)/frontend:/app" \
                        -w /app \
                        node:24-alpine \
                        sh -c "npm install --no-audit --no-fund && npm run build"
                '''
            }
            post {
                always {
                    echo 'Frontend build stage completed.'
                }
            }
        }

        // =====================================================
        // Stage 3: Docker Compose Build Verification
        // =====================================================
        stage('Docker Images Build') {
            steps {
                echo 'Preparing environment file...'
                sh '''
                    if [ ! -f .env ] && [ -f .env.example ]; then
                        cp .env.example .env
                    fi
                '''

                echo 'Building all Docker Compose images...'
                sh 'docker compose build'
            }
            post {
                always {
                    echo 'Docker images build stage completed.'
                }
            }
        }

        // =====================================================
        // Stage 4: Container Smoke Test
        // =====================================================
        stage('Smoke Test') {
            steps {
                echo 'Starting db and backend test containers...'
                sh 'docker compose up -d db backend'

                echo 'Waiting for backend healthcheck to respond OK...'
                sh '''
                    for i in $(seq 1 15); do
                        # Menggunakan curl di dalam container network agar dapat menjangkau container backend
                        if docker run --rm --network maintenance-request-log_maintenance_network curlimages/curl:latest -fs http://backend:8080/health; then
                            echo "Backend healthcheck passed!"
                            exit 0
                        fi
                        echo "Waiting for service to be healthy... ($i/15)"
                        sleep 2
                    done
                    echo "Healthcheck timed out!"
                    exit 1
                '''
            }
            post {
                always {
                    echo 'Cleaning up test containers and volumes...'
                    sh 'docker compose down -v --remove-orphans'
                }
            }
        }
    }

    // =========================================================
    // Pipeline Post Actions
    // =========================================================
    post {
        success {
            echo '==================================================='
            echo ' Pipeline Succeeded: All tests and builds passed!   '
            echo '==================================================='
        }

        failure {
            echo '==================================================='
            echo ' Pipeline Failed: Please check the build logs above.'
            echo '==================================================='
        }

        always {
            echo 'Cleaning up temporary test images...'
            sh 'docker image rm ${BACKEND_TEST_IMAGE} || true'
            echo 'Pipeline execution completed.'
        }
    }
}