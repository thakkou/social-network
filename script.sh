#!/bin/bash

set -e

MODE="start"

case "$1" in
    -i|-init)
        MODE="init"
        ;;
    -r|-refresh)
        MODE="refresh"
        ;;
    ""|-start)
        MODE="start"
        ;;
    *)
        echo "Usage: $0 [-start | -r|-refresh | -i|-init]"
        exit 1
        ;;
esac

if [ "$MODE" = "init" ]; then
    echo "Installing frontend dependencies..."
    (
        cd frontend
        npm install
    )
fi

echo "Starting backend..."

(
    cd backend
    if [ "$MODE" = "refresh" ] || [ "$MODE" = "init" ]; then
        go run . -r
    else
        go run .
    fi
) &

BACKEND_PID=$!
echo "Backend PID: $BACKEND_PID"

echo "Starting frontend..."

(
    cd frontend
    npm run dev
) &

FRONTEND_PID=$!
echo "Frontend PID: $FRONTEND_PID"

echo
echo "========================="
echo " Backend : $BACKEND_PID"
echo " Frontend: $FRONTEND_PID"
echo "========================="

wait $BACKEND_PID $FRONTEND_PID