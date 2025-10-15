import { api, handleApiResponse, handleApiError } from './api'
import type {
  Defect,
  DefectFilters,
  CreateDefectData,
  UpdateDefectData,
  DefectStatus,
  ApiResponse,
} from '../types'

// Создаем интерфейс для ответа сервера
interface DefectsResponse {
  defects: Defect[]
  pagination: {
    page: number
    page_size: number
    total: number
    total_pages: number
  }
}

export const defectService = {
  async getDefects(filters?: DefectFilters): Promise<DefectsResponse> {
    try {
      const response = await api.get<ApiResponse<DefectsResponse>>(
        '/api/defects',
        {
          params: filters,
        }
      )
      return handleApiResponse(response)
    } catch (error) {
      throw new Error(handleApiError(error))
    }
  },

  async getDefectById(id: number): Promise<{ defect: Defect }> {
    try {
      const response = await api.get<ApiResponse<{ defect: Defect }>>(
        `/api/defects/${id}`
      )
      return handleApiResponse(response)
    } catch (error) {
      throw new Error(handleApiError(error))
    }
  },

  async createDefect(
    defectData: CreateDefectData
  ): Promise<{ defect: Defect }> {
    try {
      const response = await api.post<ApiResponse<{ defect: Defect }>>(
        '/api/defects',
        defectData
      )
      return handleApiResponse(response)
    } catch (error) {
      throw new Error(handleApiError(error))
    }
  },

  async updateDefect(
    id: number,
    defectData: UpdateDefectData
  ): Promise<{ defect: Defect }> {
    try {
      const response = await api.put<ApiResponse<{ defect: Defect }>>(
        `/api/defects/${id}`,
        defectData
      )
      return handleApiResponse(response)
    } catch (error) {
      throw new Error(handleApiError(error))
    }
  },

  async updateDefectStatus(
    id: number,
    status: DefectStatus
  ): Promise<{ defect: Defect }> {
    try {
      const response = await api.patch<ApiResponse<{ defect: Defect }>>(
        `/api/defects/${id}/status`,
        { status }
      )
      return handleApiResponse(response)
    } catch (error) {
      throw new Error(handleApiError(error))
    }
  },

  async deleteDefect(id: number): Promise<{ message: string }> {
    try {
      const response = await api.delete<ApiResponse<{ message: string }>>(
        `/api/defects/${id}`
      )
      return handleApiResponse(response)
    } catch (error) {
      throw new Error(handleApiError(error))
    }
  },

  async getMyDefects(filters?: DefectFilters): Promise<DefectsResponse> {
    try {
      const response = await api.get<ApiResponse<DefectsResponse>>(
        '/api/defects/my',
        {
          params: filters,
        }
      )
      return handleApiResponse(response)
    } catch (error) {
      throw new Error(handleApiError(error))
    }
  },

  async getDefectHistory(defectId: number): Promise<any> {
    try {
      const response = await api.get<ApiResponse<any>>(
        `/api/defects/${defectId}/history`
      )
      return handleApiResponse(response)
    } catch (error) {
      throw new Error(handleApiError(error))
    }
  },
}
