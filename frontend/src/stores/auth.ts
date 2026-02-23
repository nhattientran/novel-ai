import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { authApi } from '../api/auth'

export interface User {
  id: string
  username: string
  email: string
  role: string
}

// Auth Store - uses HttpOnly cookie (no token in localStorage)
export const useAuthStore = defineStore('auth', () => {
  const user = ref<User | null>(null)
  const isLoading = ref(false)
  const error = ref<string | null>(null)

  const isAuthenticated = computed(() => !!user.value)
  const isCreator = computed(() => user.value?.role === 'creator')

  // Fetch current user on app init
  async function fetchMe() {
    isLoading.value = true
    error.value = null
    try {
      const userData = await authApi.me()
      user.value = userData
      return userData
    } catch (err) {
      user.value = null
      return null
    } finally {
      isLoading.value = false
    }
  }

  async function login(email: string, password: string) {
    isLoading.value = true
    error.value = null
    try {
      const userData = await authApi.login({ email, password })
      user.value = userData
      return userData
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Login failed'
      throw err
    } finally {
      isLoading.value = false
    }
  }

  async function register(username: string, email: string, password: string, role: 'creator' | 'reader') {
    isLoading.value = true
    error.value = null
    try {
      const userData = await authApi.register({ username, email, password, role })
      user.value = userData
      return userData
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Registration failed'
      throw err
    } finally {
      isLoading.value = false
    }
  }

  async function logout() {
    isLoading.value = true
    try {
      await authApi.logout()
      user.value = null
    } catch (err) {
      console.error('Logout error:', err)
    } finally {
      isLoading.value = false
    }
  }

  function clearError() {
    error.value = null
  }

  return {
    user,
    isLoading,
    error,
    isAuthenticated,
    isCreator,
    fetchMe,
    login,
    register,
    logout,
    clearError,
  }
})
