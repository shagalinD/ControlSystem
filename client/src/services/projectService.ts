import { api, handleApiResponse, handleApiError } from './api'
import type {
  Project,
  ProjectFilters,
  PaginatedResponse,
  CreateProjectData,
  UpdateProjectData,
  ApiResponse,
} from '../types'

export const projectService = {
  async getProjects(
    filters?: ProjectFilters
  ): Promise<PaginatedResponse<Project>> {
    try {
      const response = await api.get<ApiResponse<PaginatedResponse<Project>>>(
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
    projectData: UpdateProjectData
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
}
