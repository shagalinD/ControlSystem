import type { UserRole } from '../types'
import { ROLE_PERMISSIONS } from '../constants'

export const hasPermission = (
  userRole: UserRole | undefined,
  permission: string
): boolean => {
  if (!userRole) return false
  const permissions = ROLE_PERMISSIONS[userRole]
  return permissions.includes(permission as any)
}

export const isManagerOrAbove = (userRole: UserRole | undefined): boolean => {
  return userRole === 'manager' || userRole === 'observer'
}

export const canEditDefect = (
  userRole: UserRole | undefined,
  defectAuthorId?: number,
  currentUserId?: number
): boolean => {
  if (!userRole || !currentUserId) return false

  if (userRole === 'manager') return true
  if (userRole === 'engineer' && defectAuthorId === currentUserId) return true

  return false
}

export const getStoredToken = (): string | null => {
  return localStorage.getItem('auth_token')
}

export const setStoredToken = (token: string): void => {
  localStorage.setItem('auth_token', token)
}

export const removeStoredToken = (): void => {
  localStorage.removeItem('auth_token')
}
