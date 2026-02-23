# Phase 03 — Creator: Visual Story Map (Vue Flow) + Scenes + Choices

**Date:** 2026-02-23
**Priority:** High
**Status:** ⬜ TODO
**Estimate:** L (3–7 days)
**Depends on:** Phase 02

---

## Overview
The core creative feature. Creator can visually build branching story graphs using Vue Flow. Add/edit scenes (nodes), connect them with choices (edges), set start scene, and persist node positions.

## Key Insights
- Store `pos_x`, `pos_y` on Scene node in Neo4j (avoids separate layout model — YAGNI).
- Single `GET /graph` endpoint returns everything Vue Flow needs in one call.
- Debounce position saving on node drag to reduce API spam.
- Ensure LEADS_TO edges only connect scenes within the same story.

## Requirements

### Functional
- Scene CRUD (content, image, is_end flag, position).
- Choice (edge) CRUD (from_scene → to_scene with choice_text).
- Set exactly one start scene per story (STARTS_AT relationship).
- Load entire story graph for Vue Flow rendering.
- Scene editor panel for editing selected node content.

### Non-Functional
- Position updates debounced (500ms).
- Graph load should be fast for stories with up to ~100 scenes (MVP).

## Architecture

```
Backend:
  internal/
    domain/
      scene.go          # Scene struct
      choice.go         # Choice struct (edge data)
    repo/
      scene_repo.go     # Neo4j scene queries
      choice_repo.go    # Neo4j choice queries
      graph_repo.go     # Load full graph query
    services/
      scene_service.go
      choice_service.go
    http/
      handlers/
        scene.go
        choice.go
        graph.go

Frontend:
  src/
    api/
      scenes.ts
      choices.ts
      graph.ts
    stores/
      storyMap.ts       # Pinia store for Vue Flow state
    pages/
      creator/
        StoryMapPage.vue    # Vue Flow canvas + panels
    components/
      creator/
        SceneNode.vue       # Custom Vue Flow node
        SceneEditor.vue     # Right panel scene editor
        ChoiceModal.vue     # Modal to enter/edit choice text
```

## Related Code Files

### Files to Create
- `backend/internal/domain/scene.go`
- `backend/internal/domain/choice.go`
- `backend/internal/repo/scene_repo.go`
- `backend/internal/repo/choice_repo.go`
- `backend/internal/repo/graph_repo.go`
- `backend/internal/services/scene_service.go`
- `backend/internal/services/choice_service.go`
- `backend/internal/http/handlers/scene.go`
- `backend/internal/http/handlers/choice.go`
- `backend/internal/http/handlers/graph.go`
- `frontend/src/api/scenes.ts`
- `frontend/src/api/choices.ts`
- `frontend/src/api/graph.ts`
- `frontend/src/stores/storyMap.ts`
- `frontend/src/components/creator/SceneNode.vue`
- `frontend/src/components/creator/SceneEditor.vue`
- `frontend/src/components/creator/ChoiceModal.vue`

### Files to Modify
- `backend/internal/http/router.go` — add scene/choice/graph routes
- `frontend/src/pages/creator/StoryMapPage.vue` — implement Vue Flow
- `frontend/src/router/index.ts` — ensure StoryMapPage route exists

## Implementation Steps

### 1. Scene CRUD endpoints (M)

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/api/creator/stories/:storyId/scenes` | POST | Create scene |
| `/api/creator/stories/:storyId/scenes/:sceneId` | GET | Get scene detail |
| `/api/creator/stories/:storyId/scenes/:sceneId` | PATCH | Update scene |
| `/api/creator/stories/:storyId/scenes/:sceneId` | DELETE | Delete scene |

**Create scene request:**
```json
{ "content": "string", "image_url": "string|null", "is_end": false, "pos_x": 0, "pos_y": 0 }
```

**Cypher — Create scene:**
```cypher
MATCH (u:User {id:$user_id})-[:CREATED]->(st:Story {id:$story_id})
CREATE (sc:Scene {
  id: $scene_id, content: $content, image_url: $image_url,
  is_end: coalesce($is_end, false),
  pos_x: coalesce($pos_x, 0.0), pos_y: coalesce($pos_y, 0.0)
})
CREATE (st)-[:HAS_SCENE]->(sc)
RETURN sc { .id, .content, .image_url, .is_end, .pos_x, .pos_y } AS scene;
```

**Cypher — Update scene:**
```cypher
MATCH (u:User {id:$user_id})-[:CREATED]->(st:Story {id:$story_id})-[:HAS_SCENE]->(sc:Scene {id:$scene_id})
SET sc.content = coalesce($content, sc.content),
    sc.image_url = coalesce($image_url, sc.image_url),
    sc.is_end = coalesce($is_end, sc.is_end),
    sc.pos_x = coalesce($pos_x, sc.pos_x),
    sc.pos_y = coalesce($pos_y, sc.pos_y)
