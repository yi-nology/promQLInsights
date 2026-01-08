.PHONY: help frontend-install frontend-dev frontend-build backend-install backend-dev clean

help:
	@echo "Available commands:"
	@echo "  make frontend-install  - Install frontend dependencies"
	@echo "  make frontend-dev      - Start frontend dev server"
	@echo "  make frontend-build    - Build frontend for production"
	@echo "  make backend-install   - Install backend dependencies"
	@echo "  make backend-dev       - Start backend server"
	@echo "  make clean             - Clean build artifacts"

frontend-install:
	cd frontend && npm install

frontend-dev:
	cd frontend && npm run dev

frontend-build:
	cd frontend && npm run build

backend-install:
	cd backend && go mod download

backend-dev:
	cd backend && go run cmd/api/main.go

clean:
	rm -rf frontend/dist
	rm -rf frontend/node_modules
	rm -rf backend/bin
