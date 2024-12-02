run:
    cd frontend && npm run build
    cd backend && go run .

lint:
    cd frontend && npm run lint
    cd backend && golangci-lint run
