# Makefile (project root)

# variables
BACKEND_DIR := Backend
FRONTEND_DIR := Frontend
IMAGE_NAME := backend-dev

.PHONY: help dev back-dev back-build back-run front-dev front-build docker-build clean

help:
	@echo "❯ make dev         # start both backend+frontend watchers"
	@echo "❯ make back-dev    # start backend (Air) in dev mode"
	@echo "❯ make front-dev   # start frontend (e.g. vite/quasar) in dev"
	@echo "❯ make back-build  # build backend binary"
	@echo "❯ make front-build # build frontend for production"
	@echo "❯ make docker-build  # build backend-dev Docker image"
	@echo "❯ make clean       # remove tmp artifacts"

dev:
	@echo "Starting both backend and frontend watchers..."
	$(MAKE) -j2 back-dev front-dev

### backend tasks ###
back-dev:
	cd $(BACKEND_DIR) && air

back-build:
	cd $(BACKEND_DIR) && go build -o main ./cmd/web

back-run:
	cd $(BACKEND_DIR) && ./main

### frontend tasks ###
front-dev:
	cd $(FRONTEND_DIR) && bun dev  # or `quasar dev` etc.

front-build:
	cd $(FRONTEND_DIR) && bun run build

### docker ###
docker-build:
	docker build -f $(BACKEND_DIR)/Dockerfile -t $(IMAGE_NAME) .

clean:
	-rm -rf $(BACKEND_DIR)/tmp
	-rm -rf $(FRONTEND_DIR)/dist