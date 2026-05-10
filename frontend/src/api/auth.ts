import apiClient from './client'
import type { SignInRequest, SignInResponse } from '@/types/api'

export const authApi = {
  signIn: async (credentials: SignInRequest): Promise<SignInResponse> => {
    const response = await apiClient.post<SignInResponse>('/auth/sign-in/', credentials)
    return response.data
  },

  logout: () => {
    localStorage.removeItem('access_token')
  },

  getToken: (): string | null => {
    return localStorage.getItem('access_token')
  },

  setToken: (token: string): void => {
    localStorage.setItem('access_token', token)
  },

  isAuthenticated: (): boolean => {
    return !!localStorage.getItem('access_token')
  }
}