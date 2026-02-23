<template>
  <header class="header">
    <router-link to="/" class="back-btn">← Quay lại</router-link>
    <h1>{{ storyTitle }}</h1>
    <button class="history-btn" @click="emit('show-history')">Lịch sử</button>
  </header>
</template>

<script setup lang="ts">
defineProps<{
  storyTitle: string
}>()

const emit = defineEmits<{
  'show-history': []
}>()
</script>

<style scoped>
.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1rem 1.5rem;
  background: white;
  border-bottom: 1px solid #e5e7eb;
  position: sticky;
  top: 0;
  z-index: 100;
}

.back-btn {
  color: #667eea;
  text-decoration: none;
  font-weight: 500;
}

.back-btn:hover {
  text-decoration: underline;
}

.header h1 {
  font-size: 1.1rem;
  color: #333;
  margin: 0;
}

.history-btn {
  padding: 0.5rem 1rem;
  border: 1px solid #ddd;
  background: white;
  border-radius: 0.5rem;
  cursor: pointer;
  font-size: 0.9rem;
}

.history-btn:hover {
  background: #f9fafb;
}
</style>
