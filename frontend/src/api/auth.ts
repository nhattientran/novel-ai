import type { User } from '../stores/auth'

const API_BASE_URL = (import.meta as any).env?.VITE_API_URL || 'http://localhost:8080'

export interface RegisterRequest {
  username: string
  email: string
  password: string
  role: 'creator' | 'reader'
}

export interface LoginRequest {
  email: string
  password: string
}

export interface AuthResponse {
  id: string
  username: string
  email: string
  role: string
}

async function fetchWithAuth(url: string, options: RequestInit = {}): Promise<Response> {
  const response = await fetch(`${API_BASE_URL}${url}`, {
    ...options,
    headers: {
      'Content-Type': 'application/json',
      ...options.headers,
    },
    credentials: 'include', // Important: include cookies
  })

  if (!response.ok) {
    const error = await response.json().catch(() => ({ error: 'Unknown error' }))
    throw new Error(error.error || `HTTP ${response.status}`)
  }

  return response
}

export const authApi = {
  async register(data: RegisterRequest): Promise<AuthResponse> {
    const response = await fetchWithAuth('/api/auth/register', {
      method: 'POST',
      body: JSON.stringify(data),
    })
    return response.json()
  },

  async login(data: LoginRequest): Promise<AuthResponse> {
    const response = await fetchWithAuth('/api/auth/login', {
      method: 'POST',
      body: JSON.stringify(data),
    })
    return response.json()
  },

  async logout(): Promise<void> {
    await fetchWithAuth('/api/auth/logout', {
      method: 'POST',
    })
  },

  async me(): Promise<User> {
    const response = await fetchWithAuth('/api/me', {
      method: 'GET',
    })
    return response.json()
  },
}
