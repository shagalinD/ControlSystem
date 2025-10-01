import { api, handleApiResponse, handleApiError } from './api'
import type { Attachment, ApiResponse } from '../types'

export const attachmentService = {
  async uploadFile(
    defectId: number,
    file: File
  ): Promise<{ attachment: Attachment }> {
    try {
      const formData = new FormData()
      formData.append('file', file)

      const response = await api.post<ApiResponse<{ attachment: Attachment }>>(
        `/api/attachments/defect/${defectId}`,
        formData,
        {
          headers: {
            'Content-Type': 'multipart/form-data',
          },
        }
      )
      return handleApiResponse(response)
    } catch (error) {
      throw new Error(handleApiError(error))
    }
  },

  async getDefectAttachments(
    defectId: number
  ): Promise<{ attachments: Attachment[] }> {
    try {
      const response = await api.get<
        ApiResponse<{ attachments: Attachment[] }>
      >(`/api/attachments/defect/${defectId}`)
      return handleApiResponse(response)
    } catch (error) {
      throw new Error(handleApiError(error))
    }
  },

  async downloadAttachment(attachmentId: number): Promise<Blob> {
    try {
      const response = await api.get(
        `/api/attachments/${attachmentId}/download`,
        {
          responseType: 'blob',
        }
      )
      return response.data
    } catch (error) {
      throw new Error(handleApiError(error))
    }
  },

  async deleteAttachment(attachmentId: number): Promise<{ message: string }> {
    try {
      const response = await api.delete<ApiResponse<{ message: string }>>(
        `/api/attachments/${attachmentId}`
      )
      return handleApiResponse(response)
    } catch (error) {
      throw new Error(handleApiError(error))
    }
  },
}
