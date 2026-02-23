<template>
  <div class="story-list-page">
    <header class="header">
      <h1>Quản lý truyện</h1>
      <button class="btn btn-primary" @click="showCreateModal = true">
        + Tạo truyện mới
      </button>
    </header>

    <div v-if="isLoading" class="loading">
      Đang tải...
    </div>

    <div v-else-if="error" class="error">
      {{ error }}
      <button @click="fetchStories" class="btn btn-small">Thử lại</button>
    </div>

    <div v-else class="story-grid">
      <StoryCard
        v-for="story in stories"
        :key="story.id"
        :story="story"
        @edit="handleEdit"
        @delete="handleDelete"
      />
    </div>

    <p v-if="!isLoading && !error && stories.length === 0" class="empty-state">
      Bạn chưa có truyện nào. Hãy tạo truyện đầu tiên!
    </p>

    <!-- Create Story Modal -->
    <div v-if="showCreateModal" class="modal-overlay" @click.self="showCreateModal = false">
      <div class="modal">
        <h2>Tạo truyện mới</h2>
        <form @submit.prevent="handleCreate">
          <div class="form-group">
            <label>Tiêu đề *</label>
            <input
              v-model="createForm.title"
              type="text"
              required
              maxlength="200"
              placeholder="Nhập tiêu đề truyện"
            />
          </div>
          <div class="form-group">
            <label>Tóm tắt</label>
            <textarea
              v-model="createForm.summary"
              rows="4"
              maxlength="2000"
              placeholder="Nhập tóm tắt truyện (tùy chọn)"
            ></textarea>
          </div>
          <div class="form-actions">
            <button type="button" class="btn btn-secondary" @click="showCreateModal = false">
              Hủy
            </button>
            <button type="submit" class="btn btn-primary" :disabled="isCreating">
              {{ isCreating ? 'Đang tạo...' : 'Tạo truyện' }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useStoriesStore } from '../../stores/stories'
import StoryCard from '../../components/StoryCard.vue'

const router = useRouter()
const storiesStore = useStoriesStore()

const showCreateModal = ref(false)
const isCreating = ref(false)
const createForm = ref({
  title: '',
  summary: '',
})

const { stories, isLoading, error } = storiesStore

onMounted(() => {
  storiesStore.fetchStories()
})

async function handleCreate() {
  if (!createForm.value.title.trim()) return

  isCreating.value = true
  try {
    const newStory = await storiesStore.createStory({
      title: createForm.value.title,
      summary: createForm.value.summary,
    })
    showCreateModal.value = false
    createForm.value = { title: '', summary: '' }
    // Navigate to edit page
    router.push(`/creator/stories/${newStory.id}/edit`)
  } catch (err) {
    console.error('Failed to create story:', err)
  } finally {
    isCreating.value = false
  }
}

function handleEdit(storyId: string) {
  router.push(`/creator/stories/${storyId}/edit`)
}

async function handleDelete(storyId: string) {
  if (!confirm('Bạn có chắc muốn xóa truyện này? Hành động này không thể hoàn tác.')) {
    return
  }

  try {
    await storiesStore.deleteStory(storyId)
  } catch (err) {
    console.error('Failed to delete story:', err)
    alert('Không thể xóa truyện. Vui lòng thử lại.')
  }
}

async function fetchStories() {
  await storiesStore.fetchStories()
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

.btn-primary:hover:not(:disabled) {
  opacity: 0.9;
}

.btn-primary:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.btn-secondary {
  background: #e5e7eb;
  color: #374151;
}

.btn-secondary:hover {
  background: #d1d5db;
}

.btn-small {
  padding: 0.5rem 1rem;
  font-size: 0.875rem;
}

.story-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 1.5rem;
}

.loading, .error {
  text-align: center;
  padding: 4rem;
  color: #666;
}

.error {
  color: #dc2626;
}

.empty-state {
  text-align: center;
  padding: 4rem;
  color: #666;
  font-size: 1.1rem;
}

/* Modal Styles */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 100;
}

.modal {
  background: white;
  padding: 2rem;
  border-radius: 1rem;
  width: 100%;
  max-width: 500px;
  box-shadow: 0 20px 25px rgba(0, 0, 0, 0.15);
}

.modal h2 {
  margin-bottom: 1.5rem;
  color: #333;
}

.form-group {
  margin-bottom: 1rem;
}

.form-group label {
  display: block;
  margin-bottom: 0.5rem;
  font-weight: 500;
  color: #374151;
}

.form-group input,
.form-group textarea {
  width: 100%;
  padding: 0.75rem;
  border: 1px solid #d1d5db;
  border-radius: 0.5rem;
  font-size: 1rem;
  font-family: inherit;
}

.form-group input:focus,
.form-group textarea:focus {
  outline: none;
  border-color: #667eea;
  box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
}

.form-actions {
  display: flex;
  gap: 1rem;
  justify-content: flex-end;
  margin-top: 1.5rem;
}
</style>
