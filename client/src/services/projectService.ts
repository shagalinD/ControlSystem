import { api, handleApiResponse, handleApiError } from './api'
import type {
  Project,
  ProjectFilters,
  CreateProjectData,
  ApiResponse,
} from '../types'

// Создаем интерфейс для ответа сервера
interface ProjectsResponse {
  projects: Project[]
  pagination: {
    page: number
    page_size: number
    total: number
    total_pages: number
  }
}

export const projectService = {
  async getProjects(filters?: ProjectFilters): Promise<ProjectsResponse> {
    try {
      const response = await api.get<ApiResponse<ProjectsResponse>>(
        '/api/projects',
        {
          params: filters,
        }
      )
      return handleApiResponse(response)
    } catch (error) {
      throw new Error(handleApiError(error))
    }
  },

  async getProjectById(id: number): Promise<{ project: Project }> {
    try {
      const response = await api.get<ApiResponse<{ project: Project }>>(
        `/api/projects/${id}`
      )
      return handleApiResponse(response)
    } catch (error) {
      throw new Error(handleApiError(error))
    }
  },

  async createProject(
    projectData: CreateProjectData
  ): Promise<{ project: Project }> {
    try {
      const response = await api.post<ApiResponse<{ project: Project }>>(
        '/api/projects',
        projectData
      )
      return handleApiResponse(response)
    } catch (error) {
      throw new Error(handleApiError(error))
    }
  },

  async updateProject(
    id: number,
    projectData: Partial<CreateProjectData>
  ): Promise<{ project: Project }> {
    try {
      const response = await api.put<ApiResponse<{ project: Project }>>(
        `/api/projects/${id}`,
        projectData
      )
      return handleApiResponse(response)
    } catch (error) {
      throw new Error(handleApiError(error))
    }
  },

  async deleteProject(id: number): Promise<{ message: string }> {
    try {
      const response = await api.delete<ApiResponse<{ message: string }>>(
        `/api/projects/${id}`
      )
      return handleApiResponse(response)
    } catch (error) {
      throw new Error(handleApiError(error))
    }
  },

  async getProjectDefects(projectId: number, filters?: any): Promise<any> {
    try {
      const response = await api.get<ApiResponse<any>>(
        `/api/projects/${projectId}/defects`,
        {
          params: filters,
        }
      )
      return handleApiResponse(response)
    } catch (error) {
      throw new Error(handleApiError(error))
    }
  },
}
