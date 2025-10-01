// src/constants/index.ts

// Роли пользователей
export const USER_ROLES = {
  ENGINEER: 'engineer',
  MANAGER: 'manager',
  OBSERVER: 'observer',
} as const

// ID ролей (соответствуют бэкенду)
export const ROLE_IDS = {
  ENGINEER: 1,
  MANAGER: 2,
  OBSERVER: 3,
} as const

// Статусы дефектов
export const DEFECT_STATUSES = {
  NEW: 'new',
  IN_PROGRESS: 'in_progress',
  ON_REVIEW: 'on_review',
  CLOSED: 'closed',
  CANCELLED: 'cancelled',
} as const

// Приоритеты дефектов
export const DEFECT_PRIORITIES = {
  LOW: 'low',
  MEDIUM: 'medium',
  HIGH: 'high',
  CRITICAL: 'critical',
} as const

// Переводы для интерфейса
// export const STATUS_LABELS: Record<DefectStatus, string> = {
//   [DEFECT_STATUSES.NEW]: 'Новая',
//   [DEFECT_STATUSES.IN_PROGRESS]: 'В работе',
//   [DEFECT_STATUSES.ON_REVIEW]: 'На проверке',
//   [DEFECT_STATUSES.CLOSED]: 'Закрыта',
//   [DEFECT_STATUSES.CANCELLED]: 'Отменена',
// }

// export const PRIORITY_LABELS: Record<DefectPriority, string> = {
//   [DEFECT_PRIORITIES.LOW]: 'Низкий',
//   [DEFECT_PRIORITIES.MEDIUM]: 'Средний',
//   [DEFECT_PRIORITIES.HIGH]: 'Высокий',
//   [DEFECT_PRIORITIES.CRITICAL]: 'Критический',
// }

// export const ROLE_LABELS: Record<UserRole, string> = {
//   [USER_ROLES.ENGINEER]: 'Инженер',
//   [USER_ROLES.MANAGER]: 'Менеджер',
//   [USER_ROLES.OBSERVER]: 'Наблюдатель',
// }

// // Цвета для статусов и приоритетов
// export const STATUS_COLORS: Record<DefectStatus, string> = {
//   [DEFECT_STATUSES.NEW]: 'bg-blue-100 text-blue-800',
//   [DEFECT_STATUSES.IN_PROGRESS]: 'bg-yellow-100 text-yellow-800',
//   [DEFECT_STATUSES.ON_REVIEW]: 'bg-purple-100 text-purple-800',
//   [DEFECT_STATUSES.CLOSED]: 'bg-green-100 text-green-800',
//   [DEFECT_STATUSES.CANCELLED]: 'bg-red-100 text-red-800',
// }

// export const PRIORITY_COLORS: Record<DefectPriority, string> = {
//   [DEFECT_PRIORITIES.LOW]: 'bg-gray-100 text-gray-800',
//   [DEFECT_PRIORITIES.MEDIUM]: 'bg-green-100 text-green-800',
//   [DEFECT_PRIORITIES.HIGH]: 'bg-orange-100 text-orange-800',
//   [DEFECT_PRIORITIES.CRITICAL]: 'bg-red-100 text-red-800',
// }

// Настройки приложения
export const APP_CONFIG = {
  API_BASE_URL: import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080',
  APP_NAME: 'Система контроля дефектов',
  DEFAULT_PAGE_SIZE: 20,
  MAX_FILE_SIZE: 10 * 1024 * 1024, // 10MB
  ALLOWED_FILE_TYPES: [
    'image/jpeg',
    'image/png',
    'image/gif',
    'application/pdf',
    'text/plain',
  ],
} as const

// Ключи localStorage
export const STORAGE_KEYS = {
  AUTH_TOKEN: 'auth_token',
  USER_DATA: 'user_data',
  THEME: 'theme',
} as const

// Маршруты
export const ROUTES = {
  HOME: '/',
  LOGIN: '/login',
  REGISTER: '/register',
  PROJECTS: '/projects',
  DEFECTS: '/defects',
  MY_DEFECTS: '/my-defects',
  REPORTS: '/reports',
  PROFILE: '/profile',
} as const
