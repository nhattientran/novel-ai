# Phase 00 — Foundation (Repo, DevEnv, Baselines)

**Date:** 2026-02-23
**Priority:** Critical (blocks all phases)
**Status:** ⬜ TODO
**Estimate:** M (1–3 days)

---

## Overview
Set up monorepo, Docker dev environment with Neo4j, backend skeleton (Gin), and frontend skeleton (Vue 3 + Pinia + Router).

## Key Insights
- Keep dev setup minimal: docker-compose for Neo4j only, run backend/frontend natively.
- Single `.env` approach from day 1 to avoid config drift.

## Requirements

### Functional
- Monorepo with `frontend/`, `backend/`, `docker/`, `shared/`, `scripts/` directories.
- Neo4j accessible via docker-compose.
- Backend with health check endpoints.
- Frontend with routing structure.

### Non-Functional
- Dev setup should work with `docker compose up` + `go run` + `npm run dev`.

## Architecture

```
repo-root/
├── frontend/          # Vue 3 + Vite
├── backend/           # Go + Gin
│   ├── cmd/api/       # main.go entrypoint
│   └── internal/      # all internal packages
├── docker/            # docker-compose.yml
├── shared/openapi/    # openapi.yaml
├── scripts/           # dev helpers
└── .env.example
```

## Related Code Files

### Files to Create
- `docker/docker-compose.yml` — Neo4j service
- `.env.example` — environment variables template
- `backend/go.mod` — Go module init
- `backend/cmd/api/main.go` — server entrypoint
- `backend/internal/config/config.go` — env config loader
- `backend/internal/http/router.go` — Gin router setup
- `backend/internal/http/handlers/health.go` — health/ready handlers
- `backend/internal/repo/neo4j.go` — Neo4j driver wrapper
- `frontend/package.json` — Vue 3 + Vite + Pinia + Vue Router deps
- `frontend/vite.config.ts` — Vite config with API proxy
- `frontend/src/main.ts` — Vue app entry
- `frontend/src/router/index.ts` — route definitions
- `frontend/src/stores/index.ts` — Pinia setup
- `frontend/src/pages/HomePage.vue` — placeholder
- `frontend/src/pages/LoginPage.vue` — placeholder
- `frontend/src/pages/RegisterPage.vue` — placeholder
- `frontend/src/pages/creator/StoryListPage.vue` — placeholder
- `frontend/src/pages/creator/StoryMapPage.vue` — placeholder
- `frontend/src/pages/reader/ReadPage.vue` — placeholder
- `shared/openapi/openapi.yaml` — API contract placeholder

## Implementation Steps

### 1. Monorepo scaffolding (S)
- Create all directories as specified in architecture.
- Add root `README.md` with "How to run" instructions.

### 2. Docker Compose for Neo4j (S)
- `docker/docker-compose.yml`:
  ```yaml
  services:
    neo4j:
      image: neo4j:5
      ports:
        - "7474:7474"
        - "7687:7687"
      environment:
        NEO4J_AUTH: ${NEO4J_USER}/${NEO4J_PASSWORD}
      volumes:
        - neo4j_data:/data
  volumes:
    neo4j_data:
  ```
- `.env.example`:
  ```
  NEO4J_URI=bolt://localhost:7687
  NEO4J_USER=neo4j
  NEO4J_PASSWORD=password
  JWT_SECRET=change-me-in-production
  PORT=8080
  ```

### 3. Backend: Gin skeleton (M)
- Init Go module: `go mod init github.com/<org>/novel-ai`
- `cmd/api/main.go`: load config, init Neo4j driver, start Gin server.
- `internal/config/config.go`: read env vars (use `os.Getenv` or simple lib like `envconfig`).
- `internal/http/router.go`: create Gin engine, register routes.
- `internal/http/handlers/health.go`:
  - `GET /health` → `{"status":"ok"}`
  - `GET /ready` → check Neo4j with `RETURN 1` query

### 4. Neo4j driver wrapper (M)
- `internal/repo/neo4j.go`:
  - Create driver from `NEO4J_URI`, `NEO4J_USER`, `NEO4J_PASSWORD`
  - Helper functions for read/write transactions with context timeouts (2s default)
  - Close driver gracefully on shutdown

### 5. Frontend: Vue 3 skeleton (M)
- Scaffold with `npm create vite@latest frontend -- --template vue-ts`
- Install: `pinia`, `vue-router`, `@vue-flow/core` (install now, use in Phase 03)
- Routes:
  - `/` → HomePage
  - `/login` → LoginPage
  - `/register` → RegisterPage
  - `/creator/stories` → StoryListPage
  - `/creator/stories/:id/map` → StoryMapPage
  - `/read/:storyId` → ReadPage
- Vite config: proxy `/api` to `http://localhost:8080`

### 6. API contract placeholder (S)
- `shared/openapi/openapi.yaml` with `/health` endpoint defined.

## Todo List
- [ ] Create monorepo directory structure
- [ ] Create docker-compose.yml with Neo4j
- [ ] Create .env.example
- [ ] Init Go module + Gin server skeleton
- [ ] Implement config loader
- [ ] Implement Neo4j driver wrapper
- [ ] Implement /health and /ready endpoints
- [ ] Scaffold Vue 3 frontend with Vite
- [ ] Set up Vue Router with all route placeholders
- [ ] Set up Pinia store
- [ ] Configure Vite proxy for API
- [ ] Create openapi.yaml placeholder
- [ ] Verify: `docker compose up` starts Neo4j
- [ ] Verify: backend /health and /ready work
- [ ] Verify: frontend runs with routing

## Success Criteria
- `docker compose up` starts Neo4j successfully.
- `go run ./cmd/api` starts backend; `/health` returns 200, `/ready` confirms Neo4j connection.
- `npm run dev` starts frontend; routes navigate correctly.

## Risk Assessment
| Risk | Impact | Mitigation |
|------|--------|------------|
| Neo4j version incompatibility with Go driver | Medium | Pin Neo4j 5.x and matching Go driver version |
| Config drift between dev environments | Low | Single `.env.example` as source of truth |

## Security Considerations
- No secrets committed; `.env` in `.gitignore`.
- Neo4j default credentials changed via env vars.

## Next Steps
- → Phase 01: Authentication + Users
