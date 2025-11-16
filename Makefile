# Makefile para proyecto webPolls
.PHONY: help build run test clean stop logs restart dev install sqlc fmt vet seed seed-users setup-test

# Variables
DOCKER_COMPOSE := docker compose
PROJECT_NAME := webpolls
BINARY_NAME := webpolls
BINARY_PATH := ./bin/$(BINARY_NAME)
GO_FILES := $(shell find . -name '*.go' -not -path './vendor/*')

# Colores para output
BLUE := \033[34m
GREEN := \033[32m
YELLOW := \033[33m
RED := \033[31m
NC := \033[0m # No Color

# Target por defecto
.DEFAULT_GOAL := help

## help: Mostrar esta ayuda
help:
	@echo "$(BLUE)Comandos disponibles:$(NC)"
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' | sed -e 's/^/ /'

## install: Instalar dependencias de Go
install:
	@echo "$(BLUE)Instalando dependencias...$(NC)"
	go mod download
	go mod tidy

## build: Compilar el binario de Go
build:
	@echo "$(BLUE)Compilando binario...$(NC)"
	mkdir -p bin
	go build -o $(BINARY_PATH) ./main.go
	@echo "$(GREEN)Binario compilado en $(BINARY_PATH)$(NC)"

## run-local: Ejecutar la aplicación localmente (requiere base de datos)
run-local: build
	@echo "$(BLUE)Ejecutando aplicación localmente...$(NC)"
	$(BINARY_PATH)

## docker-build: Construir imagen Docker
docker-build:
	@echo "$(BLUE)Construyendo imagen Docker...$(NC)"
	$(DOCKER_COMPOSE) build --no-cache

## docker-up: Levantar contenedores en background
docker-up:
	@echo "$(BLUE)Levantando contenedores...$(NC)"
	$(DOCKER_COMPOSE) up -d

## docker-down: Bajar contenedores
docker-down:
	@echo "$(BLUE)Bajando contenedores...$(NC)"
	$(DOCKER_COMPOSE) down

## clean: Limpiar contenedores y volúmenes
clean:
	@echo "$(BLUE)Limpiando contenedores anteriores...$(NC)"
	-docker rm -f $$(docker ps -a -q --filter "name=$(PROJECT_NAME)-*") 2>/dev/null || true
	-docker volume prune -f 2>/dev/null || true
	@echo "$(GREEN)Limpieza completada$(NC)"

## run: Ejecutar aplicación completa (limpia, construye y ejecuta)
run: clean docker-build docker-up wait-for-app seed
	@echo "$(GREEN)Aplicación lista en http://localhost:8080$(NC)"


## wait-for-app: Esperar a que la aplicación esté lista
wait-for-app:
	@echo "$(BLUE)Esperando a que la aplicación se inicie...$(NC)"
	@while ! curl -s --fail http://localhost:8080/ > /dev/null 2>&1; do \
		echo -n "."; \
		sleep 1; \
	done
	@echo ""
	@echo "$(GREEN)La aplicación se ha iniciado correctamente$(NC)"

## logs: Mostrar logs de los contenedores
logs:
	@echo "$(BLUE)Mostrando logs...$(NC)"
	$(DOCKER_COMPOSE) logs -f

## logs-backend: Mostrar solo logs del backend
logs-backend:
	@echo "$(BLUE)Mostrando logs del backend...$(NC)"
	$(DOCKER_COMPOSE) logs -f backend

## logs-db: Mostrar solo logs de la base de datos
logs-db:
	@echo "$(BLUE)Mostrando logs de PostgreSQL...$(NC)"
	$(DOCKER_COMPOSE) logs -f postgres

## restart: Reiniciar contenedores
restart: docker-down docker-up wait-for-app
	@echo "$(GREEN)Aplicación reiniciada$(NC)"

## dev: Modo desarrollo - rebuild y restart rápido
dev: docker-build restart

## stop: Parar todos los contenedores
stop: docker-down

## status: Mostrar estado de contenedores
status:
	@echo "$(BLUE)Estado de contenedores:$(NC)"
	docker ps --filter "name=$(PROJECT_NAME)-"

## sqlc: Generar código SQLC
sqlc:
	@echo "$(BLUE) Generando código SQLC...$(NC)"
	sqlc generate
	@echo "$(GREEN) Código SQLC generado$(NC)"

templ:
	@echo "$(BLUE) Generando código Templ...$(NC)"
	templ generate
	@echo "$(GREEN) Código Templ generado$(NC)"

## shell-backend: Conectar al shell del contenedor backend
shell-backend:
	@echo "$(BLUE)Conectando al contenedor backend...$(NC)"
	docker exec -it $(PROJECT_NAME)-backend sh

## shell-db: Conectar a PostgreSQL
shell-db:
	@echo "$(BLUE)Conectando a PostgreSQL...$(NC)"
	docker exec -it $(PROJECT_NAME)-postgres psql -U postgres -d polls

## seed: Crear usuario de prueba (agus)
seed:
	@./seed.sh
