const API_BASE_URL = (import.meta as any).env?.VITE_API_URL || 'http://localhost:8080'

export interface UploadResponse {
  url: string
}

export const uploadsApi = {
  async uploadImage(file: File): Promise<UploadResponse> {
    const formData = new FormData()
    formData.append('image', file)

    const response = await fetch(`${API_BASE_URL}/api/uploads/images`, {
      method: 'POST',
      body: formData,
      credentials: 'include',
    })

    if (!response.ok) {
      const error = await response.json().catch(() => ({ error: 'Unknown error' }))
      throw new Error(error.error || `HTTP ${response.status}`)
    }

    return response.json()
  },
}
