# Makefile para proyecto webPolls
.PHONY: help build run test clean stop logs restart dev install sqlc fmt vet seed seed-users setup-test install-deps css css-watch air templ-watch dev-live db

# Variables
DOCKER_COMPOSE := docker compose
PROJECT_NAME := webpolls
BINARY_NAME := webpolls
BINARY_PATH := ./bin/$(BINARY_NAME)

# Target por defecto
.DEFAULT_GOAL := help

## help: Mostrar esta ayuda
help:
	@echo Comandos disponibles:
	@echo   docker-build  - Construir imagen Docker
	@echo   docker-up     - Levantar contenedores
	@echo   docker-down   - Bajar contenedores
	@echo   clean         - Limpiar contenedores, volúmenes e imágenes locales
	@echo   run           - Ejecutar aplicación completa
	@echo   logs          - Mostrar logs
	@echo   logs-backend  - Mostrar logs del backend
	@echo   logs-db       - Mostrar logs de la base de datos
	@echo   restart       - Reiniciar contenedores
	@echo   dev           - Modo desarrollo
	@echo   stop          - Parar contenedores
	@echo   status        - Mostrar estado
	@echo   shell-backend - Shell del backend
	@echo   shell-db      - Shell de la BD
	@echo   install-deps  - Instalar dependencias npm
	@echo   db            - Levantar solo la base de datos
	@echo   css           - Construir CSS
	@echo   css-watch     - Observar cambios en CSS
	@echo   air           - Ejecutar con live reload (Air)
	@echo   templ-watch   - Generar templates con watch
	@echo   dev-live      - Desarrollo local con Air, Tailwind y Templ

## install-deps: Instalar dependencias npm
install-deps:
	pnpm install

## css: Construir CSS con Tailwind
css:
	pnpm build:css

## css-watch: Observar cambios en CSS
css-watch:
	pnpm watch:css

## air: Ejecutar con live reload
air:
	@if ! command -v air > /dev/null; then \
		echo "Installing air..."; \
		go install github.com/air-verse/air@latest; \
	fi
	air

## templ-watch: Generar templates con watch
templ-watch:
	@if ! command -v templ > /dev/null; then \
		echo "Installing templ..."; \
		go install github.com/a-h/templ/cmd/templ@latest; \
	fi
	templ generate --watch

## dev-live: Desarrollo local con Air, Tailwind y Templ (paralelo)
dev-live:
	$(MAKE) -j3 css-watch air templ-watch

## db: Levantar solo la base de datos
db:
	$(DOCKER_COMPOSE) up -d postgres

## docker-build: Construir imagen Docker
docker-build:
	$(DOCKER_COMPOSE) build --no-cache

## docker-up: Levantar contenedores efímeros (recrear y renovar volúmenes)
docker-up:
	$(DOCKER_COMPOSE) up -d --force-recreate --renew-anon-volumes

## docker-down: Bajar contenedores y eliminar volúmenes
docker-down:
	$(DOCKER_COMPOSE) down --volumes --remove-orphans

## clean: Limpiar todo (contenedores, volúmenes, imágenes locales)
clean:
	$(DOCKER_COMPOSE) down --volumes --rmi local --remove-orphans

## run: Ejecutar aplicación completa (limpia, construye y ejecuta)
run: clean docker-build docker-up

## logs: Mostrar logs de los contenedores
logs:
	$(DOCKER_COMPOSE) logs -f

## logs-backend: Mostrar solo logs del backend
logs-backend:
	$(DOCKER_COMPOSE) logs -f backend

## logs-db: Mostrar solo logs de la base de datos
logs-db:
	$(DOCKER_COMPOSE) logs -f postgres

## restart: Reiniciar contenedores
restart: docker-down docker-up

## dev: Modo desarrollo - rebuild y restart rápido
dev: docker-build restart

## stop: Parar todos los contenedores
stop: docker-down

## status: Mostrar estado de contenedores
status:
	docker ps --filter "name=$(PROJECT_NAME)-"

## shell-backend: Conectar al shell del contenedor backend
shell-backend:
	docker exec -it $(PROJECT_NAME)-backend sh

## shell-db: Conectar a PostgreSQL
shell-db:
	docker exec -it $(PROJECT_NAME)-postgres psql -U postgres -d polls

docker-push:
	docker build -t tecdaz/webpolls:latest .
	docker push tecdaz/webpolls:latest
