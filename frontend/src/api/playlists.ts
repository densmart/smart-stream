import apiClient from './client'
import type {
  Playlist,
  CreatePlaylistRequest,
  UpdatePlaylistRequest,
  Media,
  PaginatedResponse,
  ApiResponse
} from '@/types/api'

export const playlistsApi = {
  getPlaylists: async (params?: {
    page?: number
    limit?: number
    search?: string
    type?: string
  }): Promise<PaginatedResponse<Playlist>> => {
    const response = await apiClient.get<PaginatedResponse<Playlist>>('/playlists/', { params })
    return response.data
  },

  getPlaylist: async (id: string): Promise<Playlist> => {
    const response = await apiClient.get<ApiResponse<Playlist>>(`/playlists/${id}/`)
    return response.data.result
  },

  createPlaylist: async (data: CreatePlaylistRequest): Promise<Playlist> => {
    const response = await apiClient.post<ApiResponse<Playlist>>('/playlists/', data)
    return response.data.result
  },

  updatePlaylist: async (id: string, data: UpdatePlaylistRequest): Promise<Playlist> => {
    const response = await apiClient.patch<ApiResponse<Playlist>>(`/playlists/${id}/`, data)
    return response.data.result
  },

  deletePlaylist: async (id: string): Promise<void> => {
    await apiClient.delete(`/playlists/${id}/`)
  },

  getPlaylistMedia: async (
    id: string,
    params?: { page?: number; limit?: number }
  ): Promise<PaginatedResponse<Media>> => {
    const response = await apiClient.get<PaginatedResponse<Media>>(
      `/playlists/${id}/media/`,
      { params }
    )
    return response.data
  },

  addMediaToPlaylist: async (playlistId: string, mediaId: string, order: number): Promise<void> => {
    await apiClient.post(`/playlists/${playlistId}/media/`, { media_id: mediaId, order })
  },

  removeMediaFromPlaylist: async (playlistId: string, mediaId: string): Promise<void> => {
    await apiClient.delete(`/playlists/${playlistId}/media/${mediaId}/`)
  },

  updateMediaOrder: async (playlistId: string, mediaId: string, order: number): Promise<void> => {
    await apiClient.patch(`/playlists/${playlistId}/media/${mediaId}/`, { order })
  },

  // Batch operations
  batchAddMediaToPlaylist: async (
    playlistId: string,
    media: Array<{ mediaId: string; order: number }>
  ): Promise<void> => {
    await apiClient.post(`/playlists/${playlistId}/media/batch/`, {
      media: media.map((m) => ({ media_id: m.mediaId, order: m.order }))
    })
  },

  batchRemoveMediaFromPlaylist: async (playlistId: string, mediaIds: string[]): Promise<void> => {
    await apiClient.delete(`/playlists/${playlistId}/media/batch/`, {
      data: { media_ids: mediaIds }
    })
  },

  batchUpdateMediaOrder: async (
    playlistId: string,
    updates: Array<{ mediaId: string; order: number }>
  ): Promise<void> => {
    await apiClient.patch(`/playlists/${playlistId}/media/batch/`, {
      updates: updates.map((u) => ({ media_id: u.mediaId, order: u.order }))
    })
  }
}