export type Role = 'operator' | 'supervisor' | 'admin'

export type Priority = 'Low' | 'Medium' | 'High' | 'Urgent'

export type RequestStatus = 'Submitted' | 'Approved' | 'Rejected'

export interface User {
  id: string
  username: string
  full_name: string
  role: Role
  is_active: boolean
  created_at: string
}

export interface MaintenanceRequest {
  id: string
  asset_id: string
  problem_description: string
  priority: Priority
  status: RequestStatus
  created_by: string
  creator_name?: string
  reviewed_by?: string | null
  reviewer_name?: string | null
  reviewed_at?: string | null
  review_note?: string | null
  created_at: string
  updated_at: string
}

export interface StatusLog {
  id: string
  request_id: string
  changed_by: string
  changer_name: string
  from_status: RequestStatus | null
  to_status: RequestStatus
  note: string | null
  created_at: string
}

export interface RequestDetail extends MaintenanceRequest {
  audit_logs: StatusLog[]
}

export interface AuthResponse {
  token: string
  expires_at: string
  user: User
}

export interface CreateRequestPayload {
  asset_id: string
  problem_description: string
  priority: Priority
}

export interface UpdateRequestPayload {
  asset_id?: string
  problem_description?: string
  priority?: Priority
  status?: RequestStatus
}

export interface ReviewRequestPayload {
  status: 'Approved' | 'Rejected'
  note?: string
}

export interface CreateUserPayload {
  username: string
  password: string
  full_name: string
  role: Role
}

export interface UpdateUserPayload {
  full_name: string
  role: Role
  password?: string
}

