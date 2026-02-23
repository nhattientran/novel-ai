# Phase 04 — Reader: Discovery + Interactive Reading + Progress + Undo + Preload

**Date:** 2026-02-23
**Priority:** High
**Status:** ⬜ TODO
**Estimate:** L (3–7 days)
**Depends on:** Phase 03

---

## Overview
The reader-facing experience. Homepage lists published stories. Reader starts a story, reads scenes, makes choices, with progress saved locally and undo/history support. Preload adjacent scenes for smooth transitions.

## Key Insights
- Reader endpoints only serve `status="published"` stories.
- Scene transition query returns scene + choices + adjacent image URLs in one round-trip (performance).
- Progress stored in LocalStorage keyed by storyId (no server-side for MVP).
- History is a simple stack of user steps; cycles in graph are fine.
- Preload cap: max 3 adjacent scenes' images.

## Requirements

### Functional
- Homepage lists published stories (discovery).
- Reader can start a story (loads start scene).
- Reader sees scene content + choice buttons; clicking advances to next scene.
- Progress persisted to LocalStorage; resume on revisit.
- History panel shows past choices; undo goes back one step.

### Non-Functional
- Scene transition API < 200ms.
- Preload adjacent scene images for perceived instant transitions.
- Mobile-first responsive UI for reader mode.

## Architecture

```
Backend:
  internal/
    http/
      handlers/
        reader.go       # discovery, start, scene transition
    repo/
      reader_repo.go    # reader-specific Neo4j queries
    services/
      reader_service.go

Frontend:
  src/
    api/
      reader.ts         # API client for reader endpoints
    stores/
      reader.ts         # Pinia reader store (progress, history)
    pages/
      HomePage.vue          # Story discovery
      reader/
        ReadPage.vue        # Interactive reading UI
    components/
      reader/
        SceneDisplay.vue    # Scene content + image display
        ChoiceButtons.vue   # Choice button group
        HistoryPanel.vue    # Past choices list + undo
```

## Related Code Files

### Files to Create
- `backend/internal/repo/reader_repo.go`
- `backend/internal/services/reader_service.go`
- `backend/internal/http/handlers/reader.go`
- `frontend/src/api/reader.ts`
- `frontend/src/stores/reader.ts`
- `frontend/src/components/reader/SceneDisplay.vue`
- `frontend/src/components/reader/ChoiceButtons.vue`
- `frontend/src/components/reader/HistoryPanel.vue`

### Files to Modify
- `backend/internal/http/router.go` — add reader routes
- `frontend/src/pages/HomePage.vue` — implement story listing
- `frontend/src/pages/reader/ReadPage.vue` — implement reading UI
- `frontend/src/router/index.ts` — verify reader routes

## Implementation Steps

### 1. Reader API endpoints (M–L)

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/api/stories` | GET | List published stories |
| `/api/stories/:storyId/start` | GET | Get start scene |
| `/api/stories/:storyId/scenes/:sceneId` | GET | Get scene + choices + adjacent |

**Cypher — List published stories:**
```cypher
MATCH (s:Story {status:"published"})
RETURN s { .id, .title, .summary, .cover_image, .created_at } AS story
ORDER BY s.created_at DESC;
```

**Cypher — Get start scene:**
```cypher
MATCH (st:Story {id:$story_id, status:"published"})-[:STARTS_AT]->(sc:Scene)
RETURN sc { .id, .content, .image_url, .is_end } AS scene;
```

**Cypher — Get scene + choices + adjacent (single round-trip):**
```cypher
MATCH (st:Story {id:$story_id, status:"published"})-[:HAS_SCENE]->(sc:Scene {id:$scene_id})
OPTIONAL MATCH (sc)-[r:LEADS_TO]->(next:Scene)<-[:HAS_SCENE]-(st)
RETURN
  sc { .id, .content, .image_url, .is_end } AS scene,
  collect(distinct {
    to_scene_id: next.id,
    choice_text: r.choice_text,
    image_url: next.image_url
  }) AS choices;
```

**Scene response format:**
```json
{
  "scene": { "id": "...", "content": "...", "image_url": "...", "is_end": false },
  "choices": [
    { "to_scene_id": "...", "choice_text": "Open the door", "image_url": "..." }
  ]
}
```

### 2. Frontend: Story Discovery — HomePage (M)
- Fetch `GET /api/stories` on mount.
- Display story cards in a responsive grid (mobile-first).
- Each card: cover image, title, summary preview.
- Click → navigate to `/read/:storyId`.

### 3. Frontend: Interactive Reading — ReadPage (L)

**Pinia reader store (`stores/reader.ts`):**
```typescript
interface HistoryEntry {
  sceneId: string
  choiceText: string
  nextSceneId: string
}

