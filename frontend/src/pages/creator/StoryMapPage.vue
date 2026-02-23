<template>
  <div class="story-map-page">
    <StoryMapHeader
      :title="store.storyTitle"
      @add-scene="addScene"
      @publish="publishStory"
    />

    <div class="map-container">
      <VueFlow
        v-model:nodes="store.nodes"
        v-model:edges="store.edges"
        :default-viewport="{ x: 0, y: 0, zoom: 1 }"
        :min-zoom="0.2"
        :max-zoom="4"
        :node-types="nodeTypes"
        @node-click="onNodeClick"
        @edge-click="onEdgeClick"
        @connect="onConnect"
        @node-drag-stop="onNodeDragStop"
      >
        <Background pattern-color="#aaa" :gap="16" />
        <Controls />
        <MiniMap />
      </VueFlow>

      <SceneEditor
        v-if="selectedNode"
        :node="selectedNode"
        @close="closeEditor"
        @update-content="updateNodeContent"
        @update-is-end="updateNodeIsEnd"
        @update-image="updateNodeImage"
        @set-start="setStartScene"
        @delete="deleteSelectedNode"
      />
    </div>

    <ChoiceModal
      v-if="showChoiceModal"
      :is-edit="isEditingChoice"
      :initial-text="editingChoiceText"
      @confirm="onChoiceConfirm"
      @cancel="closeChoiceModal"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, markRaw } from 'vue'
import { useRoute } from 'vue-router'
import { VueFlow, useVueFlow } from '@vue-flow/core'
import { Background } from '@vue-flow/background'
import { Controls } from '@vue-flow/controls'
import { MiniMap } from '@vue-flow/minimap'
import type { Node, Edge, Connection, NodeMouseEvent, EdgeMouseEvent } from '@vue-flow/core'
import StoryMapHeader from '@/components/story-map/story-map-header.vue'
import SceneNode from '@/components/creator/SceneNode.vue'
import SceneEditor from '@/components/creator/SceneEditor.vue'
import ChoiceModal from '@/components/creator/ChoiceModal.vue'
import { useStoryMapStore, type SceneNodeData } from '@/stores/storyMap'
import { useDebounceFn } from '@/composables/useDebounce'

import '@vue-flow/core/dist/style.css'
import '@vue-flow/core/dist/theme-default.css'
import '@vue-flow/controls/dist/style.css'
import '@vue-flow/minimap/dist/style.css'

const route = useRoute()
const store = useStoryMapStore()
const { fitView } = useVueFlow()

// Node types registration
const nodeTypes: Record<string, any> = {
  scene: markRaw(SceneNode),
}

// State
const storyId = route.params.id as string
const selectedNode = ref<Node<SceneNodeData> | null>(null)
const showChoiceModal = ref(false)
const isEditingChoice = ref(false)
const editingChoiceText = ref('')
const pendingConnection = ref<Connection | null>(null)
const editingEdge = ref<Edge | null>(null)

// Debounced position update
const debouncedUpdatePosition = useDebounceFn(
  async (sceneId: string, x: number, y: number) => {
    try {
      await store.updateScene(sceneId, { pos_x: x, pos_y: y })
    } catch (err) {
      console.error('Failed to update position:', err)
    }
  },
  500
)

// Load graph on mount
onMounted(async () => {
  try {
    await store.loadGraph(storyId)
    setTimeout(() => fitView({ padding: 0.2 }), 100)
  } catch (err) {
    console.error('Failed to load graph:', err)
    alert('Không thể tải story map')
  }
})

// Node interactions
function onNodeClick(event: NodeMouseEvent) {
  selectedNode.value = event.node as Node<SceneNodeData>
}

function closeEditor() {
  selectedNode.value = null
}

async function addScene() {
  const centerX = 250 + Math.random() * 100
  const centerY = 150 + Math.random() * 100

  try {
    await store.createScene({
      content: 'Cảnh mới',
      is_end: false,
      pos_x: centerX,
      pos_y: centerY,
    })
  } catch (err) {
    console.error('Failed to create scene:', err)
    alert('Không thể tạo cảnh mới')
  }
}

async function updateNodeContent(content: string) {
  if (!selectedNode.value) return

  try {
    await store.updateScene(selectedNode.value.id, { content })
  } catch (err) {
    console.error('Failed to update content:', err)
    alert('Không thể cập nhật nội dung')
  }
}

async function updateNodeIsEnd(isEnd: boolean) {
  if (!selectedNode.value) return

  try {
    await store.updateScene(selectedNode.value.id, { is_end: isEnd })
  } catch (err) {
    console.error('Failed to update is_end:', err)
    alert('Không thể cập nhật trạng thái')
  }
}

async function updateNodeImage(imageUrl: string | null) {
  if (!selectedNode.value) return

  try {
    await store.updateScene(selectedNode.value.id, { image_url: imageUrl })
  } catch (err) {
    console.error('Failed to update image:', err)
    alert('Không thể cập nhật hình ảnh')
  }
}

async function setStartScene() {
  if (!selectedNode.value) return

  try {
    await store.setStartScene(selectedNode.value.id)
  } catch (err) {
    console.error('Failed to set start scene:', err)
    alert('Không thể đặt cảnh bắt đầu')
  }
}

async function deleteSelectedNode() {
  if (!selectedNode.value) return

  if (!confirm('Bạn có chắc muốn xóa cảnh này?')) return

  try {
    await store.deleteScene(selectedNode.value.id)
    closeEditor()
  } catch (err) {
    console.error('Failed to delete scene:', err)
    alert('Không thể xóa cảnh')
  }
}

// Edge interactions
function onEdgeClick(event: EdgeMouseEvent) {
  editingEdge.value = event.edge
  editingChoiceText.value = event.edge.label as string || ''
  isEditingChoice.value = true
  showChoiceModal.value = true
}

function onConnect(connection: Connection) {
  pendingConnection.value = connection
  isEditingChoice.value = false
  editingChoiceText.value = ''
  showChoiceModal.value = true
}

async function onChoiceConfirm(text: string) {
  if (isEditingChoice.value && editingEdge.value) {
    // Update existing choice
    try {
      await store.updateChoice({
        from_scene_id: editingEdge.value.source,
        to_scene_id: editingEdge.value.target,
        choice_text: text,
      })
    } catch (err) {
      console.error('Failed to update choice:', err)
      alert('Không thể cập nhật lựa chọn')
    }
  } else if (pendingConnection.value) {
    // Create new choice
    try {
      await store.createChoice({
        from_scene_id: pendingConnection.value.source!,
        to_scene_id: pendingConnection.value.target!,
        choice_text: text,
      })
    } catch (err) {
      console.error('Failed to create choice:', err)
      alert('Không thể tạo lựa chọn')
    }
  }

  closeChoiceModal()
}

function closeChoiceModal() {
  showChoiceModal.value = false
  pendingConnection.value = null
  editingEdge.value = null
  editingChoiceText.value = ''
}

// Node drag
function onNodeDragStop({ node }: { node: Node }) {
  debouncedUpdatePosition(node.id, node.position.x, node.position.y)
}

function publishStory() {
  console.log('Publish story:', storyId)
  alert('Tính năng publish sẽ được triển khai ở phase sau')
}
</script>

<style scoped>
.story-map-page {
  height: 100vh;
  display: flex;
  flex-direction: column;
}

.map-container {
  flex: 1;
  display: flex;
  background: #f9fafb;
  overflow: hidden;
}

.map-container :deep(.vue-flow) {
  flex: 1;
}
</style>
