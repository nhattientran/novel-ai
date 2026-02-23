<template>
  <div class="story-edit-page">
    <header class="header">
      <button class="btn btn-back" @click="goBack">
        ← Quay lại
      </button>
      <h1>{{ isNew ? 'Tạo truyện mới' : 'Chỉnh sửa truyện' }}</h1>
      <div class="header-actions">
        <button
          v-if="!isNew"
          class="btn btn-danger"
          @click="handleDelete"
        >
          Xóa
        </button>
        <button
          class="btn btn-primary"
          :disabled="isSaving || !canSave"
          @click="handleSave"
        >
          {{ isSaving ? 'Đang lưu...' : 'Lưu' }}
        </button>
      </div>
    </header>

    <div v-if="isLoading" class="loading">
      Đang tải...
    </div>

    <div v-else-if="error" class="error">
      {{ error }}
    </div>

    <div v-else class="edit-form">
      <div class="form-section">
        <h2>Thông tin cơ bản</h2>

        <div class="form-group">
          <label>Tiêu đề *</label>
          <input
            v-model="form.title"
            type="text"
            maxlength="200"
            placeholder="Nhập tiêu đề truyện"
          />
        </div>

        <div class="form-group">
          <label>Tóm tắt</label>
          <textarea
            v-model="form.summary"
            rows="6"
            maxlength="2000"
            placeholder="Nhập tóm tắt truyện..."
          ></textarea>
          <span class="char-count">{{ form.summary?.length || 0 }}/2000</span>
        </div>
      </div>

      <div class="form-section">
        <h2>Ảnh bìa</h2>
        <ImageUpload v-model="form.cover_image" />
      </div>

      <div v-if="!isNew" class="form-section">
        <h2>Trạng thái</h2>
        <div class="status-display">
          <span class="status-badge" :class="form.status">
            {{ form.status === 'draft' ? 'Bản nháp' : 'Đã xuất bản' }}
          </span>
          <p class="status-hint">
            {{ form.status === 'draft'
              ? 'Truyện đang ở chế độ bản nháp. Ngườii đọc chưa thể xem được.'
              : 'Truyện đã được xuất bản và có thể được đọc bởi ngườii dùng.'
            }}
          </p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useStoriesStore } from '../../stores/stories'
import ImageUpload from '../../components/ImageUpload.vue'

const route = useRoute()
const router = useRouter()
const storiesStore = useStoriesStore()

const storyId = computed(() => route.params.storyId as string)
const isNew = computed(() => storyId.value === 'new')

const isLoading = ref(false)
const isSaving = ref(false)
const error = ref<string | null>(null)

const form = ref({
  title: '',
  summary: '',
  cover_image: '',
  status: 'draft' as 'draft' | 'published',
})

const canSave = computed(() => form.value.title.trim().length > 0)

onMounted(async () => {
  if (!isNew.value) {
    await loadStory()
  }
})

async function loadStory() {
  isLoading.value = true
  error.value = null

  try {
    await storiesStore.fetchStory(storyId.value)
    const story = storiesStore.currentStory

    if (story) {
      form.value = {
        title: story.title,
        summary: story.summary,
        cover_image: story.cover_image,
        status: story.status,
      }
    } else {
      error.value = 'Không tìm thấy truyện'
    }
  } catch (err: any) {
    error.value = err.message || 'Không thể tải thông tin truyện'
  } finally {
    isLoading.value = false
  }
}

async function handleSave() {
  if (!canSave.value) return

  isSaving.value = true
  error.value = null

  try {
    const data = {
      title: form.value.title.trim(),
      summary: form.value.summary.trim(),
      cover_image: form.value.cover_image,
    }

    if (isNew.value) {
      const newStory = await storiesStore.createStory(data)
      router.replace(`/creator/stories/${newStory.id}/edit`)
    } else {
      await storiesStore.updateStory(storyId.value, data)
    }
  } catch (err: any) {
    error.value = err.message || 'Không thể lưu truyện'
  } finally {
    isSaving.value = false
  }
}

async function handleDelete() {
  if (!confirm('Bạn có chắc muốn xóa truyện này? Hành động này không thể hoàn tác.')) {
    return
  }

  try {
    await storiesStore.deleteStory(storyId.value)
    router.push('/creator/stories')
  } catch (err: any) {
    error.value = err.message || 'Không thể xóa truyện'
  }
}

function goBack() {
  router.push('/creator/stories')
}
</script>

<style scoped>
.story-edit-page {
  padding: 2rem;
  max-width: 800px;
  margin: 0 auto;
}

.header {
  display: flex;
  align-items: center;
  gap: 1rem;
  margin-bottom: 2rem;
}

.header h1 {
  flex: 1;
  font-size: 1.5rem;
  color: #333;
}

.header-actions {
  display: flex;
  gap: 0.75rem;
}

.btn {
  padding: 0.75rem 1.5rem;
  border-radius: 0.5rem;
  border: none;
  font-size: 1rem;
  cursor: pointer;
  transition: all 0.3s;
  text-decoration: none;
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
}

.btn-back {
  background: #f3f4f6;
  color: #374151;
}

.btn-back:hover {
  background: #e5e7eb;
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

.btn-danger {
  background: #fee2e2;
  color: #dc2626;
}

.btn-danger:hover {
  background: #fecaca;
}

.loading, .error {
  text-align: center;
  padding: 4rem;
  color: #666;
}

.error {
  color: #dc2626;
}

.edit-form {
  display: flex;
  flex-direction: column;
  gap: 2rem;
}

.form-section {
  background: white;
  padding: 1.5rem;
  border-radius: 1rem;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.form-section h2 {
  font-size: 1.25rem;
  color: #333;
  margin-bottom: 1.5rem;
}

.form-group {
  margin-bottom: 1.5rem;
}

.form-group:last-child {
  margin-bottom: 0;
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
  transition: border-color 0.3s, box-shadow 0.3s;
}

.form-group input:focus,
.form-group textarea:focus {
  outline: none;
  border-color: #667eea;
  box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
}

.form-group textarea {
  resize: vertical;
  min-height: 120px;
}

.char-count {
  display: block;
  text-align: right;
  font-size: 0.875rem;
  color: #9ca3af;
  margin-top: 0.25rem;
}

.status-display {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.status-badge {
  display: inline-flex;
  align-self: flex-start;
  padding: 0.5rem 1rem;
  border-radius: 1rem;
  font-size: 0.875rem;
  font-weight: 500;
  text-transform: uppercase;
}

.status-badge.draft {
  background: #fef3c7;
  color: #92400e;
}

.status-badge.published {
  background: #d1fae5;
  color: #065f46;
}

.status-hint {
  color: #6b7280;
  font-size: 0.875rem;
  margin: 0;
}
</style>
