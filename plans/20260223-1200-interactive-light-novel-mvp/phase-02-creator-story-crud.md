# Phase 02 — Creator: Story Management (CRUD + Draft/Published)

**Date:** 2026-02-23
**Priority:** High
**Status:** ✅ COMPLETED
**Estimate:** M–L (2–5 days)
**Depends on:** Phase 01

---

## Overview
Creator can create, edit, delete stories with title, cover image, summary. Stories default to `draft` status. Image upload via simple filesystem storage.

## Key Insights
- `DETACH DELETE` for story deletion in MVP — acceptable for now.
- Cover image stored on local filesystem; serve static files via Gin.
- Keep upload logic simple: multipart → disk → return URL.

## Requirements

### Functional
- CRUD operations for stories (owned by authenticated creator).
- Image upload endpoint for cover images.
- Stories default to `draft` status on creation.
- Creator can only see/edit their own stories.

### Non-Functional
- File size limit for uploads (e.g., 5MB).
- Restrict mime types to images (jpeg, png, webp).

## Architecture

```
Backend:
  internal/
    domain/
      story.go          # Story struct
    repo/
      story_repo.go     # Neo4j story queries
    services/
      story_service.go  # business logic
    http/
      handlers/
        story.go        # CRUD handlers
    storage/
      local.go          # filesystem image storage
  uploads/              # stored images directory

Frontend:
  src/
    api/
      stories.ts        # API client for story endpoints
      uploads.ts        # API client for file uploads
    stores/
      stories.ts        # Pinia story store
    pages/
      creator/
        StoryListPage.vue   # list + create
        StoryEditPage.vue   # edit title/summary/cover
    components/
      StoryCard.vue         # story card component
      ImageUpload.vue       # upload component
```

## Related Code Files

### Files to Create
- `backend/internal/domain/story.go`
- `backend/internal/repo/story_repo.go`
- `backend/internal/services/story_service.go`
- `backend/internal/http/handlers/story.go`
- `backend/internal/storage/local.go`
- `backend/migrations/002_story_constraints.cypher`
- `frontend/src/api/stories.ts`
- `frontend/src/api/uploads.ts`
- `frontend/src/stores/stories.ts`
- `frontend/src/pages/creator/StoryEditPage.vue`
- `frontend/src/components/StoryCard.vue`
- `frontend/src/components/ImageUpload.vue`

### Files to Modify
- `backend/internal/http/router.go` — add story + upload routes
- `frontend/src/pages/creator/StoryListPage.vue` — implement UI
- `frontend/src/router/index.ts` — add StoryEditPage route

## Implementation Steps

### 1. Neo4j constraints + indexes (S)
`backend/migrations/002_story_constraints.cypher`:
```cypher
CREATE CONSTRAINT story_id_unique IF NOT EXISTS FOR (s:Story) REQUIRE s.id IS UNIQUE;
CREATE CONSTRAINT scene_id_unique IF NOT EXISTS FOR (s:Scene) REQUIRE s.id IS UNIQUE;
CREATE INDEX story_status_idx IF NOT EXISTS FOR (s:Story) ON (s.status);
CREATE INDEX story_created_at_idx IF NOT EXISTS FOR (s:Story) ON (s.created_at);
```

### 2. Story domain model (S)
```go
type Story struct {
    ID         string    `json:"id"`
    Title      string    `json:"title"`
    Summary    string    `json:"summary"`
    CoverImage string    `json:"cover_image"`
    Status     string    `json:"status"` // "draft" or "published"
    CreatedAt  time.Time `json:"created_at"`
}
```

### 3. Image upload storage (M)
`backend/internal/storage/local.go`:
- Accept multipart file.
- Validate: size ≤ 5MB, mime type in [image/jpeg, image/png, image/webp].
- Generate unique filename (UUID + extension).
- Save to `backend/uploads/`.
- Return public URL path `/uploads/<filename>`.

`POST /api/uploads/images`:
- Requires auth.
- Returns `{ "url": "/uploads/abc123.jpg" }`.

Gin static serving: `router.Static("/uploads", "./uploads")`

