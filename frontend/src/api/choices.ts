export interface Choice {
  from_scene_id: string
  to_scene_id: string
  choice_text: string
}

export interface CreateChoiceRequest {
  from_scene_id: string
  to_scene_id: string
  choice_text: string
}

export interface UpdateChoiceRequest {
  from_scene_id: string
  to_scene_id: string
  choice_text: string
}

export interface DeleteChoiceRequest {
  from_scene_id: string
  to_scene_id: string
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

export const choicesApi = {
  async create(storyId: string, data: CreateChoiceRequest): Promise<Choice> {
    const response = await fetchWithAuth(`/api/creator/stories/${storyId}/choices`, {
      method: 'POST',
      body: JSON.stringify(data),
    })
    return response.json()
  },

  async update(storyId: string, data: UpdateChoiceRequest): Promise<Choice> {
    const response = await fetchWithAuth(`/api/creator/stories/${storyId}/choices`, {
      method: 'PATCH',
      body: JSON.stringify(data),
    })
    return response.json()
  },

  async delete(storyId: string, data: DeleteChoiceRequest): Promise<void> {
    await fetchWithAuth(`/api/creator/stories/${storyId}/choices`, {
      method: 'DELETE',
      body: JSON.stringify(data),
    })
  },
}
