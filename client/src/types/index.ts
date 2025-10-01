// src/types/index.ts

// Базовые типы
export type UserRole = 'engineer' | 'manager' | 'observer'
export type DefectStatus =
  | 'new'
  | 'in_progress'
  | 'on_review'
  | 'closed'
  | 'cancelled'
export type DefectPriority = 'low' | 'medium' | 'high' | 'critical'

// Пользователь
export interface User {
  id: number
  email: string
  full_name: string
  role_id: number
  role_name: UserRole
  created_at?: string
  updated_at?: string
}

// Проект
export interface Project {
  id: number
  name: string
  description: string
  manager_id: number
  created_at: string
  updated_at: string
  manager: User
  defects?: Defect[]
  defects_count?: number
}

// Дефект
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

// Комментарий
export interface Comment {
  id: number
  text: string
  defect_id: number
  author_id: number
  created_at: string
  author: User
}

// Вложение
export interface Attachment {
  id: number
  filename: string
  file_size: number
  mime_type: string
  uploaded_by: number
  created_at: string
  uploader: User
  download_url?: string
}

// История изменений дефекта
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

// Пагинация
export interface Pagination {
  page: number
  page_size: number
  total: number
  total_pages: number
}

// API Responses
export interface ApiResponse<T = any> {
  success: boolean
  message: string
  data: T
  error?: string
}

export interface LoginResponse {
  token: string
  user: User
}

export interface DefectsResponse {
  defects: Defect[]
  pagination: Pagination
}

export interface ProjectsResponse {
  projects: Project[]
  page: number
  pageSize: number
}

// Формы и DTO
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

// Фильтры
export interface DefectFilters {
  project_id?: number
  status?: DefectStatus
  priority?: DefectPriority
  assignee_id?: number
  author_id?: number
  page?: number
  page_size?: number
  sort_by?: string
  order?: 'asc' | 'desc'
  search?: string
}
