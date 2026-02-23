import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

// Auth Store
export const useAuthStore = defineStore('auth', () => {
  const token = ref<string | null>(localStorage.getItem('token'))
  const user = ref<{ id: string; username: string; email: string; role: string } | null>(null)

  const isAuthenticated = computed(() => !!token.value)
  const isCreator = computed(() => user.value?.role === 'creator')

  function setToken(newToken: string) {
    token.value = newToken
    localStorage.setItem('token', newToken)
  }

  function clearToken() {
    token.value = null
    user.value = null
    localStorage.removeItem('token')
  }

  function setUser(userData: { id: string; username: string; email: string; role: string }) {
    user.value = userData
  }

  return {
    token,
    user,
    isAuthenticated,
    isCreator,
    setToken,
    clearToken,
    setUser,
  }
})

// Reading Progress Store
export const useReadingStore = defineStore('reading', () => {
  const currentStoryId = ref<string | null>(null)
  const currentSceneId = ref<string | null>(null)
  const history = ref<string[]>([])

  function startReading(storyId: string, sceneId: string) {
    currentStoryId.value = storyId
    currentSceneId.value = sceneId
    history.value = [sceneId]
    saveToLocalStorage()
  }

  function goToScene(sceneId: string) {
    currentSceneId.value = sceneId
    history.value.push(sceneId)
    saveToLocalStorage()
  }

  function goBack() {
    if (history.value.length > 1) {
      history.value.pop()
      currentSceneId.value = history.value[history.value.length - 1]
      saveToLocalStorage()
    }
  }

  function saveToLocalStorage() {
    if (currentStoryId.value) {
      localStorage.setItem(
        `reading:${currentStoryId.value}`,
        JSON.stringify({
          sceneId: currentSceneId.value,
          history: history.value,
        })
      )
    }
  }

  function loadFromLocalStorage(storyId: string) {
    const data = localStorage.getItem(`reading:${storyId}`)
    if (data) {
      const parsed = JSON.parse(data)
      currentStoryId.value = storyId
      currentSceneId.value = parsed.sceneId
      history.value = parsed.history || []
      return true
    }
    return false
  }

  return {
    currentStoryId,
    currentSceneId,
    history,
    startReading,
    goToScene,
    goBack,
    loadFromLocalStorage,
  }
})
