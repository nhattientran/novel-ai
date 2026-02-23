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
- **Package Manager**: pnpm (frontend)

## Common Commands

### Frontend
```bash
cd frontend
pnpm install       # or npm install
pnpm dev         # Start dev server at http://localhost:5173
pnpm build       # Production build with TypeScript check
pnpm preview     # Preview production build
pnpm lint        # ESLint on .vue,.ts,.tsx files
```

### Backend
```bash
cd backend
go mod tidy           # Download dependencies
go build ./cmd/api    # Build binary
go run ./cmd/api      # Start API server at http://localhost:8080
```

### Environment Setup
Backend requires `.env` file in project root (not backend/). Copy from `.env.example`:
```bash
cp .env backend/.env
# Or export directly:
export NEO4J_URI=bolt://localhost:7687 && export NEO4J_USER=neo4j && export NEO4J_PASSWORD=... && export JWT_SECRET=...
```

### Database
```bash
cd docker && docker compose up -d   # Start Neo4j at http://localhost:7474
```

### Quick Start
```bash
./scripts/dev.sh   # Setup: start Neo4j, create .env if missing
```

## Architecture

### Project Structure
- `docs/` - Documentation (PRD.md)
- `plans/` - Implementation plans with TODO tasks
- `scripts/` - Development utility scripts
- `shared/openapi/` - OpenAPI specification

### Frontend Structure
- `src/pages/` - Route-level components (HomePage, LoginPage, StoryListPage, StoryMapPage, ReadPage)
- `src/components/` - Reusable components organized by feature (reading/, scene-editor/, story-map/)
- `src/router/` - Vue Router configuration
- `src/stores/` - Pinia state management
- `src/api/` - API client wrappers

### Backend Structure
- `cmd/api/main.go` - Application entry point with graceful shutdown
- `internal/config/` - Environment configuration (NEO4J_URI, JWT_SECRET, PORT)
- `internal/domain/` - Domain models (User, Story, Scene)
- `internal/http/` - HTTP handlers and Gin router setup
- `internal/repo/` - Neo4j data access layer
- `internal/auth/` - JWT authentication logic
- `internal/storage/` - File storage (local filesystem)

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
- Creator routes: `/api/creator/*` (requires auth + creator role)
- Static files: `/uploads/*` served from `./uploads/`

## Development Workflow

1. Start Neo4j: `cd docker && docker compose up -d`
2. Start backend: `cd backend && go run ./cmd/api`
3. Start frontend: `cd frontend && npm run dev`
4. Access points:
   - Frontend: http://localhost:5173
   - Backend API: http://localhost:8080
   - Neo4j Browser: http://localhost:7474

### Build & Type Check
```bash
# Frontend - must pass TypeScript check
cd frontend && pnpm build

# Backend
cd backend && go build ./cmd/api
```

## Environment Variables

Copy `.env.example` to `.env` and configure:
- `NEO4J_URI` - Neo4j connection string (default: bolt://localhost:7687)
- `NEO4J_USER` / `NEO4J_PASSWORD` - Neo4j credentials
- `JWT_SECRET` - Secret for JWT signing
- `PORT` - API server port (default: 8080)
