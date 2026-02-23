<template>
  <div
    class="scene-node"
    :class="{
      'is-start': data?.is_start,
      'is-end': data?.is_end,
    }"
  >
    <div class="scene-node-header">
      <span v-if="data?.is_start" class="badge start-badge">Bắt đầu</span>
      <span v-if="data?.is_end" class="badge end-badge">Kết thúc</span>
    </div>
    <div class="scene-node-content">
      <p class="content-preview">{{ contentPreview }}</p>
    </div>
    <Handle type="target" :position="Position.Top" class="handle" />
    <Handle type="source" :position="Position.Bottom" class="handle" />
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { Handle, Position } from '@vue-flow/core'
import type { SceneNodeData } from '@/stores/storyMap'

interface Props {
  id: string
  data?: SceneNodeData
}

const props = defineProps<Props>()

const contentPreview = computed(() => {
  if (!props.data?.content) return 'Chưa có nội dung'
  const maxLength = 100
  return props.data.content.length > maxLength
    ? props.data.content.slice(0, maxLength) + '...'
    : props.data.content
})
</script>

<style scoped>
.scene-node {
  background: white;
  border: 2px solid #e5e7eb;
  border-radius: 8px;
  padding: 12px;
  min-width: 180px;
  max-width: 240px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  transition: all 0.2s ease;
}

.scene-node:hover {
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.15);
}

.scene-node.is-start {
  border-color: #10b981;
  box-shadow: 0 0 0 2px rgba(16, 185, 129, 0.2);
}

.scene-node.is-end {
  border-color: #f59e0b;
}

.scene-node-header {
  display: flex;
  gap: 4px;
  margin-bottom: 8px;
}

.badge {
  font-size: 10px;
  padding: 2px 6px;
  border-radius: 4px;
  font-weight: 600;
  text-transform: uppercase;
}

.start-badge {
  background: #d1fae5;
  color: #065f46;
}

.end-badge {
  background: #fef3c7;
  color: #92400e;
}

.scene-node-content {
  font-size: 13px;
  line-height: 1.5;
  color: #374151;
}

.content-preview {
  margin: 0;
  white-space: pre-wrap;
  word-break: break-word;
}

.handle {
  width: 8px;
  height: 8px;
  background: #6b7280;
  border: 2px solid white;
}
</style>
