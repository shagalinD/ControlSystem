// Базовые типы
export interface User {
  id: number
  email: string
  full_name: string
  role_id: number
  role_name: UserRole
  created_at?: string
  updated_at?: string
}

export interface Project {
  id: number
  name: string
  description: string
  manager_id: number
  created_at: string
  updated_at: string
  manager: User
  defects?: Defect[]
}

export interface Defect {
  id: number
  title: string
  description: string
  status: DefectStatus
  priority: DefectPriority
  deadline?: string
  project_id: number
  author_id: number
  assignee_id?: number
  created_at: string
  updated_at?: string
  project: Project
  author: User
  assignee?: User
  comments?: Comment[]
  attachments?: Attachment[]
  history?: DefectHistory[]
}

export interface Comment {
  id: number
  text: string
  defect_id: number
  author_id: number
  created_at: string
  author: User
}

export interface Attachment {
  id: number
  filename: string
  file_size: number
  mime_type: string
  uploaded_by: number
  created_at?: string
  uploader: User
}

export interface DefectHistory {
  id: number
  defect_id: number
  field_name: string
  old_value: string
  new_value: string
  changed_by: number
  created_at: string
  user?: User
}

// Типы для форм
export interface LoginFormData {
  email: string
  password: string
}

export interface RegisterFormData {
  email: string
  password: string
  full_name: string
  role_id: number
}

export interface CreateDefectData {
  title: string
  description: string
  priority: DefectPriority
  deadline?: string
  project_id: number
  assignee_id?: number
}

export interface UpdateDefectData {
  title?: string
  description?: string
  status?: DefectStatus
  priority?: DefectPriority
  deadline?: string
  assignee_id?: number
}

// Типы для фильтров и пагинации
export interface PaginationParams {
  page?: number
  page_size?: number
}

export interface DefectFilters extends PaginationParams {
  project_id?: number
  status?: DefectStatus
  priority?: DefectPriority
  assignee_id?: number
  sort_by?: string
  order?: 'asc' | 'desc'
}

export interface ProjectFilters extends PaginationParams {
  manager_id?: number
}

// Response types
export interface ApiResponse<T = unknown> {
  success: boolean
  message: string
  data: T
}

export interface AuthResponse {
  token: string
  user: User
}

export interface PaginatedResponse<T = unknown> {
  data: T[]
  pagination: {
    page: number
    page_size: number
    total: number
    total_pages: number
  }
}

// Union types
export type UserRole = 'engineer' | 'manager' | 'observer'
export type DefectStatus =
  | 'new'
  | 'in_progress'
  | 'on_review'
  | 'closed'
  | 'cancelled'
export type DefectPriority = 'low' | 'medium' | 'high' | 'critical'

// Store types
export interface AuthState {
  user: User | null
  token: string | null
  isAuthenticated: boolean
  isLoading: boolean
}
