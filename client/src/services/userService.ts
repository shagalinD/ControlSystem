import { api, handleApiResponse, handleApiError } from './api'
import type {
  User,
  ApiResponse,
  PaginatedResponse,
  UpdateProfileData,
} from '../types'

export const userService = {
  async getUsers(): Promise<PaginatedResponse<User>> {
    try {
      const response = await api.get<ApiResponse<PaginatedResponse<User>>>(
        '/api/users'
      )
      return handleApiResponse(response)
    } catch (error) {
      throw new Error(handleApiError(error))
    }
  },

  async getUserById(id: number): Promise<{ user: User }> {
    try {
      const response = await api.get<ApiResponse<{ user: User }>>(
        `/api/users/${id}`
      )
      return handleApiResponse(response)
    } catch (error) {
      throw new Error(handleApiError(error))
    }
  },

  async getEngineers(): Promise<PaginatedResponse<User>> {
    try {
      const response = await api.get<ApiResponse<PaginatedResponse<User>>>(
        '/api/users/engineers'
      )
      return handleApiResponse(response)
    } catch (error) {
      throw new Error(handleApiError(error))
    }
  },

  async getManagers(): Promise<PaginatedResponse<User>> {
    try {
      const response = await api.get<ApiResponse<PaginatedResponse<User>>>(
        '/api/users/managers'
      )
      return handleApiResponse(response)
    } catch (error) {
      throw new Error(handleApiError(error))
    }
  },
  // Добавляем методы в userService:

  async updateUserProfile(
    userId: number,
    userData: UpdateProfileData
  ): Promise<{ user: User }> {
    try {
      const response = await api.put<ApiResponse<{ user: User }>>(
        `/api/users/${userId}`,
        userData
      )
      return handleApiResponse(response)
    } catch (error) {
      throw new Error(handleApiError(error))
    }
  },

  async changePassword(
    currentPassword: string,
    newPassword: string
  ): Promise<{ message: string }> {
    try {
      const response = await api.post<ApiResponse<{ message: string }>>(
        '/api/users/change-password',
        {
          current_password: currentPassword,
          new_password: newPassword,
        }
      )
      return handleApiResponse(response)
    } catch (error) {
      throw new Error(handleApiError(error))
    }
  },
}
