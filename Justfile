run:
    #!/bin/bash
    set -e

    cd backend
    go generate ./...
    DEV=true go run github.com/air-verse/air@latest &

    cd ../frontend
    npm run dev

lint:
    cd frontend && npm run lint
    cd backend && golangci-lint run

format:
    cd frontend && npm run format
    cd backend && go fmt ./...