RETURN sc { .id, .content, .image_url, .is_end, .pos_x, .pos_y } AS scene;
```

**Cypher — Delete scene:**
```cypher
MATCH (u:User {id:$user_id})-[:CREATED]->(st:Story {id:$story_id})-[:HAS_SCENE]->(sc:Scene {id:$scene_id})
DETACH DELETE sc;
```

### 2. Choice (edge) management endpoints (M)

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/api/creator/stories/:storyId/choices` | POST | Create/upsert LEADS_TO edge |
| `/api/creator/stories/:storyId/choices` | PATCH | Update choice_text |
| `/api/creator/stories/:storyId/choices` | DELETE | Remove LEADS_TO edge |

**Request body (all three):**
```json
{ "from_scene_id": "string", "to_scene_id": "string", "choice_text": "string" }
```

**Cypher — Upsert choice:**
```cypher
MATCH (u:User {id:$user_id})-[:CREATED]->(st:Story {id:$story_id})
MATCH (st)-[:HAS_SCENE]->(a:Scene {id:$from_scene_id})
MATCH (st)-[:HAS_SCENE]->(b:Scene {id:$to_scene_id})
MERGE (a)-[r:LEADS_TO]->(b)
SET r.choice_text = $choice_text
RETURN a.id AS from, b.id AS to, r.choice_text AS choice_text;
```

**Cypher — Delete choice:**
```cypher
MATCH (u:User {id:$user_id})-[:CREATED]->(st:Story {id:$story_id})
MATCH (st)-[:HAS_SCENE]->(a:Scene {id:$from_scene_id})-[r:LEADS_TO]->(b:Scene {id:$to_scene_id})
DELETE r;
```

### 3. Start scene endpoint (S–M)

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/api/creator/stories/:storyId/start/:sceneId` | PUT | Set start scene |

**Cypher — Set start (ensure single STARTS_AT):**
```cypher
MATCH (u:User {id:$user_id})-[:CREATED]->(st:Story {id:$story_id})
OPTIONAL MATCH (st)-[old:STARTS_AT]->(:Scene)
DELETE old
WITH st
MATCH (st)-[:HAS_SCENE]->(sc:Scene {id:$scene_id})
MERGE (st)-[:STARTS_AT]->(sc)
RETURN sc.id AS start_scene_id;
```

### 4. Graph load endpoint (M)

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/api/creator/stories/:storyId/graph` | GET | Full graph for Vue Flow |

**Response format:**
```json
{
  "story": { "id": "...", "title": "...", "status": "draft" },
  "nodes": [
    { "id": "scene-1", "type": "scene", "position": {"x": 0, "y": 0},
      "data": { "content": "...", "is_start": true, "is_end": false } }
  ],
  "edges": [
    { "id": "scene-1->scene-2", "source": "scene-1", "target": "scene-2",
      "label": "Open the door" }
  ]
}
```

**Cypher — Load graph:**
```cypher
MATCH (u:User {id:$user_id})-[:CREATED]->(st:Story {id:$story_id})
OPTIONAL MATCH (st)-[:STARTS_AT]->(start:Scene)
MATCH (st)-[:HAS_SCENE]->(sc:Scene)
OPTIONAL MATCH (sc)-[r:LEADS_TO]->(to:Scene)<-[:HAS_SCENE]-(st)
RETURN
  st { .id, .title, .status } AS story,
  collect(distinct sc {
    .id, .pos_x, .pos_y, .is_end, .content,
    is_start: (start IS NOT NULL AND sc.id = start.id)
  }) AS scenes,
  collect(distinct {
    from: sc.id, to: to.id, choice_text: r.choice_text
  }) AS choices;
```

