<template>
  <div class="story-card">
    <div class="story-cover">
      <img v-if="story.cover_image" :src="getCoverImageUrl(story.cover_image)" :alt="story.title" />
      <div v-else class="placeholder-cover">📚</div>
    </div>
    <div class="story-info">
      <h3>{{ story.title }}</h3>
      <span class="status" :class="story.status">{{ formatStatus(story.status) }}</span>
      <p class="summary">{{ story.summary || 'Chưa có tóm tắt' }}</p>
      <div class="actions">
        <button class="btn btn-small btn-edit" @click="$emit('edit', story.id)">
          Chỉnh sửa
        </button>
        <button class="btn btn-small btn-delete" @click="$emit('delete', story.id)">
          Xóa
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { Story } from '../stores/stories'

interface Props {
  story: Story
}

defineProps<Props>()

defineEmits<{
  edit: [storyId: string]
  delete: [storyId: string]
}>()

const API_BASE_URL = (import.meta as any).env?.VITE_API_URL || 'http://localhost:8080'

function getCoverImageUrl(path: string): string {
  if (path.startsWith('http')) {
    return path
  }
  return `${API_BASE_URL}${path}`
}

function formatStatus(status: string): string {
  return status === 'draft' ? 'Bản nháp' : 'Đã xuất bản'
}
</script>

<style scoped>
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

.actions {
  display: flex;
  gap: 0.5rem;
}

.btn {
  padding: 0.5rem 1rem;
  border-radius: 0.5rem;
  border: none;
  font-size: 0.875rem;
  cursor: pointer;
  transition: all 0.3s;
}

.btn-edit {
  background: #667eea;
  color: white;
}

.btn-edit:hover {
  background: #5a67d8;
}

.btn-delete {
  background: #fee2e2;
  color: #dc2626;
}

.btn-delete:hover {
  background: #fecaca;
}
</style>
