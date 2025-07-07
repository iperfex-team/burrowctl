# Makefile para burrowctl
# Versión por defecto
VERSION ?= v1.2.2

# Configuración del proyecto
PROJECT_NAME = burrowctl
MODULE_NAME = github.com/lordbasex/burrowctl
GO_VERSION = 1.22.0

# Colores para output
RED = \033[0;31m
GREEN = \033[0;32m
YELLOW = \033[0;33m
BLUE = \033[0;34m
NC = \033[0m # No Color

# Ayuda por defecto
.PHONY: help
help: ## Muestra esta ayuda
	@echo "$(BLUE)$(PROJECT_NAME) - Makefile$(NC)"
	@echo "Uso: make [target]"
	@echo ""
	@echo "Targets disponibles:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(YELLOW)%-15s$(NC) %s\n", $$1, $$2}'

# Tareas de desarrollo
.PHONY: install
install: ## Instala dependencias
	@echo "$(GREEN)📦 Instalando dependencias...$(NC)"
	go mod download
	go mod tidy

.PHONY: build
build: build-examples ## Compila el proyecto y ejemplos
	@echo "$(GREEN)🔨 Compilando proyecto...$(NC)"
	go build -v ./...

.PHONY: test
test: test-examples ## Ejecuta tests del proyecto y ejemplos
	@echo "$(GREEN)🧪 Ejecutando tests...$(NC)"
	go test -v ./...

.PHONY: test-coverage
test-coverage: ## Ejecuta tests con cobertura
	@echo "$(GREEN)📊 Ejecutando tests con cobertura...$(NC)"
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)Reporte de cobertura generado: coverage.html$(NC)"

.PHONY: lint
lint: ## Ejecuta el linter
	@echo "$(GREEN)🔍 Ejecutando linter...$(NC)"
	@which golangci-lint > /dev/null || (echo "$(RED)golangci-lint no está instalado. Instálalo con: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest$(NC)" && exit 1)
	golangci-lint run

.PHONY: fmt
fmt: ## Formatea el código
	@echo "$(GREEN)✨ Formateando código...$(NC)"
	go fmt ./...

.PHONY: vet
vet: ## Ejecuta go vet
	@echo "$(GREEN)🔍 Ejecutando go vet...$(NC)"
	go vet ./...

# Tareas de limpieza
.PHONY: clean
clean: clean-examples ## Limpia archivos generados
	@echo "$(GREEN)🧹 Limpiando archivos generados...$(NC)"
	go clean
	rm -f coverage.out coverage.html

# Tareas de release
.PHONY: check-git-clean
check-git-clean: ## Verifica que el git esté limpio
	@echo "$(GREEN)🔍 Verificando estado de git...$(NC)"
	@if [ -n "$$(git status --porcelain)" ]; then \
		echo "$(RED)❌ Error: Hay cambios sin commit$(NC)"; \
		git status --short; \
		exit 1; \
	fi
	@echo "$(GREEN)✅ Git está limpio$(NC)"

.PHONY: check-main-branch
check-main-branch: ## Verifica que estés en la rama main
	@echo "$(GREEN)🔍 Verificando rama actual...$(NC)"
	@CURRENT_BRANCH=$$(git branch --show-current); \
	if [ "$$CURRENT_BRANCH" != "main" ]; then \
		echo "$(RED)❌ Error: No estás en la rama main (actual: $$CURRENT_BRANCH)$(NC)"; \
		exit 1; \
	fi
	@echo "$(GREEN)✅ Estás en la rama main$(NC)"

.PHONY: pre-release-checks
pre-release-checks: check-main-branch check-git-clean test vet ## Ejecuta todas las verificaciones antes del release
	@echo "$(GREEN)✅ Todas las verificaciones pasaron$(NC)"

.PHONY: tag
tag: pre-release-checks ## Crea tag de versión (usar VERSION=vX.Y.Z)
	@echo "$(GREEN)🏷️  Creando tag $(VERSION)...$(NC)"
	@if git tag -l | grep -q "^$(VERSION)$$"; then \
		echo "$(RED)❌ Error: El tag $(VERSION) ya existe$(NC)"; \
		exit 1; \
	fi
	git tag -a $(VERSION) -m "Release $(VERSION)"
	@echo "$(GREEN)✅ Tag $(VERSION) creado$(NC)"

.PHONY: push
push: ## Hace push del código y tags
	@echo "$(GREEN)🚀 Haciendo push del código...$(NC)"
	git push origin main
	@echo "$(GREEN)🚀 Haciendo push de los tags...$(NC)"
	git push origin --tags
	@echo "$(GREEN)✅ Push completado$(NC)"

.PHONY: release
release: tag push ## Crea tag y hace push (versión completa)
	@echo "$(GREEN)🎉 Release $(VERSION) completado exitosamente!$(NC)"
	@echo "$(BLUE)Para verificar: git ls-remote --tags origin$(NC)"

.PHONY: quick-release
quick-release: ## Release rápido con VERSION por defecto
	@echo "$(GREEN)🚀 Iniciando release rápido con versión $(VERSION)...$(NC)"
	$(MAKE) release VERSION=$(VERSION)

# Información del proyecto
.PHONY: info
info: ## Muestra información del proyecto
	@echo "$(BLUE)Información del proyecto:$(NC)"
	@echo "  Nombre: $(PROJECT_NAME)"
	@echo "  Módulo: $(MODULE_NAME)"
	@echo "  Versión Go: $(GO_VERSION)"
	@echo "  Versión por defecto: $(VERSION)"
	@echo "  Rama actual: $$(git branch --show-current)"
	@echo "  Último commit: $$(git log -1 --oneline)"
	@echo "  Tags existentes: $$(git tag -l | tail -5 | tr '\n' ' ')"

# Tareas para ejemplos
.PHONY: build-examples
build-examples: ## Compila todos los ejemplos
	@echo "$(GREEN)🔨 Compilando ejemplos...$(NC)"
	@echo "$(BLUE)  → Compilando client example...$(NC)"
	cd examples/client && go build -o client_example client_example.go
	@echo "$(BLUE)  → Compilando server example...$(NC)"
	cd examples/server && go build -o server_example server_example.go
	@echo "$(GREEN)✅ Ejemplos compilados$(NC)"

.PHONY: test-examples
test-examples: ## Ejecuta tests de los ejemplos
	@echo "$(GREEN)🧪 Testeando ejemplos...$(NC)"
	@echo "$(BLUE)  → Testeando client example...$(NC)"
	cd examples/client && go test -v ./...
	@echo "$(BLUE)  → Testeando server example...$(NC)"
	cd examples/server && go test -v ./...
	@echo "$(GREEN)✅ Tests de ejemplos completados$(NC)"

.PHONY: run-server-example
run-server-example: ## Ejecuta el ejemplo del servidor
	@echo "$(GREEN)🚀 Ejecutando server example...$(NC)"
	cd examples/server && go run server_example.go

.PHONY: run-client-example
run-client-example: ## Ejecuta el ejemplo del cliente
	@echo "$(GREEN)🚀 Ejecutando client example...$(NC)"
	cd examples/client && go run client_example.go

.PHONY: docker-up
docker-up: ## Levanta el entorno Docker para los ejemplos
	@echo "$(GREEN)🐳 Levantando entorno Docker...$(NC)"
	cd examples/server && docker-compose up -d
	@echo "$(GREEN)✅ Entorno Docker iniciado$(NC)"

.PHONY: docker-down
docker-down: ## Detiene el entorno Docker
	@echo "$(GREEN)🐳 Deteniendo entorno Docker...$(NC)"
	cd examples/server && docker-compose down
	@echo "$(GREEN)✅ Entorno Docker detenido$(NC)"

.PHONY: docker-logs
docker-logs: ## Muestra logs del entorno Docker
	@echo "$(GREEN)📋 Mostrando logs de Docker...$(NC)"
	cd examples/server && docker-compose logs -f

.PHONY: clean-examples
clean-examples: ## Limpia binarios de ejemplos
	@echo "$(GREEN)🧹 Limpiando ejemplos...$(NC)"
	rm -f examples/client/client_example
	rm -f examples/server/server_example
	@echo "$(GREEN)✅ Ejemplos limpiados$(NC)"

# Target por defecto
.DEFAULT_GOAL := help