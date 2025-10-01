export const API_BASE_URL = 'http://localhost:8080'

export const DEFECT_STATUSES = {
  new: { label: 'Новая', color: 'bg-gray-100 text-gray-800' },
  in_progress: { label: 'В работе', color: 'bg-blue-100 text-blue-800' },
  on_review: { label: 'На проверке', color: 'bg-yellow-100 text-yellow-800' },
  closed: { label: 'Закрыта', color: 'bg-green-100 text-green-800' },
  cancelled: { label: 'Отменена', color: 'bg-red-100 text-red-800' },
} as const

export const DEFECT_PRIORITIES = {
  low: { label: 'Низкий', color: 'bg-gray-100 text-gray-800' },
  medium: { label: 'Средний', color: 'bg-blue-100 text-blue-800' },
  high: { label: 'Высокий', color: 'bg-orange-100 text-orange-800' },
  critical: { label: 'Критический', color: 'bg-red-100 text-red-800' },
} as const

export const USER_ROLES = {
  engineer: { label: 'Инженер', level: 1 },
  manager: { label: 'Менеджер', level: 2 },
  observer: { label: 'Наблюдатель', level: 0 },
} as const

export const ROLE_PERMISSIONS = {
  engineer: [
    'defects:create',
    'defects:edit_own',
    'defects:change_status',
    'comments:create',
  ],
  manager: [
    'defects:create',
    'defects:edit',
    'defects:change_status',
    'projects:manage',
    'reports:view',
    'reports:export',
    'users:assign',
  ],
  observer: ['defects:view', 'projects:view', 'reports:view'],
} as const

export const PAGINATION_DEFAULTS = {
  page: 1,
  page_size: 20,
} as const
