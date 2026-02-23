# Novel AI - Nền tảng Light Novel Tương Tác

Một nền tảng web cho phép ngườii dùng sáng tạo, xuất bản và trải nghiệm các bộ light novel có cốt truyện rẽ nhánh.

## Kiến trúc

```
repo-root/
  frontend/          # Vue 3 + Vite + TypeScript
  backend/           # Go + Gin
  docker/            # Docker Compose cho Neo4j
  shared/openapi/    # API contract
```

## Yêu cầu

- Go 1.23+
- Node.js 20+
- Docker & Docker Compose

## Cài đặt

### 1. Clone repository

```bash
git clone <repo-url>
cd novel-ai
```

### 2. Setup môi trường

```bash
cp .env.example .env
# Chỉnh sửa .env nếu cần
```

### 3. Khởi động Neo4j

```bash
cd docker
docker compose up -d
```

### 4. Khởi động Backend

```bash
cd backend
go mod tidy
go run ./cmd/api
```

### 5. Khởi động Frontend

```bash
cd frontend
npm install
npm run dev
```

## Truy cập

- Frontend: http://localhost:5173
- Backend API: http://localhost:8080
- Neo4j Browser: http://localhost:7474

## API Endpoints

### Health
- `GET /api/health` - Health check
- `GET /api/ready` - Readiness check

## Phát triển

Xem chi tiết tại `plans/20260223-1200-interactive-light-novel-mvp/`.
