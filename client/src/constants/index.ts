export const API_BASE_URL = 'http://localhost:8080'

export const ROLES = {
  engineer: 'engineer',
  manager: 'manager',
  observer: 'observer',
} as const

export const DEFECT_STATUSES = {
  new: { value: 'new', label: 'Новая', color: 'bg-gray-100 text-gray-800' },
  in_progress: {
    value: 'in_progress',
    label: 'В работе',
    color: 'bg-blue-100 text-blue-800',
  },
  on_review: {
    value: 'on_review',
    label: 'На проверке',
    color: 'bg-yellow-100 text-yellow-800',
  },
  closed: {
    value: 'closed',
    label: 'Закрыта',
    color: 'bg-green-100 text-green-800',
  },
  cancelled: {
    value: 'cancelled',
    label: 'Отменена',
    color: 'bg-red-100 text-red-800',
  },
} as const

export const DEFECT_PRIORITIES = {
  low: { value: 'low', label: 'Низкий', color: 'bg-gray-100 text-gray-800' },
  medium: {
    value: 'medium',
    label: 'Средний',
    color: 'bg-yellow-100 text-yellow-800',
  },
  high: {
    value: 'high',
    label: 'Высокий',
    color: 'bg-orange-100 text-orange-800',
  },
  critical: {
    value: 'critical',
    label: 'Критический',
    color: 'bg-red-100 text-red-800',
  },
} as const

export const ROLE_PERMISSIONS = {
  engineer: [
    'defects:create',
    'defects:edit_own',
    'defects:view',
    'comments:create',
  ],
  manager: [
    'defects:create',
    'defects:edit',
    'defects:view',
    'projects:manage',
    'reports:view',
    'comments:create',
  ],
  observer: ['defects:view', 'reports:view'],
} as const

export const STORAGE_KEYS = {
  AUTH_TOKEN: 'auth_token',
  USER_DATA: 'user_data',
} as const

export const ROUTES = {
  HOME: '/',
  LOGIN: '/login',
  REGISTER: '/register',
  DEFECTS: '/defects',
  DEFECTS_CREATE: '/defects/create',
  DEFECTS_EDIT: '/defects/:id/edit',
  DEFECTS_VIEW: '/defects/:id',
  PROJECTS: '/projects',
  PROJECTS_CREATE: '/projects/create',
  PROJECTS_VIEW: '/projects/:id',
  PROJECTS_EDIT: '/projects/:id/edit',
  REPORTS: '/reports',
  PROFILE: '/profile',
} as const
