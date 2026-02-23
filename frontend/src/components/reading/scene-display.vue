<template>
  <div class="scene">
    <div v-if="scene.imageUrl" class="scene-image">
      <img :src="scene.imageUrl" :alt="sceneTitle" />
    </div>

    <div class="scene-content">
      <h2 v-if="sceneTitle">{{ sceneTitle }}</h2>
      <p v-for="(paragraph, index) in paragraphs" :key="index">
        {{ paragraph }}
      </p>
    </div>

    <div class="choices">
      <button
        v-for="choice in choices"
        :key="choice.id"
        class="choice-btn"
        @click="emit('choose', choice)"
      >
        {{ choice.text }}
      </button>
    </div>

    <div v-if="scene.isEnd" class="end-scene">
      <p>🎉 Bạn đã hoàn thành câu chuyện!</p>
      <button class="btn btn-primary" @click="emit('restart')">Đọc lại từ đầu</button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

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

const props = defineProps<{
  scene: Scene
  choices: Choice[]
}>()

const emit = defineEmits<{
  choose: [choice: Choice]
  restart: []
}>()

const sceneTitle = computed(() => props.scene.title)
const paragraphs = computed(() =>
  props.scene.content.split('\n').filter((p) => p.trim())
)
</script>

<style scoped>
.scene {
  background: white;
  border-radius: 1rem;
  overflow: hidden;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
}

.scene-image {
  width: 100%;
  height: 300px;
  overflow: hidden;
}

.scene-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.scene-content {
  padding: 2rem;
}

.scene-content h2 {
  font-size: 1.5rem;
  margin-bottom: 1rem;
  color: #333;
}

.scene-content p {
  line-height: 1.8;
  color: #444;
  margin-bottom: 1rem;
  font-size: 1.1rem;
}

.choices {
  padding: 0 2rem 2rem;
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.choice-btn {
  padding: 1rem 1.5rem;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border: none;
  border-radius: 0.75rem;
  font-size: 1rem;
  cursor: pointer;
  transition: transform 0.2s, box-shadow 0.2s;
  text-align: left;
}

.choice-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.4);
}

.end-scene {
  padding: 2rem;
  text-align: center;
  border-top: 1px solid #e5e7eb;
}

.end-scene p {
  font-size: 1.25rem;
  margin-bottom: 1rem;
  color: #333;
}

.btn {
  padding: 0.75rem 1.5rem;
  border-radius: 0.5rem;
  border: none;
  font-size: 1rem;
  cursor: pointer;
  transition: all 0.3s;
}

.btn-primary {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
}

.btn-primary:hover {
  opacity: 0.9;
}

@media (max-width: 640px) {
  .scene-content {
    padding: 1.5rem;
  }

  .scene-content p {
    font-size: 1rem;
  }

  .choices {
    padding: 0 1.5rem 1.5rem;
  }
}
</style>