interface ReaderState {
  currentStoryId: string | null
  currentScene: Scene | null
  choices: Choice[]
  history: HistoryEntry[]
}
```

**LocalStorage persistence:**
- Key: `progress:<storyId>`
- Value: `{ sceneId: string, history: HistoryEntry[] }`
- Save on every scene transition.
- Load on page mount if exists; otherwise call `/start`.

**Reading flow:**
1. Mount → check LocalStorage for saved progress.
2. If saved → load scene by saved `sceneId`.
3. If not → call `GET /stories/:id/start`.
4. Display scene content + image + choice buttons.
5. User clicks choice → push to history → load next scene → save to LocalStorage.
6. If `is_end` → show "The End" with options to restart or go home.

### 4. SceneDisplay component (M)
- Image banner (if `image_url` exists): full-width, aspect-ratio maintained.
- Content text: readable typography, appropriate line height, max-width for readability.
- Mobile-first: padding, font-size optimized for phone screens.

### 5. ChoiceButtons component (S–M)
- List of buttons at bottom of scene.
- Each button shows `choice_text`.
- Click emits event with `to_scene_id` and `choice_text`.
- If no choices + `is_end` → show end screen.

### 6. History/Undo — HistoryPanel (M)
- Expandable panel showing list of past decisions.
- Each entry: `choiceText` → clicked to navigate back.
- **Undo button**: pop last history entry → set `currentSceneId` to previous scene → save to LocalStorage.
- Undo multiple steps: pop multiple entries.

### 7. Preload strategy (M)
After each scene loads:
```typescript
// Preload adjacent scene images (max 3)
const preloadImages = (choices: Choice[]) => {
  choices.slice(0, 3).forEach(choice => {
    if (choice.image_url) {
      const img = new Image()
      img.src = choice.image_url
    }
  })
}
```
- Images only for MVP (scene content is small text, no need to prefetch).
- Cap at 3 to avoid excessive bandwidth on scenes with many choices.

### 8. Mobile-first responsive design (M)
- Reader UI: single column, full-width on mobile.
- Image: `width: 100%`, `max-height: 40vh`, `object-fit: cover`.
- Content: `padding: 1rem`, `font-size: 1.1rem`, `line-height: 1.8`.
- Choice buttons: stacked vertically, full-width, `min-height: 48px` for touch targets.
- History panel: collapsible drawer from bottom on mobile.

## Todo List
- [ ] Implement reader repository (list published, get start, get scene+choices)
- [ ] Implement reader service
- [ ] Implement reader handlers
- [ ] Update router with reader routes
- [ ] Implement frontend reader API client
- [ ] Implement Pinia reader store with LocalStorage persistence
- [ ] Build HomePage with story discovery listing
- [ ] Build ReadPage with reading flow
- [ ] Build SceneDisplay component (mobile-first)
- [ ] Build ChoiceButtons component
- [ ] Build HistoryPanel with undo
- [ ] Implement image preloading
- [ ] Implement LocalStorage save/restore
- [ ] Implement "The End" screen
- [ ] Test: start story → make choices → refresh → resume
- [ ] Test: undo goes back correctly
- [ ] Test: mobile viewport works well

## Success Criteria
- Reader can discover published stories on homepage.
- Reader can start a story and navigate through choices.
- Page reload continues from saved progress.
- Undo goes back at least one step.
- Images preload for next scenes.
- UI works well on mobile screens.

## Risk Assessment
| Risk | Impact | Mitigation |
|------|--------|------------|
| Cycles in story graph causing infinite history | Low | History is just a stack; cycles are fine, user can undo through them |
| Too many choices causing preload overload | Low | Cap preload to 3 adjacent scenes |
| LocalStorage quota exceeded | Very Low | Story progress is tiny (<1KB per story) |
| Scene transition > 200ms | Medium | Single Cypher query with indexes; monitor and optimize |

## Security Considerations
- Reader endpoints are public (no auth required for published stories).
- Only `published` stories served — enforce in every query.
- No user data exposed in reader mode (progress is client-side only).

## Next Steps
- → Phase 05: Publishing, Validation, Performance & Polish
