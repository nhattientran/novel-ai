export interface Scene {
  id: string
  content: string
  image_url?: string
  is_end: boolean
  pos_x: number
  pos_y: number
}

export interface CreateSceneRequest {
  content: string
  image_url?: string
  is_end?: boolean
  pos_x?: number
  pos_y?: number
}

export interface UpdateSceneRequest {
  content?: string
  image_url?: string | null
  is_end?: boolean
  pos_x?: number
  pos_y?: number
}

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

export const scenesApi = {
  async create(storyId: string, data: CreateSceneRequest): Promise<Scene> {
    const response = await fetchWithAuth(`/api/creator/stories/${storyId}/scenes`, {
      method: 'POST',
      body: JSON.stringify(data),
    })
    return response.json()
  },

  async get(storyId: string, sceneId: string): Promise<Scene> {
    const response = await fetchWithAuth(`/api/creator/stories/${storyId}/scenes/${sceneId}`, {
      method: 'GET',
    })
    return response.json()
  },

  async update(storyId: string, sceneId: string, data: UpdateSceneRequest): Promise<Scene> {
    const response = await fetchWithAuth(`/api/creator/stories/${storyId}/scenes/${sceneId}`, {
      method: 'PATCH',
      body: JSON.stringify(data),
    })
    return response.json()
  },

  async delete(storyId: string, sceneId: string): Promise<void> {
    await fetchWithAuth(`/api/creator/stories/${storyId}/scenes/${sceneId}`, {
      method: 'DELETE',
    })
  },

  async setStartScene(storyId: string, sceneId: string): Promise<{ message: string; start_scene_id: string }> {
    const response = await fetchWithAuth(`/api/creator/stories/${storyId}/start/${sceneId}`, {
      method: 'PUT',
    })
    return response.json()
  },
}
