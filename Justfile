build:
    cd backend && go generate ./...
    cd backend && go build

run:
    #!/bin/bash
    set -e

    cd backend
    go generate ./...
    DEV=true go run github.com/air-verse/air@latest &
    trap "killall air" EXIT

    cd ../frontend
    npm run dev

lint:
    cd frontend && npm run lint
    cd backend && golangci-lint run

format:
    cd frontend && npm run format
    cd backend && go fmt ./...

sqlc:
    cd backend && go generate ./database
