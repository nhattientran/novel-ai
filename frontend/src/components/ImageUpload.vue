<template>
  <div class="image-upload">
    <div
      class="upload-area"
      :class="{ 'has-image': modelValue, 'is-dragging': isDragging }"
      @dragenter.prevent="isDragging = true"
      @dragleave.prevent="isDragging = false"
      @dragover.prevent
      @drop.prevent="handleDrop"
      @click="triggerFileInput"
    >
      <img v-if="modelValue" :src="getImageUrl(modelValue)" alt="Preview" class="preview" />
      <div v-else class="upload-placeholder">
        <span class="icon">📷</span>
        <p>Click hoặc kéo thả ảnh vào đây</p>
        <span class="hint">JPEG, PNG, WebP (tối đa 5MB)</span>
      </div>
      <input
        ref="fileInput"
        type="file"
        accept="image/jpeg,image/png,image/webp"
        style="display: none"
        @change="handleFileChange"
      />
    </div>

    <div v-if="isUploading" class="upload-progress">
      <div class="progress-bar">
        <div class="progress-fill" :style="{ width: uploadProgress + '%' }"></div>
      </div>
      <span>Đang tải lên...</span>
    </div>

    <div v-if="error" class="upload-error">
      {{ error }}
    </div>

    <button v-if="modelValue && !isUploading" type="button" class="btn-remove" @click="removeImage">
      Xóa ảnh
    </button>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { uploadsApi } from '../api/uploads'

interface Props {
  modelValue: string
}

defineProps<Props>()

const emit = defineEmits<{
  'update:modelValue': [value: string]
}>()

const API_BASE_URL = (import.meta as any).env?.VITE_API_URL || 'http://localhost:8080'

const fileInput = ref<HTMLInputElement | null>(null)
const isDragging = ref(false)
const isUploading = ref(false)
const uploadProgress = ref(0)
const error = ref<string | null>(null)

function getImageUrl(path: string): string {
  if (path.startsWith('http')) {
    return path
  }
  return `${API_BASE_URL}${path}`
}

function triggerFileInput() {
  fileInput.value?.click()
}

function handleDrop(event: DragEvent) {
  isDragging.value = false
  const files = event.dataTransfer?.files
  if (files && files.length > 0) {
    uploadFile(files[0])
  }
}

function handleFileChange(event: Event) {
  const target = event.target as HTMLInputElement
  const file = target.files?.[0]
  if (file) {
    uploadFile(file)
  }
}

async function uploadFile(file: File) {
  // Validate file size (5MB)
  const maxSize = 5 * 1024 * 1024
  if (file.size > maxSize) {
    error.value = 'Kích thước file không được vượt quá 5MB'
    return
  }

  // Validate file type
  const allowedTypes = ['image/jpeg', 'image/png', 'image/webp']
  if (!allowedTypes.includes(file.type)) {
    error.value = 'Chỉ chấp nhận file JPEG, PNG hoặc WebP'
    return
  }

  error.value = null
  isUploading.value = true
  uploadProgress.value = 0

  // Simulate progress
  const progressInterval = setInterval(() => {
    if (uploadProgress.value < 90) {
      uploadProgress.value += 10
    }
  }, 100)

  try {
    const response = await uploadsApi.uploadImage(file)
    uploadProgress.value = 100
    emit('update:modelValue', response.url)
  } catch (err: any) {
    error.value = err.message || 'Tải ảnh lên thất bại'
  } finally {
    clearInterval(progressInterval)
    isUploading.value = false
    // Reset file input
    if (fileInput.value) {
      fileInput.value.value = ''
    }
  }
}

function removeImage() {
  emit('update:modelValue', '')
  error.value = null
}
</script>

<style scoped>
.image-upload {
  width: 100%;
}

.upload-area {
  border: 2px dashed #d1d5db;
  border-radius: 0.75rem;
  padding: 2rem;
  text-align: center;
  cursor: pointer;
  transition: all 0.3s;
  min-height: 200px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.upload-area:hover {
  border-color: #667eea;
  background: #f9fafb;
}

.upload-area.is-dragging {
  border-color: #667eea;
  background: #eef2ff;
}

.upload-area.has-image {
  padding: 0;
  border-style: solid;
}

.preview {
  width: 100%;
  height: 200px;
  object-fit: cover;
  border-radius: 0.5rem;
}

.upload-placeholder {
  color: #6b7280;
}

.upload-placeholder .icon {
  font-size: 3rem;
  display: block;
  margin-bottom: 0.5rem;
}

.upload-placeholder p {
  margin: 0 0 0.25rem;
  font-weight: 500;
}

.upload-placeholder .hint {
  font-size: 0.875rem;
  color: #9ca3af;
}

.upload-progress {
  margin-top: 1rem;
  text-align: center;
}

.progress-bar {
  height: 8px;
  background: #e5e7eb;
  border-radius: 4px;
  overflow: hidden;
  margin-bottom: 0.5rem;
}

.progress-fill {
  height: 100%;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  transition: width 0.3s;
}

.upload-error {
  margin-top: 0.75rem;
  padding: 0.75rem;
  background: #fee2e2;
  color: #dc2626;
  border-radius: 0.5rem;
  font-size: 0.875rem;
}

.btn-remove {
  margin-top: 0.75rem;
  padding: 0.5rem 1rem;
  background: #fee2e2;
  color: #dc2626;
  border: none;
  border-radius: 0.5rem;
  cursor: pointer;
  font-size: 0.875rem;
  transition: background 0.3s;
}

.btn-remove:hover {
  background: #fecaca;
}
</style>
