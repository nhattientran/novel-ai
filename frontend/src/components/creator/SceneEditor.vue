<template>
  <div class="scene-editor">
    <div class="editor-header">
      <h3>Chỉnh sửa cảnh</h3>
      <button class="close-btn" @click="$emit('close')">
        <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <line x1="18" y1="6" x2="6" y2="18"></line>
          <line x1="6" y1="6" x2="18" y2="18"></line>
        </svg>
      </button>
    </div>

    <div class="editor-body">
      <div class="form-group">
        <label>Nội dung</label>
        <textarea
          v-model="form.content"
          rows="8"
          placeholder="Nhập nội dung cảnh..."
          @blur="saveContent"
        ></textarea>
      </div>

      <div class="form-group">
        <label>Hình ảnh</label>
        <div class="image-upload">
          <input
            ref="fileInput"
            type="file"
            accept="image/*"
            style="display: none"
            @change="handleImageUpload"
          />
          <button
            v-if="!form.image_url"
            class="upload-btn"
            @click="openFilePicker"
            :disabled="uploading"
          >
            {{ uploading ? 'Đang tải...' : 'Tải ảnh lên' }}
          </button>
          <div v-else class="image-preview">
            <img :src="form.image_url" alt="Scene image" />
            <button class="remove-btn" @click="removeImage">Xóa</button>
          </div>
        </div>
      </div>

      <div class="form-group checkbox-group">
        <label class="checkbox-label">
          <input
            type="checkbox"
            v-model="form.is_end"
            @change="saveIsEnd"
          />
          <span>Cảnh kết thúc</span>
        </label>
      </div>

      <div class="actions">
        <button
          v-if="!node?.data?.is_start"
          class="btn-set-start"
          @click="$emit('set-start')"
        >
          Đặt làm cảnh bắt đầu
        </button>
        <span v-else class="start-indicator">
          Cảnh bắt đầu
        </span>

        <button class="btn-delete" @click="$emit('delete')">
          Xóa cảnh
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import type { Node } from '@vue-flow/core'
import type { SceneNodeData } from '@/stores/storyMap'
import { uploadsApi } from '@/api/uploads'

interface Props {
  node: Node<SceneNodeData> | null
}

const props = defineProps<Props>()

const emit = defineEmits<{
  close: []
  'update-content': [content: string]
  'update-is-end': [isEnd: boolean]
  'update-image': [imageUrl: string | null]
  'set-start': []
  delete: []
}>()

const form = ref({
  content: '',
  image_url: '',
  is_end: false,
})

const uploading = ref(false)
const fileInput = ref<HTMLInputElement | null>(null)

watch(() => props.node, (node) => {
  if (node?.data) {
    form.value.content = node.data.content || ''
    form.value.image_url = node.data.image_url || ''
    form.value.is_end = node.data.is_end || false
  }
}, { immediate: true })

function saveContent() {
  emit('update-content', form.value.content)
}

function saveIsEnd() {
  emit('update-is-end', form.value.is_end)
}

async function handleImageUpload(event: Event) {
  const target = event.target as HTMLInputElement
  const file = target.files?.[0]
  if (!file) return

  uploading.value = true
  try {
    const result = await uploadsApi.uploadImage(file)
    form.value.image_url = result.url
    emit('update-image', result.url)
  } catch (err) {
    console.error('Upload failed:', err)
    alert('Tải ảnh lên thất bại')
  } finally {
    uploading.value = false
    if (fileInput.value) {
      fileInput.value.value = ''
    }
  }
}

function removeImage() {
  form.value.image_url = ''
  emit('update-image', null)
}

function openFilePicker() {
  fileInput.value?.click()
}
</script>

<style scoped>
.scene-editor {
  width: 320px;
  background: white;
  border-left: 1px solid #e5e7eb;
  height: 100%;
  display: flex;
  flex-direction: column;
}

.editor-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px;
  border-bottom: 1px solid #e5e7eb;
}

.editor-header h3 {
  margin: 0;
  font-size: 16px;
  font-weight: 600;
}

.close-btn {
  background: none;
  border: none;
  cursor: pointer;
  padding: 4px;
  color: #6b7280;
  border-radius: 4px;
}

.close-btn:hover {
  background: #f3f4f6;
  color: #111827;
}

.editor-body {
  flex: 1;
  padding: 16px;
  overflow-y: auto;
}

.form-group {
  margin-bottom: 16px;
}

.form-group label {
  display: block;
  font-size: 13px;
  font-weight: 500;
  color: #374151;
  margin-bottom: 6px;
}

.form-group textarea {
  width: 100%;
  padding: 10px 12px;
  border: 1px solid #d1d5db;
  border-radius: 6px;
  font-size: 14px;
  resize: vertical;
  font-family: inherit;
}

.form-group textarea:focus {
  outline: none;
  border-color: #3b82f6;
  box-shadow: 0 0 0 2px rgba(59, 130, 246, 0.2);
}

.checkbox-group {
  margin-top: 20px;
}

.checkbox-label {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  font-size: 14px;
}

.checkbox-label input {
  width: 16px;
  height: 16px;
}

.upload-btn {
  width: 100%;
  padding: 10px;
  border: 2px dashed #d1d5db;
  border-radius: 6px;
  background: #f9fafb;
  color: #6b7280;
  cursor: pointer;
  font-size: 13px;
}

.upload-btn:hover:not(:disabled) {
  border-color: #3b82f6;
  color: #3b82f6;
}

.upload-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.image-preview {
  position: relative;
}

.image-preview img {
  width: 100%;
  border-radius: 6px;
  max-height: 150px;
  object-fit: cover;
}

.image-preview .remove-btn {
  position: absolute;
  top: 8px;
  right: 8px;
  padding: 4px 8px;
  background: rgba(0, 0, 0, 0.7);
  color: white;
  border: none;
  border-radius: 4px;
  font-size: 12px;
  cursor: pointer;
}

.actions {
  margin-top: 24px;
  padding-top: 16px;
  border-top: 1px solid #e5e7eb;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.btn-set-start {
  padding: 10px 16px;
  background: #10b981;
  color: white;
  border: none;
  border-radius: 6px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
}

.btn-set-start:hover {
  background: #059669;
}

.start-indicator {
  padding: 10px 16px;
  background: #d1fae5;
  color: #065f46;
  border-radius: 6px;
  font-size: 14px;
  font-weight: 500;
  text-align: center;
}

.btn-delete {
  padding: 10px 16px;
  background: #fee2e2;
  color: #dc2626;
  border: none;
  border-radius: 6px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
}

.btn-delete:hover {
  background: #fecaca;
}
</style>
