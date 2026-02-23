# Phase 05 — Publishing, Validation, Performance & Polish

**Date:** 2026-02-23
**Priority:** Medium
**Status:** ⬜ TODO
**Estimate:** M (1–3 days)
**Depends on:** All previous phases

---

## Overview
Final MVP phase. Creator can publish/unpublish stories with validation (must have start scene). Performance guardrails, seed data, and basic testing.

## Key Insights
- Publish validation prevents broken reader experiences (no start scene = can't read).
- Keep testing minimal but targeted for MVP (KISS).
- Seed data enables manual QA and demo.

## Requirements

### Functional
- Publish endpoint validates story before changing status.
- Unpublish reverts to draft.
- Seed script creates demo data.

### Non-Functional
- API response times consistently < 200ms.
- Request logging with duration.
- Context timeouts on all DB queries.

## Architecture

```
Backend:
  internal/
    http/
      handlers/
        publish.go      # publish/unpublish handlers
    services/
      publish_service.go # validation + status change
  migrations/
    seed.cypher         # demo data
  test/
    integration/        # integration tests
```

## Related Code Files

### Files to Create
- `backend/internal/http/handlers/publish.go`
- `backend/internal/services/publish_service.go`
- `backend/migrations/seed.cypher`
- `backend/test/integration/publish_test.go`
- `backend/test/integration/reader_test.go`

### Files to Modify
- `backend/internal/http/router.go` — add publish/unpublish routes
- `frontend/src/pages/creator/StoryMapPage.vue` — add publish button
- `frontend/src/api/stories.ts` — add publish/unpublish calls

## Implementation Steps

### 1. Publish/unpublish endpoints (M)

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/api/creator/stories/:storyId/publish` | PUT | Validate + set published |
| `/api/creator/stories/:storyId/unpublish` | PUT | Set draft |

**Publish validation rules:**
- Story must have a STARTS_AT relationship (start scene exists).
- Start scene must belong to the story.
- Optional: warn if any scenes have empty content.

**Cypher — Check story has start scene:**
```cypher
MATCH (u:User {id:$user_id})-[:CREATED]->(st:Story {id:$story_id})
OPTIONAL MATCH (st)-[:STARTS_AT]->(sc:Scene)
RETURN st.id AS story_id, (sc IS NOT NULL) AS has_start;
```

**Cypher — Set status:**
```cypher
MATCH (u:User {id:$user_id})-[:CREATED]->(st:Story {id:$story_id})
SET st.status = $status
RETURN st { .id, .status } AS story;
```

**Error responses:**
- `400 Bad Request`: `{ "error": "Story must have a start scene before publishing" }`

### 2. Frontend: Publish button (S)
- Add "Publish" / "Unpublish" button to StoryMapPage toolbar.
- Show confirmation dialog before publishing.
- Display validation errors if publish fails.
- Update story status badge after success.

### 3. Performance guardrails (M)
- **Request timeouts**: Gin middleware with 5s timeout per request.
- **Neo4j query timeouts**: 2s context timeout (set in Phase 00 driver wrapper).
- **Request logging**: middleware that logs method, path, status, duration.
- **Verify indexes**: ensure all constraints/indexes from migrations are applied.

**Log format example:**
```
[GIN] 200 | 12.5ms | GET /api/stories/abc/scenes/xyz
```

### 4. Seed data (S)
`backend/migrations/seed.cypher`:
```cypher
// Create demo user
CREATE (u:User {
  id: "user-demo-001",
  username: "demo_author",
  email: "demo@example.com",
  password_hash: "$2a$10$...", // bcrypt hash of "password123"
  role: "creator"
});

// Create demo story
MATCH (u:User {id: "user-demo-001"})
CREATE (s:Story {
  id: "story-demo-001",
  title: "The Haunted Mansion",
  summary: "You find yourself at the entrance of a mysterious mansion...",
  cover_image: null,
  status: "published",
  created_at: datetime()
})
CREATE (u)-[:CREATED]->(s);

// Create scenes
MATCH (s:Story {id: "story-demo-001"})
CREATE (sc1:Scene {id: "scene-001", content: "You stand before a dark, towering mansion. The front door creaks in the wind.", image_url: null, is_end: false, pos_x: 0.0, pos_y: 0.0})
CREATE (sc2:Scene {id: "scene-002", content: "You push open the heavy door and step inside. A grand hallway stretches before you.", image_url: null, is_end: false, pos_x: 200.0, pos_y: 0.0})
CREATE (sc3:Scene {id: "scene-003", content: "You walk around the mansion but find nothing of interest. You go home.", image_url: null, is_end: true, pos_x: 200.0, pos_y: 200.0})
CREATE (sc4:Scene {id: "scene-004", content: "In the library, you find an ancient book that reveals the mansion's dark secret!", image_url: null, is_end: false, pos_x: 400.0, pos_y: -100.0})
CREATE (sc5:Scene {id: "scene-005", content: "Armed with knowledge, you banish the ghost. The mansion is free! THE END.", image_url: null, is_end: true, pos_x: 600.0, pos_y: -100.0})
CREATE (s)-[:HAS_SCENE]->(sc1)
CREATE (s)-[:HAS_SCENE]->(sc2)
CREATE (s)-[:HAS_SCENE]->(sc3)
CREATE (s)-[:HAS_SCENE]->(sc4)
CREATE (s)-[:HAS_SCENE]->(sc5)
CREATE (s)-[:STARTS_AT]->(sc1)
CREATE (sc1)-[:LEADS_TO {choice_text: "Enter the mansion"}]->(sc2)
CREATE (sc1)-[:LEADS_TO {choice_text: "Walk away"}]->(sc3)
CREATE (sc2)-[:LEADS_TO {choice_text: "Explore the library"}]->(sc4)
CREATE (sc2)-[:LEADS_TO {choice_text: "Leave the mansion"}]->(sc3)
CREATE (sc4)-[:LEADS_TO {choice_text: "Use the knowledge"}]->(sc5);
```

### 5. Testing (M)

**Backend integration tests:**
- `publish_test.go`:
  - Publish story without start scene → expect 400.
  - Publish story with start scene → expect 200, status = "published".
  - Unpublish → status = "draft".
- `reader_test.go`:
  - List stories → only published returned.
  - Get start scene → returns correct scene.
  - Get scene + choices → returns correct data.

**Manual QA checklist:**
- [ ] Register new user as creator
- [ ] Create story with title and summary
- [ ] Add 3+ scenes with content
- [ ] Connect scenes with choices
- [ ] Set start scene
- [ ] Publish story
- [ ] Open homepage → story visible
- [ ] Start reading → make choices → reach end
- [ ] Refresh page → progress restored
- [ ] Undo → goes back correctly
- [ ] Unpublish → story hidden from homepage

## Todo List
- [ ] Implement publish validation logic
- [ ] Implement publish/unpublish handlers
- [ ] Update router with publish routes
- [ ] Add publish button to StoryMapPage
- [ ] Implement request timeout middleware
- [ ] Implement request logging middleware
- [ ] Verify all Neo4j indexes are applied
- [ ] Create seed data Cypher script
- [ ] Write publish integration tests
- [ ] Write reader integration tests
- [ ] Run full manual QA checklist
- [ ] Performance check: scene transition < 200ms

## Success Criteria
- Publishing fails with clear error if start scene missing.
- Published stories appear in discovery; drafts do not.
- Unpublished stories disappear from reader view.
- All API responses < 200ms in local environment.
- Seed data loads correctly and produces a playable demo story.

## Risk Assessment
| Risk | Impact | Mitigation |
|------|--------|------------|
| Story published without valid graph | High | Validate STARTS_AT exists before publish |
| Slow queries on larger datasets | Medium | Indexes on id/status; monitor durations |
| Test environment Neo4j differs from dev | Low | Use same docker-compose for both |

## Security Considerations
- Publish/unpublish requires authentication + ownership verification.
- Status changes logged for audit trail (via request logging).

## Next Steps
- MVP complete! 🎉
- Consider Phase 2 features: variables/conditions, inventory/stats, AI writing assistance.
- Consider server-side progress tracking for cross-device support.
- Consider S3 for image storage as scale increases.
