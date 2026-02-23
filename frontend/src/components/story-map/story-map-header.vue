<template>
  <header class="header">
    <div class="header-left">
      <router-link to="/creator/stories" class="back-btn">← Quay lại</router-link>
      <h1>{{ title }}</h1>
    </div>
    <div class="header-actions">
      <button class="btn" @click="emit('add-scene')">+ Thêm cảnh</button>
      <button class="btn btn-primary" @click="emit('publish')">Xuất bản</button>
    </div>
  </header>
</template>

<script setup lang="ts">
defineProps<{
  title: string
}>()

const emit = defineEmits<{
  'add-scene': []
  publish: []
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
}

.header-left {
  display: flex;
  align-items: center;
  gap: 1rem;
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
  font-size: 1.25rem;
  color: #333;
}

.header-actions {
  display: flex;
  gap: 0.75rem;
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
</style>
