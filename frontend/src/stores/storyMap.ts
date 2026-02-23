import { ref, computed } from 'vue'
import { defineStore } from 'pinia'
import type { Node, Edge } from '@vue-flow/core'
import { graphApi, type GraphResponse, type GraphNode, type GraphEdge } from '@/api/graph'
import { scenesApi, type Scene, type CreateSceneRequest, type UpdateSceneRequest } from '@/api/scenes'
import { choicesApi, type CreateChoiceRequest, type UpdateChoiceRequest, type DeleteChoiceRequest } from '@/api/choices'

export interface SceneNodeData {
  content: string
  is_start: boolean
  is_end: boolean
  image_url?: string
}

export const useStoryMapStore = defineStore('storyMap', () => {
  // State
  const storyId = ref<string>('')
  const storyTitle = ref<string>('')
  const storyStatus = ref<string>('draft')
  const nodes = ref<Node<SceneNodeData>[]>([])
  const edges = ref<Edge[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)

  // Getters
  const startNodeId = computed(() => {
    const startNode = nodes.value.find(n => n.data?.is_start)
    return startNode?.id || null
  })

  const endNodes = computed(() => {
    return nodes.value.filter(n => n.data?.is_end)
  })

  // Actions
  async function loadGraph(id: string) {
    loading.value = true
    error.value = null
    storyId.value = id

    try {
      const response: GraphResponse = await graphApi.loadGraph(id)
      storyTitle.value = response.story.title
      storyStatus.value = response.story.status

      // Convert API nodes to Vue Flow nodes
      nodes.value = response.nodes.map((n: GraphNode) => ({
        id: n.id,
        type: 'scene',
        position: { x: n.position.x, y: n.position.y },
        data: {
          content: n.data.content,
          is_start: n.data.is_start,
          is_end: n.data.is_end,
          image_url: n.data.image_url,
        },
      }))

      // Convert API edges to Vue Flow edges
      edges.value = response.edges.map((e: GraphEdge) => ({
        id: e.id,
        source: e.source,
        target: e.target,
        label: e.label,
        type: 'default',
      }))
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to load graph'
      throw err
    } finally {
      loading.value = false
    }
  }

  async function createScene(data: CreateSceneRequest) {
    const scene: Scene = await scenesApi.create(storyId.value, data)

    const newNode: Node<SceneNodeData> = {
      id: scene.id,
      type: 'scene',
      position: { x: scene.pos_x, y: scene.pos_y },
      data: {
        content: scene.content,
        is_start: false,
        is_end: scene.is_end,
        image_url: scene.image_url,
      },
    }

    nodes.value.push(newNode)
    return scene
  }

  async function updateScene(sceneId: string, data: UpdateSceneRequest) {
    const scene: Scene = await scenesApi.update(storyId.value, sceneId, data)

    const nodeIndex = nodes.value.findIndex(n => n.id === sceneId)
    if (nodeIndex !== -1) {
      const node = nodes.value[nodeIndex]
      const nodeData = node.data
      if (!nodeData) return scene

      nodes.value[nodeIndex] = {
        ...node,
        position: {
          x: data.pos_x ?? node.position.x,
          y: data.pos_y ?? node.position.y,
        },
        data: {
          content: data.content ?? nodeData.content,
          is_start: nodeData.is_start,
          is_end: data.is_end ?? nodeData.is_end,
          image_url: data.image_url !== undefined
            ? (data.image_url ?? undefined)
            : nodeData.image_url,
        },
      }
    }

    return scene
  }

  async function deleteScene(sceneId: string) {
    await scenesApi.delete(storyId.value, sceneId)

    nodes.value = nodes.value.filter(n => n.id !== sceneId)
    edges.value = edges.value.filter(e => e.source !== sceneId && e.target !== sceneId)
  }

  async function setStartScene(sceneId: string) {
    await scenesApi.setStartScene(storyId.value, sceneId)

    // Update local state
    nodes.value = nodes.value.map(n => {
      const data = n.data!
      return {
        ...n,
        data: {
          content: data.content,
          is_start: n.id === sceneId,
          is_end: data.is_end,
          image_url: data.image_url,
        },
      }
    })
  }

  async function createChoice(data: CreateChoiceRequest) {
    const choice = await choicesApi.create(storyId.value, data)

    const newEdge: Edge = {
      id: `${choice.from_scene_id}->${choice.to_scene_id}`,
      source: choice.from_scene_id,
      target: choice.to_scene_id,
      label: choice.choice_text,
      type: 'default',
    }

    edges.value.push(newEdge)
    return choice
  }

  async function updateChoice(data: UpdateChoiceRequest) {
    const choice = await choicesApi.update(storyId.value, data)

    const edgeIndex = edges.value.findIndex(
      e => e.source === choice.from_scene_id && e.target === choice.to_scene_id
    )
    if (edgeIndex !== -1) {
      edges.value[edgeIndex] = {
        ...edges.value[edgeIndex],
        label: choice.choice_text,
      }
    }

    return choice
  }

  async function deleteChoice(data: DeleteChoiceRequest) {
    await choicesApi.delete(storyId.value, data)

    edges.value = edges.value.filter(
      e => !(e.source === data.from_scene_id && e.target === data.to_scene_id)
    )
  }

  function updateNodePosition(sceneId: string, x: number, y: number) {
    const nodeIndex = nodes.value.findIndex(n => n.id === sceneId)
    if (nodeIndex !== -1) {
      nodes.value[nodeIndex].position = { x, y }
    }
  }

  return {
    // State
    storyId,
    storyTitle,
    storyStatus,
    nodes,
    edges,
    loading,
    error,

    // Getters
    startNodeId,
    endNodes,

    // Actions
    loadGraph,
    createScene,
    updateScene,
    deleteScene,
    setStartScene,
    createChoice,
    updateChoice,
    deleteChoice,
    updateNodePosition,
  }
})
