<template>
  <div class="story-list-page">
    <header class="header">
      <h1>Quản lý truyện</h1>
      <button class="btn btn-primary" @click="createStory">
        + Tạo truyện mới
      </button>
    </header>

    <div class="story-grid">
      <div v-for="story in stories" :key="story.id" class="story-card">
        <div class="story-cover">
          <img v-if="story.coverImage" :src="story.coverImage" :alt="story.title" />
          <div v-else class="placeholder-cover">📚</div>
        </div>
        <div class="story-info">
          <h3>{{ story.title }}</h3>
          <p class="status" :class="story.status">{{ story.status }}</p>
          <p class="summary">{{ story.summary }}</p>
          <div class="actions">
            <router-link :to="`/creator/stories/${story.id}/map`" class="btn btn-small">
              Chỉnh sửa
            </router-link>
          </div>
        </div>
      </div>
    </div>

    <p v-if="stories.length === 0" class="empty-state">
      Bạn chưa có truyện nào. Hãy tạo truyện đầu tiên!
    </p>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

interface Story {
  id: string
  title: string
  summary: string
  coverImage?: string
  status: 'draft' | 'published'
}

const stories = ref<Story[]>([])

const createStory = () => {
  // TODO: Implement create story modal/navigation
  console.log('Create new story')
}
</script>

<style scoped>
.story-list-page {
  padding: 2rem;
  max-width: 1200px;
  margin: 0 auto;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;
}

.header h1 {
  font-size: 1.75rem;
  color: #333;
}

.btn {
  padding: 0.75rem 1.5rem;
  border-radius: 0.5rem;
  border: none;
  font-size: 1rem;
  cursor: pointer;
  transition: all 0.3s;
  text-decoration: none;
  display: inline-block;
}

.btn-primary {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
}

.btn-primary:hover {
  opacity: 0.9;
}

.btn-small {
  padding: 0.5rem 1rem;
  font-size: 0.875rem;
  background: #667eea;
  color: white;
}

.story-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 1.5rem;
}

.story-card {
  background: white;
  border-radius: 1rem;
  overflow: hidden;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  transition: transform 0.3s, box-shadow 0.3s;
}

.story-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 10px 20px rgba(0, 0, 0, 0.15);
}

.story-cover {
  height: 180px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  display: flex;
  align-items: center;
  justify-content: center;
}

.story-cover img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.placeholder-cover {
  font-size: 4rem;
}

.story-info {
  padding: 1.5rem;
}

.story-info h3 {
  font-size: 1.25rem;
  margin-bottom: 0.5rem;
  color: #333;
}

.status {
  display: inline-block;
  padding: 0.25rem 0.75rem;
  border-radius: 1rem;
  font-size: 0.75rem;
  font-weight: 500;
  text-transform: uppercase;
  margin-bottom: 0.75rem;
}

.status.draft {
  background: #fef3c7;
  color: #92400e;
}

.status.published {
  background: #d1fae5;
  color: #065f46;
}

.summary {
  color: #666;
  font-size: 0.9rem;
  line-height: 1.5;
  margin-bottom: 1rem;
  display: -webkit-box;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.empty-state {
  text-align: center;
  padding: 4rem;
  color: #666;
  font-size: 1.1rem;
}
</style>
