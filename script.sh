#!/bin/sh

echo "Building social network containers..."

echo ""
echo "=== Building Backend ==="
docker build -t sn-backend -f Dockerfile.backend .

echo ""
echo "=== Building Frontend ==="
docker build -t sn-frontend -f Dockerfile.frontend .

echo ""
echo "=== Starting services ==="
docker compose up -d

echo ""
echo "Done! Backend running on http://localhost:8080"
echo "Frontend running on http://localhost:3000"
