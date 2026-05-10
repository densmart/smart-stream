import { defineStore } from 'pinia'
import { ref } from 'vue'
import { authApi } from '@/api'
import type { SignInRequest } from '@/types/api'

export const useAuthStore = defineStore('auth', () => {
  const isAuthenticated = ref(authApi.isAuthenticated())
  const loading = ref(false)
  const error = ref<string | null>(null)

  const signIn = async (credentials: SignInRequest): Promise<boolean> => {
    try {
      loading.value = true
      error.value = null

      const response = await authApi.signIn(credentials)
      authApi.setToken(response.access_token)
      isAuthenticated.value = true

      return true
    } catch (err: any) {
      error.value = err.response?.data?.error || 'Failed to sign in'
      isAuthenticated.value = false
      return false
    } finally {
      loading.value = false
    }
  }

  const logout = () => {
    authApi.logout()
    isAuthenticated.value = false
  }

  const checkAuth = () => {
    isAuthenticated.value = authApi.isAuthenticated()
  }

  return {
    isAuthenticated,
    loading,
    error,
    signIn,
    logout,
    checkAuth
  }
})