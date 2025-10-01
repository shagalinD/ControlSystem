export interface ApiResponse<T> {
  success: boolean
  message: string
  data: T
}

export interface PaginationMeta {
  page: number
  page_size: number
  total: number
  total_pages: number
}

// Пользователь
export interface User {
  id: number
  email: string
  full_name: string
  role_id: number
  role_name: 'engineer' | 'manager' | 'observer'
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
  user: User
}

// Отчеты
export interface DefectsReport {
  total_defects: number
  defects_by_status: Array<{ status: DefectStatus; count: number }>
  defects_by_priority: Array<{ priority: DefectPriority; count: number }>
  overdue_defects: number
  avg_resolution_time: number
}

// Фильтры
export interface DefectFilters {
  project_id?: number
  status?: DefectStatus
  priority?: DefectPriority
  assignee_id?: number
  page?: number
  page_size?: number
  sort_by?: string
  order?: 'asc' | 'desc'
  search?: string
}

// Формы данных
export interface LoginData {
  email: string
  password: string
}

export interface RegisterData {
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
  assignee_id?: number
  deadline?: string
}

export interface CreateProjectData {
  name: string
  description: string
  manager_id: number
}

export interface CreateCommentData {
  text: string
  defect_id: number
}

// Enums
export type DefectStatus =
  | 'new'
  | 'in_progress'
  | 'on_review'
  | 'closed'
  | 'cancelled'
export type DefectPriority = 'low' | 'medium' | 'high' | 'critical'
export type UserRole = 'engineer' | 'manager' | 'observer'