### 5. Vue Flow integration (L)

**StoryMapPage.vue — Core setup:**
```typescript
import { VueFlow, useVueFlow } from '@vue-flow/core'

const { nodes, edges, onConnect, onNodeDragStop } = useVueFlow()

// Load graph on mount
onMounted(async () => {
  const graph = await graphApi.loadGraph(storyId)
  nodes.value = graph.nodes.map(n => ({
    id: n.id, type: 'scene',
    position: { x: n.position.x, y: n.position.y },
    data: n.data
  }))
  edges.value = graph.edges.map(e => ({
    id: e.id, source: e.source, target: e.target, label: e.label
  }))
})
```

**Interactions:**
- **Add node**: button → create scene via API → add to `nodes` array
- **Drag node**: `onNodeDragStop` → debounce 500ms → PATCH scene `pos_x/pos_y`
- **Connect nodes**: `onConnect` → open ChoiceModal → enter `choice_text` → POST choice → add to `edges`
- **Edit edge**: click edge → ChoiceModal → PATCH choice
- **Delete node**: confirm dialog → DELETE scene API → remove from `nodes`
- **Delete edge**: DELETE choice API → remove from `edges`
- **Select node**: click → show SceneEditor panel on right side

**SceneNode.vue — Custom node:**
- Display truncated scene content preview.
- Badge for start/end status.
- Handle connections.

**SceneEditor.vue — Right panel:**
- Textarea for content (simple text, no rich editor for MVP).
- Image upload field (reuse ImageUpload component).
- Toggle: `is_end` checkbox.
- Button: "Set as Start Scene".
- Auto-save on blur or debounced input.

**ChoiceModal.vue:**
- Text input for `choice_text`.
- Confirm/cancel buttons.

### 6. Position debouncing
```typescript
const debouncedUpdatePosition = useDebounceFn(
  (sceneId: string, x: number, y: number) => {
    scenesApi.updateScene(storyId, sceneId, { pos_x: x, pos_y: y })
  },
  500
)

onNodeDragStop(({ node }) => {
  debouncedUpdatePosition(node.id, node.position.x, node.position.y)
})
```

## Todo List
- [ ] Implement Scene domain model
- [ ] Implement Choice domain model
- [ ] Implement scene repository (create, get, update, delete)
- [ ] Implement choice repository (upsert, delete)
- [ ] Implement graph repository (load full graph)
- [ ] Implement start scene repository
- [ ] Implement scene service
- [ ] Implement choice service
- [ ] Implement scene handlers
- [ ] Implement choice handlers
- [ ] Implement graph handler
- [ ] Update router with scene/choice/graph routes
- [ ] Implement frontend scenes API client
- [ ] Implement frontend choices API client
- [ ] Implement frontend graph API client
- [ ] Implement Pinia storyMap store
- [ ] Build StoryMapPage with Vue Flow
- [ ] Build SceneNode custom component
- [ ] Build SceneEditor panel
- [ ] Build ChoiceModal component
- [ ] Implement node drag position debouncing
- [ ] Implement add/delete node interactions
- [ ] Implement connect/disconnect edge interactions
- [ ] Implement set-as-start-scene interaction
- [ ] Test: create scenes → connect → set start → reload graph

## Success Criteria
- Map page renders story graph from Neo4j data.
- Creator can add scenes, connect them with choice text, reposition nodes.
- Start scene can be set; only one STARTS_AT exists per story.
- Graph loads correctly after page refresh.
- Position changes persist after drag.

## Risk Assessment
| Risk | Impact | Mitigation |
|------|--------|------------|
| Cross-story edge creation | High | Always match both scenes via `(st)-[:HAS_SCENE]->` |
| Multiple STARTS_AT relationships | Medium | Delete old before creating new in same transaction |
| Vue Flow complexity for juniors | Medium | Keep it simple: one custom node type, controlled state, minimal custom styling |
| Excessive API calls on drag | Low | Debounce 500ms |
| Large graphs slow to render | Low (MVP) | Acceptable for ≤100 scenes; address in Phase 2 if needed |

## Security Considerations
- All endpoints verify story ownership via `(u:User {id:$user_id})-[:CREATED]->`.
- Scene/choice operations scoped to the story they belong to.

## Next Steps
- → Phase 04: Reader Mode
