<template>
  <div class="modal-overlay" @click.self="emit('close')">
    <div class="modal">
      <h2>Chỉnh sửa cảnh</h2>

      <div class="form-group">
        <label>Nội dung</label>
        <textarea v-model="content" rows="6" />
      </div>

      <div class="form-group">
        <label>Hình ảnh</label>
        <input type="file" accept="image/*" @change="handleImageUpload" />
      </div>

      <div class="modal-actions">
        <button class="btn" @click="emit('close')">Đóng</button>
        <button class="btn btn-danger" @click="emit('delete')">Xóa cảnh</button>
        <button class="btn btn-primary" @click="emit('save', content)">Lưu</button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'

const props = defineProps<{
  node: { data: { content: string } } | null
}>()

const emit = defineEmits<{
  close: []
  save: [content: string]
  delete: []
}>()

const content = ref('')

watch(
  () => props.node,
  (newNode) => {
    content.value = newNode?.data.content || ''
  },
  { immediate: true }
)

const handleImageUpload = (event: Event) => {
  const target = event.target as HTMLInputElement
  if (target.files && target.files[0]) {
    console.log('Upload image:', target.files[0])
  }
}
</script>

<style scoped>
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
  z-index: 1000;
}

.modal {
  background: white;
  padding: 2rem;
  border-radius: 1rem;
  width: 90%;
  max-width: 500px;
}

.modal h2 {
  margin-bottom: 1.5rem;
  color: #333;
}

.form-group {
  margin-bottom: 1.25rem;
}

.form-group label {
  display: block;
  margin-bottom: 0.5rem;
  color: #555;
  font-weight: 500;
}

.form-group textarea,
.form-group input {
  width: 100%;
  padding: 0.75rem;
  border: 1px solid #ddd;
  border-radius: 0.5rem;
  font-size: 1rem;
  font-family: inherit;
}

.form-group textarea:focus,
.form-group input:focus {
  outline: none;
  border-color: #667eea;
}

.modal-actions {
  display: flex;
  gap: 0.75rem;
  justify-content: flex-end;
  margin-top: 1.5rem;
}

.btn {
  padding: 0.5rem 1rem;
  border-radius: 0.5rem;
  border: 1px solid #ddd;
  background: white;
  font-size: 0.9rem;
  cursor: pointer;
  transition: all 0.3s;
}

.btn:hover {
  background: #f9fafb;
}

.btn-primary {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  border: none;
}

.btn-primary:hover {
  opacity: 0.9;
}

.btn-danger {
  background: #ef4444;
  color: white;
  border: none;
}

.btn-danger:hover {
  background: #dc2626;
}
</style>
