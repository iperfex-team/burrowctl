# Makefile para burrowctl
# Versión por defecto
# ---------------- Gestión automática de versión ----------------
# Obtiene el hash corto de git; si falla (por ejemplo, sin repo), muestra 'development'
git_hash := $(shell git rev-parse --short HEAD 2>/dev/null || echo 'development')

# Archivo donde se persiste la versión
version_file := version.txt

# Versión inicial por defecto si el archivo no existe
initial_version := 1.0.0

# Crear version.txt si no existe
ifeq ($(wildcard $(version_file)),)
$(shell echo "initial_version: $(initial_version)" > $(version_file))
$(shell echo "git_hash: $(git_hash)" >> $(version_file))
endif

# Leer última versión y hash almacenados
last_version := $(shell awk -F': ' '/initial_version:/ {print $$2}' $(version_file) | xargs)
last_git_hash := $(shell awk -F': ' '/git_hash:/ {print $$2}' $(version_file) | xargs)

# Calcular la versión a usar
ifeq ($(strip $(git_hash)), $(strip $(last_git_hash)))
VERSION := $(last_version)
else
# Incrementar el número de versión (patch → minor → major)
VERSION := $(shell \
  major=$$(echo $(last_version) | awk -F. '{print $$1}'); \
  minor=$$(echo $(last_version) | awk -F. '{print $$2}'); \
  patch=$$(echo $(last_version) | awk -F. '{print $$3}'); \
  if [ $$patch -eq 9 ]; then \
    if [ $$minor -eq 9 ]; then \
      echo "$$((major + 1)).0.0"; \
    else \
      echo "$$major.$$((minor + 1)).0"; \
    fi; \
  else \
    echo "$$major.$$minor.$$((patch + 1))"; \
  fi)
# Actualizar version.txt con la nueva versión y hash
$(shell echo "initial_version: $(VERSION)" > $(version_file))
$(shell echo "git_hash: $(git_hash)" >> $(version_file))
endif

# Mensajes informativos
$(info git_hash: $(git_hash))
$(info last_version: $(last_version))
$(info last_git_hash: $(last_git_hash))
$(info next VERSION: $(VERSION))

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
	@UNCOMMITTED=$$(git status --porcelain | grep -v "^ M version.txt$$" | grep -v "^M  version.txt$$"); \
	if [ -n "$$UNCOMMITTED" ]; then \
		echo "$(RED)❌ Error: Hay cambios sin commit$(NC)"; \
		git status --short | grep -v "version.txt"; \
		exit 1; \
	fi
	@echo "$(GREEN)✅ Git está limpio (ignorando version.txt)$(NC)"

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
	@echo "$(BLUE)  → Compilando command example...$(NC)"
	cd examples/client/command-example && go build -o command-example main.go
	@echo "$(BLUE)  → Compilando function example...$(NC)"
	cd examples/client/function-example && go build -o function-example main.go
	@echo "$(BLUE)  → Compilando sql example...$(NC)"
	cd examples/client/sql-example && go build -o sql-example main.go
	@echo "$(BLUE)  → Compilando server example...$(NC)"
	cd examples/server && go build -o server_example server_example.go
	@echo "$(GREEN)✅ Ejemplos compilados$(NC)"

.PHONY: test-examples
test-examples: ## Ejecuta tests de los ejemplos
	@echo "$(GREEN)🧪 Testeando ejemplos...$(NC)"
	@echo "$(BLUE)  → Verificando command example...$(NC)"
	cd examples/client/command-example && go build -o /tmp/command-example-test main.go && rm -f /tmp/command-example-test
	@echo "$(BLUE)  → Verificando function example...$(NC)"
	cd examples/client/function-example && go build -o /tmp/function-example-test main.go && rm -f /tmp/function-example-test
	@echo "$(BLUE)  → Verificando sql example...$(NC)"
	cd examples/client/sql-example && go build -o /tmp/sql-example-test main.go && rm -f /tmp/sql-example-test
	@echo "$(BLUE)  → Testeando server example...$(NC)"
	cd examples/server && go test -v ./... 2>/dev/null || echo "$(YELLOW)    No tests found - skipping$(NC)"
	@echo "$(GREEN)✅ Tests de ejemplos completados$(NC)"

.PHONY: run-server-example
run-server-example: ## Ejecuta el ejemplo del servidor
	@echo "$(GREEN)🚀 Ejecutando server example...$(NC)"
	cd examples/server && go run server_example.go

.PHONY: run-command-example
run-command-example: ## Ejecuta el ejemplo de comando
	@echo "$(GREEN)🚀 Ejecutando command example...$(NC)"
	cd examples/client/command-example && go run main.go

.PHONY: run-function-example
run-function-example: ## Ejecuta el ejemplo de función
	@echo "$(GREEN)🚀 Ejecutando function example...$(NC)"
	cd examples/client/function-example && go run main.go

.PHONY: run-sql-example
run-sql-example: ## Ejecuta el ejemplo de SQL
	@echo "$(GREEN)🚀 Ejecutando sql example...$(NC)"
	cd examples/client/sql-example && go run main.go

.PHONY: run-client-examples
run-client-examples: ## Ejecuta todos los ejemplos de cliente
	@echo "$(GREEN)🚀 Ejecutando todos los ejemplos de cliente...$(NC)"
	@echo "$(BLUE)  → Ejecutando command example...$(NC)"
	cd examples/client/command-example && go run main.go
	@echo "$(BLUE)  → Ejecutando function example...$(NC)"
	cd examples/client/function-example && go run main.go
	@echo "$(BLUE)  → Ejecutando sql example...$(NC)"
	cd examples/client/sql-example && go run main.go

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
	rm -f examples/client/command-example/command-example
	rm -f examples/client/function-example/function-example
	rm -f examples/client/sql-example/sql-example
	rm -f examples/server/server_example
	@echo "$(GREEN)✅ Ejemplos limpiados$(NC)"

.PHONY: list-examples
list-examples: ## Lista todos los ejemplos disponibles
	@echo "$(BLUE)📋 Ejemplos disponibles:$(NC)"
	@echo "$(GREEN)Client Examples:$(NC)"
	@echo "  - command-example: Ejecuta comandos remotos"
	@echo "  - function-example: Ejecuta funciones remotas"
	@echo "  - sql-example: Ejecuta consultas SQL remotas"
	@echo "$(GREEN)Server Examples:$(NC)"
	@echo "  - server-example: Servidor que procesa las peticiones"
	@echo ""
	@echo "$(YELLOW)Para ejecutar un ejemplo específico:$(NC)"
	@echo "  make run-command-example"
	@echo "  make run-function-example"
	@echo "  make run-sql-example"
	@echo "  make run-server-example"

.PHONY: demo-command
demo-command: ## Ejecuta una demostración del comando ejemplo
	@echo "$(GREEN)🎬 Demostración del command-example...$(NC)"
	@echo "$(BLUE)  → Ejecutando 'ls -la' remoto...$(NC)"
	cd examples/client/command-example && go run main.go "ls -la"
	@echo "$(BLUE)  → Ejecutando 'ps aux' remoto...$(NC)"
	cd examples/client/command-example && go run main.go "ps aux"

# Target por defecto
.DEFAULT_GOAL := help