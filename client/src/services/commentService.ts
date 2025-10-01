import { api, handleApiResponse, handleApiError } from './api'
import type { Comment, ApiResponse } from '../types'

export const commentService = {
  async getDefectComments(defectId: number): Promise<{ comments: Comment[] }> {
    try {
      const response = await api.get<ApiResponse<{ comments: Comment[] }>>(
        `/api/comments/defect/${defectId}`
      )
      return handleApiResponse(response)
    } catch (error) {
      throw new Error(handleApiError(error))
    }
  },

  async createComment(
    defectId: number,
    text: string
  ): Promise<{ comment: Comment }> {
    try {
      const response = await api.post<ApiResponse<{ comment: Comment }>>(
        `/api/comments/defect/${defectId}`,
        {
          text,
          defect_id: defectId,
        }
      )
      return handleApiResponse(response)
    } catch (error) {
      throw new Error(handleApiError(error))
    }
  },

  async updateComment(
    commentId: number,
    text: string
  ): Promise<{ comment: Comment }> {
    try {
      const response = await api.put<ApiResponse<{ comment: Comment }>>(
        `/api/comments/${commentId}`,
        { text }
      )
      return handleApiResponse(response)
    } catch (error) {
      throw new Error(handleApiError(error))
    }
  },

  async deleteComment(commentId: number): Promise<{ message: string }> {
    try {
      const response = await api.delete<ApiResponse<{ message: string }>>(
        `/api/comments/${commentId}`
      )
      return handleApiResponse(response)
    } catch (error) {
      throw new Error(handleApiError(error))
    }
  },
}
