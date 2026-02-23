export interface StorySummary {
  id: string
  title: string
  status: string
}

export interface GraphNode {
  id: string
  type: string
  position: {
    x: number
    y: number
  }
  data: {
    content: string
    is_start: boolean
    is_end: boolean
    image_url?: string
  }
}

export interface GraphEdge {
  id: string
  source: string
  target: string
  label: string
}

export interface GraphResponse {
  story: StorySummary
  nodes: GraphNode[]
  edges: GraphEdge[]
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

export const graphApi = {
  async loadGraph(storyId: string): Promise<GraphResponse> {
    const response = await fetchWithAuth(`/api/creator/stories/${storyId}/graph`, {
      method: 'GET',
    })
    return response.json()
  },
}