### 4. Story repository (M)
`backend/internal/repo/story_repo.go`:

**Create story:**
```cypher
MATCH (u:User {id: $user_id})
CREATE (s:Story {
  id: $story_id, title: $title, summary: $summary,
  cover_image: $cover_image, status: "draft", created_at: datetime()
})
CREATE (u)-[:CREATED]->(s)
RETURN s { .id, .title, .summary, .cover_image, .status, .created_at } AS story;
```

**List my stories:**
```cypher
MATCH (u:User {id: $user_id})-[:CREATED]->(s:Story)
RETURN s { .id, .title, .summary, .cover_image, .status, .created_at } AS story
ORDER BY s.created_at DESC;
```

**Get story by ID (owned):**
```cypher
MATCH (u:User {id: $user_id})-[:CREATED]->(s:Story {id: $story_id})
RETURN s { .id, .title, .summary, .cover_image, .status, .created_at } AS story;
```

**Update story (partial):**
```cypher
MATCH (u:User {id: $user_id})-[:CREATED]->(s:Story {id: $story_id})
SET s.title = coalesce($title, s.title),
    s.summary = coalesce($summary, s.summary),
    s.cover_image = coalesce($cover_image, s.cover_image)
RETURN s { .id, .title, .summary, .cover_image, .status, .created_at } AS story;
```

**Delete story (and all related):**
```cypher
MATCH (u:User {id: $user_id})-[:CREATED]->(s:Story {id: $story_id})
OPTIONAL MATCH (s)-[:HAS_SCENE]->(sc:Scene)
DETACH DELETE sc, s;
```

### 5. Story CRUD endpoints (M)
| Endpoint | Method | Description |
|----------|--------|-------------|
| `/api/creator/stories` | POST | Create story |
| `/api/creator/stories` | GET | List my stories |
| `/api/creator/stories/:storyId` | GET | Get story detail |
| `/api/creator/stories/:storyId` | PATCH | Update story |
| `/api/creator/stories/:storyId` | DELETE | Delete story + scenes |

**Create request:**
```json
{ "title": "string", "summary": "string", "cover_image": "string|null" }
```

**Create response:**
```json
{ "id":"...", "title":"...", "summary":"...", "cover_image":"...", "status":"draft", "created_at":"..." }
```

### 6. Frontend: Story list + editor (M)
- **StoryListPage**: fetch stories, display cards, "Create Story" button.
- **StoryEditPage**: form with title, summary, cover image upload.
- **StoryCard**: thumbnail + title + status badge.
- **ImageUpload**: file input → upload → show preview.

## Todo List
- [x] Create Neo4j story/scene constraints migration
- [x] Implement Story domain model
- [x] Implement local file storage for images
- [x] Implement upload endpoint with validation
- [x] Implement story repository (create, list, get, update, delete)
- [x] Implement story service
- [x] Implement story CRUD handlers
- [x] Update router with story + upload routes
- [x] Implement frontend stories API client
- [x] Implement frontend uploads API client
- [x] Implement Pinia stories store
- [x] Build StoryListPage UI
- [x] Build StoryEditPage UI
- [x] Build StoryCard component
- [x] Build ImageUpload component
- [x] Add StoryEditPage route
- [x] Test: create → edit → delete story flow

## Success Criteria
- Creator can create a story with title, summary, and cover image.
- Creator can list, edit, and delete their own stories.
- Stories are always created with `draft` status.
- Image upload works and returns accessible URL.

## Risk Assessment
| Risk | Impact | Mitigation |
|------|--------|------------|
| Orphaned Scene nodes after story deletion | Medium | Use `DETACH DELETE` on both story and its scenes |
| Large file uploads | Low | Enforce 5MB size limit |
| Unauthorized access to other users' stories | High | Always match `(u:User {id:$user_id})-[:CREATED]->` |

## Security Considerations
- All story endpoints require authentication.
- Story operations scoped to authenticated user via CREATED relationship.
- File upload validates mime type and size.

## Next Steps
- → Phase 03: Visual Story Map + Scenes + Choices
