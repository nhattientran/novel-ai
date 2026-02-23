<template>
  <div class="modal-overlay" @click.self="emit('close')">
    <div class="modal">
      <h2>Lịch sử đọc</h2>
      <ul class="history-list">
        <li v-for="(item, index) in history" :key="index">
          {{ index + 1 }}. {{ item.sceneTitle }}
          <span v-if="item.choiceText">→ {{ item.choiceText }}</span>
        </li>
      </ul>
      <div class="modal-actions">
        <button class="btn btn-primary" @click="emit('close')">Đóng</button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
interface HistoryItem {
  sceneId: string
  sceneTitle?: string
  choiceText?: string
}

defineProps<{
  history: HistoryItem[]
}>()

const emit = defineEmits<{
  close: []
}>()
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
  max-height: 80vh;
  overflow-y: auto;
}

.modal h2 {
  margin-bottom: 1.5rem;
  color: #333;
}

.history-list {
  list-style: none;
  padding: 0;
}

.history-list li {
  padding: 0.75rem;
  border-bottom: 1px solid #e5e7eb;
  color: #444;
}

.history-list li:last-child {
  border-bottom: none;
}

.history-list span {
  color: #667eea;
  font-weight: 500;
}

.modal-actions {
  display: flex;
  justify-content: flex-end;
  margin-top: 1.5rem;
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
</style>
