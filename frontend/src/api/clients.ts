import apiClient from './client'
import type {
  Client,
  CreateClientRequest,
  UpdateClientRequest,
  PaginatedResponse,
  ApiResponse
} from '@/types/api'

export const clientsApi = {
  getClients: async (params?: {
    page?: number
    limit?: number
    search?: string
  }): Promise<PaginatedResponse<Client>> => {
    const response = await apiClient.get<PaginatedResponse<Client>>('/clients/', { params })
    return response.data
  },

  getClient: async (id: string): Promise<Client> => {
    const response = await apiClient.get<ApiResponse<Client>>(`/clients/${id}/`)
    return response.data.result
  },

  createClient: async (data: CreateClientRequest): Promise<Client> => {
    const response = await apiClient.post<ApiResponse<Client>>('/clients/', data)
    return response.data.result
  },

  updateClient: async (id: string, data: UpdateClientRequest): Promise<Client> => {
    const response = await apiClient.patch<ApiResponse<Client>>(`/clients/${id}/`, data)
    return response.data.result
  },

  deleteClient: async (id: string): Promise<void> => {
    await apiClient.delete(`/clients/${id}/`)
  }
}