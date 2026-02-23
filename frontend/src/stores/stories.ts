import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { storiesApi } from '../api/stories'

export interface Story {
  id: string
  title: string
  summary: string
  cover_image: string
  status: 'draft' | 'published'
  created_at: string
}

export interface CreateStoryRequest {
  title: string
  summary?: string
  cover_image?: string
}

export interface UpdateStoryRequest {
  title?: string
  summary?: string
  cover_image?: string
}

export const useStoriesStore = defineStore('stories', () => {
  // State
  const stories = ref<Story[]>([])
  const currentStory = ref<Story | null>(null)
  const isLoading = ref(false)
  const error = ref<string | null>(null)

  // Getters
  const draftStories = computed(() => stories.value.filter(s => s.status === 'draft'))
  const publishedStories = computed(() => stories.value.filter(s => s.status === 'published'))
  const hasStories = computed(() => stories.value.length > 0)

  // Actions
  async function fetchStories() {
    isLoading.value = true
    error.value = null

    try {
      stories.value = await storiesApi.list()
    } catch (err: any) {
      error.value = err.message || 'Failed to fetch stories'
      throw err
    } finally {
      isLoading.value = false
    }
  }

  async function fetchStory(id: string) {
    isLoading.value = true
    error.value = null

    try {
      currentStory.value = await storiesApi.get(id)
    } catch (err: any) {
      error.value = err.message || 'Failed to fetch story'
      throw err
    } finally {
      isLoading.value = false
    }
  }

  async function createStory(data: CreateStoryRequest) {
    isLoading.value = true
    error.value = null

    try {
      const newStory = await storiesApi.create(data)
      stories.value.unshift(newStory)
      return newStory
    } catch (err: any) {
      error.value = err.message || 'Failed to create story'
      throw err
    } finally {
      isLoading.value = false
    }
  }

  async function updateStory(id: string, data: UpdateStoryRequest) {
    isLoading.value = true
    error.value = null

    try {
      const updatedStory = await storiesApi.update(id, data)

      // Update in stories list
      const index = stories.value.findIndex(s => s.id === id)
      if (index !== -1) {
        stories.value[index] = updatedStory
      }

      // Update current story if it's the same
      if (currentStory.value?.id === id) {
        currentStory.value = updatedStory
      }

      return updatedStory
    } catch (err: any) {
      error.value = err.message || 'Failed to update story'
      throw err
    } finally {
      isLoading.value = false
    }
  }

  async function deleteStory(id: string) {
    isLoading.value = true
    error.value = null

    try {
      await storiesApi.delete(id)
      stories.value = stories.value.filter(s => s.id !== id)

      if (currentStory.value?.id === id) {
        currentStory.value = null
      }
    } catch (err: any) {
      error.value = err.message || 'Failed to delete story'
      throw err
    } finally {
      isLoading.value = false
    }
  }

  function clearError() {
    error.value = null
  }

  function clearCurrentStory() {
    currentStory.value = null
  }

  return {
    // State
    stories,
    currentStory,
    isLoading,
    error,
    // Getters
    draftStories,
    publishedStories,
    hasStories,
    // Actions
    fetchStories,
    fetchStory,
    createStory,
    updateStory,
    deleteStory,
    clearError,
    clearCurrentStory,
  }
})
