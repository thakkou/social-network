#!/bin/sh

echo "Building social network containers..."

echo ""
echo "=== rmoving old containers ==="
docker compose down -v

echo ""
echo "=== Building new ==="
docker compose up --build


echo "Done! Backend running on http://localhost:8080"
echo "Frontend running on http://localhost:3000"
