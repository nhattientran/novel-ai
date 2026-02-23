<template>
  <div class="story-map-page">
    <StoryMapHeader
      :title="storyTitle"
      @add-scene="addScene"
      @publish="publishStory"
    />

    <div class="map-container">
      <VueFlow
        v-model:nodes="nodes"
        v-model:edges="edges"
        :default-viewport="{ x: 0, y: 0, zoom: 1 }"
        :min-zoom="0.2"
        :max-zoom="4"
        @node-click="onNodeClick"
        @edge-click="onEdgeClick"
        @connect="onConnect"
      >
        <Background pattern-color="#aaa" :gap="16" />
        <Controls />
        <MiniMap />
      </VueFlow>
    </div>

    <SceneEditorModal
      v-if="selectedNode"
      :node="selectedNode"
      @close="closeModal"
      @save="saveNode"
      @delete="deleteNode"
    />
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRoute } from 'vue-router'
import { VueFlow } from '@vue-flow/core'
import { Background } from '@vue-flow/background'
import { Controls } from '@vue-flow/controls'
import { MiniMap } from '@vue-flow/minimap'
import type { Node, Edge, Connection, NodeMouseEvent, EdgeMouseEvent } from '@vue-flow/core'
import StoryMapHeader from '@/components/story-map/story-map-header.vue'
import SceneEditorModal from '@/components/scene-editor/scene-editor-modal.vue'

import '@vue-flow/core/dist/style.css'
import '@vue-flow/core/dist/theme-default.css'
import '@vue-flow/controls/dist/style.css'
import '@vue-flow/minimap/dist/style.css'

const route = useRoute()
const storyId = route.params.id as string
const storyTitle = ref('Tên truyện')

const nodes = ref<Node[]>([
  {
    id: '1',
    type: 'default',
    position: { x: 250, y: 100 },
    data: { label: 'Cảnh bắt đầu', content: '' },
  },
])

const edges = ref<Edge[]>([])
interface SceneNode extends Node {
  data: {
    label: string
    content: string
  }
}

const selectedNode = ref<SceneNode | null>(null)

const addScene = () => {
  const newNode: Node = {
    id: String(nodes.value.length + 1),
    type: 'default',
    position: { x: Math.random() * 400, y: Math.random() * 400 },
    data: { label: `Cảnh ${nodes.value.length + 1}`, content: '' },
  }
  nodes.value.push(newNode)
}

const onNodeClick = (event: NodeMouseEvent) => {
  selectedNode.value = event.node as SceneNode
}

const onEdgeClick = (event: EdgeMouseEvent) => {
  console.log('Edge clicked:', event.edge)
}

const onConnect = (connection: Connection) => {
  const newEdge: Edge = {
    id: `e${connection.source}-${connection.target}`,
    source: connection.source!,
    target: connection.target!,
    label: 'Lựa chọn',
  }
  edges.value.push(newEdge)
}

const closeModal = () => {
  selectedNode.value = null
}

const saveNode = (content: string) => {
  if (selectedNode.value) {
    selectedNode.value.data.content = content
  }
  closeModal()
}

const deleteNode = () => {
  if (selectedNode.value) {
    nodes.value = nodes.value.filter((n) => n.id !== selectedNode.value!.id)
    closeModal()
  }
}

const publishStory = () => {
  console.log('Publish story:', storyId)
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
  background: #f9fafb;
}
</style>
