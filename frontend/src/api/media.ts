import apiClient from './client'
import type {
  Media,
  CreateMediaRequest,
  UpdateMediaRequest,
  PaginatedResponse,
  ApiResponse,
  BrowseMediaFilesResponse
} from '@/types/api'

export const mediaApi = {
  getMedia: async (params?: {
    page?: number
    limit?: number
    offset?: number
    name?: string
    format?: string
    only_unassigned?: boolean
  }): Promise<PaginatedResponse<Media>> => {
    const response = await apiClient.get<PaginatedResponse<Media>>('/media/', { params })
    return response.data
  },

  getMediaById: async (id: string): Promise<Media> => {
    const response = await apiClient.get<ApiResponse<Media>>(`/media/${id}/`)
    return response.data.result
  },

  createMedia: async (data: CreateMediaRequest): Promise<Media> => {
    const response = await apiClient.post<ApiResponse<Media>>('/media/', data)
    return response.data.result
  },

  updateMedia: async (id: string, data: UpdateMediaRequest): Promise<Media> => {
    const response = await apiClient.patch<ApiResponse<Media>>(`/media/${id}/`, data)
    return response.data.result
  },

  deleteMedia: async (id: string): Promise<void> => {
    await apiClient.delete(`/media/${id}/`)
  },

  uploadPoster: async (file: File): Promise<{ filename: string }> => {
    const formData = new FormData()
    formData.append('file', file)

    const response = await apiClient.post<ApiResponse<{ filename: string }>>('/upload/poster/', formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    })
    return response.data.result
  },

  browseMediaFiles: async (path: string = ''): Promise<ApiResponse<BrowseMediaFilesResponse>> => {
    const response = await apiClient.get<ApiResponse<BrowseMediaFilesResponse>>('/media/browse', {
      params: { path }
    })
    return response.data
  }
}