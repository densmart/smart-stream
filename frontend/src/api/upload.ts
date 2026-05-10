import apiClient from './client'
import type { ApiResponse } from '@/types/api'

export interface PosterUploadResponse {
  filename: string
  path: string
}

export const uploadApi = {
  uploadPoster: async (file: File): Promise<PosterUploadResponse> => {
    const formData = new FormData()
    formData.append('file', file)

    const response = await apiClient.post<ApiResponse<PosterUploadResponse>>('/upload/poster/', formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    })
    return response.data.result
  },

  getPosterUrl: (posterPath: string | null): string | null => {
    if (!posterPath) return null
    const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1'
    return `${API_BASE_URL}/static/posters/${posterPath}`
  }
}
