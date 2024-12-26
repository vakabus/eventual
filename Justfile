run:
    cd backend && go generate && DEV=true go run .

lint:
    cd frontend && npm run lint
    cd backend && golangci-lint run

format:
    cd frontend && npm run format
    cd backend && go fmt ./...
