<template>
  <div class="modal-overlay" @click.self="$emit('cancel')">
    <div class="modal-content">
      <h3>{{ isEdit ? 'Chỉnh sửa lựa chọn' : 'Thêm lựa chọn' }}</h3>

      <div class="form-group">
        <label>Nội dung lựa chọn</label>
        <input
          v-model="choiceText"
          type="text"
          placeholder="Ví dụ: Mở cánh cửa"
          @keyup.enter="confirm"
        />
      </div>

      <div class="modal-actions">
        <button class="btn-cancel" @click="$emit('cancel')">Hủy</button>
        <button class="btn-confirm" @click="confirm" :disabled="!choiceText.trim()">
          {{ isEdit ? 'Lưu' : 'Thêm' }}
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'

interface Props {
  isEdit?: boolean
  initialText?: string
}

const props = withDefaults(defineProps<Props>(), {
  isEdit: false,
  initialText: '',
})

const emit = defineEmits<{
  confirm: [text: string]
  cancel: []
}>()

const choiceText = ref('')

watch(() => props.initialText, (text) => {
  choiceText.value = text
}, { immediate: true })

function confirm() {
  const trimmed = choiceText.value.trim()
  if (trimmed) {
    emit('confirm', trimmed)
  }
}
</script>

<style scoped>
.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal-content {
  background: white;
  padding: 24px;
  border-radius: 8px;
  width: 400px;
  max-width: 90vw;
}

.modal-content h3 {
  margin: 0 0 16px 0;
  font-size: 18px;
  font-weight: 600;
}

.form-group {
  margin-bottom: 20px;
}

.form-group label {
  display: block;
  font-size: 13px;
  font-weight: 500;
  color: #374151;
  margin-bottom: 6px;
}

.form-group input {
  width: 100%;
  padding: 10px 12px;
  border: 1px solid #d1d5db;
  border-radius: 6px;
  font-size: 14px;
}

.form-group input:focus {
  outline: none;
  border-color: #3b82f6;
  box-shadow: 0 0 0 2px rgba(59, 130, 246, 0.2);
}

.modal-actions {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

.btn-cancel {
  padding: 8px 16px;
  background: #f3f4f6;
  color: #374151;
  border: none;
  border-radius: 6px;
  font-size: 14px;
  cursor: pointer;
}

.btn-cancel:hover {
  background: #e5e7eb;
}

.btn-confirm {
  padding: 8px 16px;
  background: #3b82f6;
  color: white;
  border: none;
  border-radius: 6px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
}

.btn-confirm:hover:not(:disabled) {
  background: #2563eb;
}

.btn-confirm:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}
</style>
