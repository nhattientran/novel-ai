#!/bin/bash

# Development startup script

echo "Starting Novel AI development environment..."

# Start Neo4j
echo "Starting Neo4j..."
cd docker
docker compose up -d
cd ..

# Check if .env exists
if [ ! -f .env ]; then
    echo "Creating .env from .env.example..."
    cp .env.example .env
fi

# Load environment variables
export $(cat .env | xargs)

echo ""
echo "Development environment ready!"
echo ""
echo "Next steps:"
echo "1. Start backend: cd backend && go run ./cmd/api"
echo "2. Start frontend: cd frontend && npm run dev"
echo ""
echo "Access points:"
echo "- Frontend: http://localhost:5173"
echo "- Backend:  http://localhost:8080"
echo "- Neo4j:    http://localhost:7474"
