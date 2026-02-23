# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Novel AI is an interactive light novel platform with a Vue 3 frontend and Go backend using Neo4j as a graph database. The platform supports:
- **Creator mode**: Visual story mapping with branching narratives using Vue Flow
- **Reader mode**: Interactive reading with choices and progress tracking

## Technology Stack

- **Frontend**: Vue 3 + Vite + TypeScript + Pinia + Vue Router + Vue Flow
- **Backend**: Go 1.23+ + Gin + Neo4j (neo4j-go-driver/v5)
- **Database**: Neo4j (graph database for story/scene/choice relationships)
- **Infrastructure**: Docker Compose for Neo4j

## Common Commands

### Frontend
```bash
cd frontend
npm install
npm run dev      # Start dev server at http://localhost:5173
npm run build    # Production build (vue-tsc && vite build)
npm run preview  # Preview production build
npm run lint     # ESLint on .vue,.ts,.tsx files
```

### Backend
```bash
cd backend
go mod tidy           # Download dependencies
go run ./cmd/api      # Start API server at http://localhost:8080
```

### Database
```bash
cd docker && docker compose up -d   # Start Neo4j at http://localhost:7474
```

## Architecture

### Frontend Structure
- `src/pages/` - Route-level components (HomePage, LoginPage, StoryListPage, StoryMapPage, ReadPage)
- `src/components/` - Reusable components organized by feature (reading/, scene-editor/, story-map/)
- `src/router/` - Vue Router configuration
- `src/stores/` - Pinia state management
- `src/api/` - API client wrappers

### Backend Structure
- `cmd/api/main.go` - Application entry point with graceful shutdown
- `internal/config/` - Environment configuration (NEO4J_URI, JWT_SECRET, PORT)
- `internal/http/` - HTTP handlers and Gin router setup
- `internal/repo/` - Neo4j data access layer
- `internal/auth/` - JWT authentication logic

### Data Model (Neo4j Graph)
- **Nodes**: `User`, `Story`, `Scene`
- **Relationships**:
  - `CREATED` (User → Story)
  - `HAS_SCENE` (Story → Scene)
  - `STARTS_AT` (Story → Scene) - marks starting scene
  - `LEADS_TO` (Scene → Scene) - choices with `choice_text` property

### API Design
- REST API with `/api` prefix
- JWT authentication using HttpOnly cookies
- OpenAPI spec at `shared/openapi/openapi.yaml`
- Health endpoints: `GET /api/health`, `GET /api/ready`

## Development Workflow

1. Start Neo4j: `cd docker && docker compose up -d`
2. Start backend: `cd backend && go run ./cmd/api`
3. Start frontend: `cd frontend && npm run dev`
4. Access points:
   - Frontend: http://localhost:5173
   - Backend API: http://localhost:8080
   - Neo4j Browser: http://localhost:7474

## Environment Variables

Copy `.env.example` to `.env` and configure:
- `NEO4J_URI` - Neo4j connection string (default: bolt://localhost:7687)
- `NEO4J_USER` / `NEO4J_PASSWORD` - Neo4j credentials
- `JWT_SECRET` - Secret for JWT signing
- `PORT` - API server port (default: 8080)
