import { api, handleApiResponse, handleApiError } from './api'
import type { User, ApiResponse, UpdateProfileData } from '../types'
interface UsersResponse {
  data: User[]
  pagination: {
    page: number
    page_size: number
    total: number
    total_pages: number
  }
}
interface EngineersResponse {
  engineers: User[]
  pagination: {
    page: number
    page_size: number
    total: number
    total_pages: number
  }
}
interface MaganersResponse {
  managers: User[]
  pagination: {
    page: number
    page_size: number
    total: number
    total_pages: number
  }
}
export const userService = {
  async getUsers(): Promise<UsersResponse> {
    try {
      const response = await api.get<ApiResponse<UsersResponse>>('/api/users')
      return handleApiResponse(response)
    } catch (error) {
      throw new Error(handleApiError(error))
    }
  },

  async getEngineers(): Promise<EngineersResponse> {
    try {
      const response = await api.get<ApiResponse<EngineersResponse>>(
        '/api/users/engineers'
      )
      return handleApiResponse(response)
    } catch (error) {
      throw new Error(handleApiError(error))
    }
  },

  async getManagers(): Promise<MaganersResponse> {
    try {
      const response = await api.get<ApiResponse<MaganersResponse>>(
        '/api/users/managers'
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
