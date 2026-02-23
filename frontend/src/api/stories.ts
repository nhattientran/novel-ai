import type { Story, CreateStoryRequest, UpdateStoryRequest } from '../stores/stories'

const API_BASE_URL = (import.meta as any).env?.VITE_API_URL || 'http://localhost:8080'

async function fetchWithAuth(url: string, options: RequestInit = {}): Promise<Response> {
  const response = await fetch(`${API_BASE_URL}${url}`, {
    ...options,
    headers: {
      'Content-Type': 'application/json',
      ...options.headers,
    },
    credentials: 'include',
  })

  if (!response.ok) {
    const error = await response.json().catch(() => ({ error: 'Unknown error' }))
    throw new Error(error.error || `HTTP ${response.status}`)
  }

  return response
}

export const storiesApi = {
  async list(): Promise<Story[]> {
    const response = await fetchWithAuth('/api/creator/stories', {
      method: 'GET',
    })
    return response.json()
  },

  async create(data: CreateStoryRequest): Promise<Story> {
    const response = await fetchWithAuth('/api/creator/stories', {
      method: 'POST',
      body: JSON.stringify(data),
    })
    return response.json()
  },

  async get(id: string): Promise<Story> {
    const response = await fetchWithAuth(`/api/creator/stories/${id}`, {
      method: 'GET',
    })
    return response.json()
  },

  async update(id: string, data: UpdateStoryRequest): Promise<Story> {
    const response = await fetchWithAuth(`/api/creator/stories/${id}`, {
      method: 'PATCH',
      body: JSON.stringify(data),
    })
    return response.json()
  },

  async delete(id: string): Promise<void> {
    await fetchWithAuth(`/api/creator/stories/${id}`, {
      method: 'DELETE',
    })
  },
}
