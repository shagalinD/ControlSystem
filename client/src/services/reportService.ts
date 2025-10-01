import { api, handleApiResponse, handleApiError } from './api'
import type { DefectFilters, ApiResponse } from '../types'

export interface DefectsReport {
  total_defects: number
  defects_by_status: Array<{ status: string; count: number }>
  defects_by_priority: Array<{ priority: string; count: number }>
  overdue_defects: number
  avg_resolution_time: number
}

export interface ProjectReport {
  project_id: number
  project_name: string
  total_defects: number
  defects_by_status: Array<{ status: string; count: number }>
  defects_by_priority: Array<{ priority: string; count: number }>
  completion_percentage: number
}

export const reportService = {
  async getDefectsReport(
    filters?: DefectFilters
  ): Promise<{ report: DefectsReport }> {
    try {
      const response = await api.get<ApiResponse<{ report: DefectsReport }>>(
        '/api/reports/defects',
        {
          params: filters,
        }
      )
      return handleApiResponse(response)
    } catch (error) {
      throw new Error(handleApiError(error))
    }
  },

  async getProjectReport(
    projectId: number
  ): Promise<{ report: ProjectReport }> {
    try {
      const response = await api.get<ApiResponse<{ report: ProjectReport }>>(
        `/api/reports/projects/${projectId}`
      )
      return handleApiResponse(response)
    } catch (error) {
      throw new Error(handleApiError(error))
    }
  },

  async exportDefectsToCSV(filters?: DefectFilters): Promise<Blob> {
    try {
      const response = await api.get('/api/reports/defects/export', {
        params: filters,
        responseType: 'blob',
      })
      return response.data
    } catch (error) {
      throw new Error(handleApiError(error))
    }
  },

  async exportProjectReport(projectId: number): Promise<Blob> {
    try {
      const response = await api.get(
        `/api/reports/projects/${projectId}/export`,
        {
          responseType: 'blob',
        }
      )
      return response.data
    } catch (error) {
      throw new Error(handleApiError(error))
    }
  },
}
