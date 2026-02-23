<template>
  <div class="read-page">
    <ReadingHeader
      :story-title="storyTitle"
      @show-history="showHistory = true"
    />

    <main class="content">
      <SceneDisplay
        v-if="currentScene"
        :scene="currentScene"
        :choices="choices"
        @choose="makeChoice"
        @restart="restart"
      />

      <div v-else class="loading">
        <p>Đang tải...</p>
      </div>
    </main>

    <HistoryModal
      v-if="showHistory"
      :history="readingHistory"
      @close="showHistory = false"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import ReadingHeader from '@/components/reading/reading-header.vue'
import SceneDisplay from '@/components/reading/scene-display.vue'
import HistoryModal from '@/components/reading/history-modal.vue'

interface Scene {
  id: string
  title?: string
  content: string
  imageUrl?: string
  isEnd: boolean
}

interface Choice {
  id: string
  text: string
  nextSceneId: string
}

interface HistoryItem {
  sceneId: string
  sceneTitle?: string
  choiceText?: string
}

const route = useRoute()
const storyId = route.params.storyId as string
void storyId // Will be used when API is implemented

const storyTitle = ref('Tên truyện')
const currentScene = ref<Scene | null>(null)
const choices = ref<Choice[]>([])
const showHistory = ref(false)
const readingHistory = ref<HistoryItem[]>([])

const loadScene = async (sceneId: string) => {
  currentScene.value = {
    id: sceneId,
    title: 'Cảnh 1',
    content:
      'Đây là nội dung của cảnh đầu tiên.\n\nBạn đang đứng trước một cánh cửa bí ẩn. Ánh sáng yếu ớt từ bên trong chiếu ra.',
    isEnd: false,
  }

  choices.value = [
    { id: '1', text: 'Mở cửa và bước vào', nextSceneId: 'scene-2' },
    { id: '2', text: 'Quay lại', nextSceneId: 'scene-3' },
  ]
}

const makeChoice = (choice: Choice) => {
  readingHistory.value.push({
    sceneId: currentScene.value!.id,
    sceneTitle: currentScene.value?.title,
    choiceText: choice.text,
  })
  loadScene(choice.nextSceneId)
}

const restart = () => {
  readingHistory.value = []
  loadScene('scene-1')
}

onMounted(() => {
  loadScene('scene-1')
})
</script>

<style scoped>
.read-page {
  min-height: 100vh;
  background: #f9fafb;
  display: flex;
  flex-direction: column;
}

.content {
  flex: 1;
  max-width: 800px;
  width: 100%;
  margin: 0 auto;
  padding: 2rem 1rem;
}

.loading {
  text-align: center;
  padding: 4rem;
  color: #666;
}
</style>
