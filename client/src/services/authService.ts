import { api, handleApiResponse, handleApiError } from './api'
import type {
  LoginFormData,
  RegisterFormData,
  AuthResponse,
  User,
  ApiResponse,
} from '../types'

export const authService = {
  async login(credentials: LoginFormData): Promise<AuthResponse> {
    try {
      const response = await api.post<ApiResponse<AuthResponse>>(
        '/auth/login',
        credentials
      )
      return handleApiResponse(response)
    } catch (error) {
      throw new Error(handleApiError(error))
    }
  },

  async register(userData: RegisterFormData): Promise<{ user: User }> {
    try {
      const response = await api.post<ApiResponse<{ user: User }>>(
        '/auth/register',
        userData
      )
      return handleApiResponse(response)
    } catch (error) {
      throw new Error(handleApiError(error))
    }
  },

  async getCurrentUser(): Promise<{ user: User }> {
    try {
      const response = await api.get<ApiResponse<{ user: User }>>('/api/me')
      return handleApiResponse(response)
    } catch (error) {
      throw new Error(handleApiError(error))
    }
  },

  async refreshToken(): Promise<{ token: string }> {
    try {
      const response = await api.post<ApiResponse<{ token: string }>>(
        '/auth/refresh'
      )
      return handleApiResponse(response)
    } catch (error) {
      throw new Error(handleApiError(error))
    }
  },

  async logout(): Promise<void> {
    try {
      await api.post('/auth/logout')
    } catch (error) {
      // Даже если запрос не удался, очищаем локальное хранилище
      console.error('Logout error:', error)
    } finally {
      // Всегда очищаем локальные данные
      localStorage.removeItem('auth_token')
      localStorage.removeItem('user_data')
    }
  },

  async requestPasswordReset(email: string): Promise<{ message: string }> {
    try {
      const response = await api.post<ApiResponse<{ message: string }>>(
        '/auth/forgot-password',
        { email }
      )
      return handleApiResponse(response)
    } catch (error) {
      throw new Error(handleApiError(error))
    }
  },

  async resetPassword(
    token: string,
    newPassword: string
  ): Promise<{ message: string }> {
    try {
      const response = await api.post<ApiResponse<{ message: string }>>(
        '/auth/reset-password',
        {
          token,
          new_password: newPassword,
        }
      )
      return handleApiResponse(response)
    } catch (error) {
      throw new Error(handleApiError(error))
    }
  },
}
