# Interactive Light Novel Platform (MVP) — Implementation Plan

## Goals (MVP)
Build a web platform that supports:
- **Creator mode**: create/manage stories, build branching graph visually (Vue Flow), edit scenes (text + image upload), create choices (edges), publish story.
- **Reader mode**: discover published stories, read scenes interactively, choose branches, track progress + undo/history.
- **Neo4j graph model**: `User`, `Story`, `Scene` nodes; relationships `CREATED`, `HAS_SCENE`, `STARTS_AT`, `LEADS_TO(choice_text)`.

## Non-goals (explicitly out of scope)
- Variables/conditions, inventory/stats, AI writing, monetization, multi-language.

## Key Architecture Decisions (KISS/YAGNI)
1. **Backend framework**: **Gin**
   - Mature ecosystem, common patterns, strong middleware support, easy team onboarding, stable documentation.
   - Fiber is also fast, but Gin is the "safe default" for maintainability.
2. **Image storage (MVP)**: local filesystem served by backend (or via nginx in prod)
   - Simple path: upload → store on disk → return URL.
   - Future: S3-compatible object storage.
3. **Auth (MVP)**: JWT access token in **HttpOnly cookie**
   - Avoid storing tokens in LocalStorage (XSS risk).
   - Simple RBAC: roles `creator` and `reader`.
4. **Progress tracking (MVP)**:
   - Primary: browser **LocalStorage** (per PRD).
   - Optional later: server-side progress per user/story.

## Monorepo Structure
```
repo-root/
  frontend/
    index.html
    package.json
    vite.config.ts
    src/
      app/
      pages/
      components/
      stores/            # Pinia
      api/               # API client wrappers
      router/
      assets/
      styles/
  backend/
    go.mod
    cmd/
      api/               # main.go
    internal/
      config/
      http/              # Gin router, handlers, middleware
      domain/            # core types (Story, Scene, User)
      services/          # business logic
      repo/              # Neo4j data access
      auth/              # JWT helpers
      storage/           # image storage (filesystem)
    migrations/          # cypher migration scripts
    test/
  shared/
    openapi/             # openapi.yaml (source of truth)
  scripts/
    dev.sh
  docker/
    docker-compose.yml   # neo4j + dev services
  README.md
```

## Implementation Phases

| Phase | Name | Depends On | Estimate | Status |
|-------|------|------------|----------|--------|
| [Phase 00](phase-00-foundation.md) | Foundation (Repo, DevEnv, Baselines) | — | M | ⬜ TODO |
| [Phase 01](phase-01-auth-users.md) | Authentication + Users (JWT + Roles) | Phase 00 | M | ⬜ TODO |
| [Phase 02](phase-02-creator-story-crud.md) | Creator: Story Management (CRUD) | Phase 01 | M–L | ⬜ TODO |
| [Phase 03](phase-03-creator-story-map-scenes-choices.md) | Creator: Visual Story Map + Scenes + Choices | Phase 02 | L | ⬜ TODO |
| [Phase 04](phase-04-reader-reading-progress-preload.md) | Reader: Discovery + Reading + Progress + Preload | Phase 03 | L | ⬜ TODO |
| [Phase 05](phase-05-publishing-perf-polish.md) | Publishing, Validation, Performance & Polish | All | M | ⬜ TODO |

## Performance Approach (meet <200ms scene transitions)
- Neo4j constraints + indexes on IDs/status.
- Scene transition endpoint returns **scene + choices + adjacent scene metadata** in one query.
- Frontend preloads adjacent scenes & images after each transition.

## Testing Strategy (minimal but effective)
- Backend: unit tests for services + repository integration tests (Neo4j in docker).
- Frontend: component tests for story map editor + e2e smoke test for reading flow.
- Add simple seed data for manual QA.

## Estimates Legend
- **S**: < 1 day
- **M**: 1–3 days
- **L**: 3–7 days
- **XL**: > 7 days

## REST API Endpoints (Consolidated)

### Auth
- `POST /api/auth/register`
- `POST /api/auth/login`
- `POST /api/auth/logout`
- `GET /api/me`

### Uploads
- `POST /api/uploads/images` (multipart form)
- `GET /uploads/:filename` (static)

### Creator — Stories
- `POST /api/creator/stories`
- `GET /api/creator/stories`
- `GET /api/creator/stories/:storyId`
- `PATCH /api/creator/stories/:storyId`
- `DELETE /api/creator/stories/:storyId`
- `PUT /api/creator/stories/:storyId/publish`
- `PUT /api/creator/stories/:storyId/unpublish`

### Creator — Scenes & Graph
- `POST /api/creator/stories/:storyId/scenes`
- `GET /api/creator/stories/:storyId/scenes/:sceneId`
- `PATCH /api/creator/stories/:storyId/scenes/:sceneId`
- `DELETE /api/creator/stories/:storyId/scenes/:sceneId`
- `PUT /api/creator/stories/:storyId/start/:sceneId`
- `GET /api/creator/stories/:storyId/graph`

### Creator — Choices (Edges)
- `POST /api/creator/stories/:storyId/choices`
- `PATCH /api/creator/stories/:storyId/choices`
- `DELETE /api/creator/stories/:storyId/choices`

### Reader
- `GET /api/stories` (published)
- `GET /api/stories/:storyId/start`
- `GET /api/stories/:storyId/scenes/:sceneId`

## Top Risks & Mitigations

| Risk | Mitigation |
|------|------------|
| Graph consistency / cross-story edges | Constrain edge creation: both scenes must belong to same story via `(st)-[:HAS_SCENE]->` |
| Performance regression on scene transitions | Single Cypher query per transition; indexes on `Story.id`, `Story.status`, `Scene.id` |
| Vue Flow complexity for juniors | Keep MVP simple: controlled nodes/edges, one custom node type, modal for choice text, debounce position saving |
| Auth security pitfalls | HttpOnly cookie JWT; server-side middleware; no tokens in LocalStorage |
| Image handling (storage, size, bandwidth) | Basic file size limits; restrict mime types; preload only adjacent images |
