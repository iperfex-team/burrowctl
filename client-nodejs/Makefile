.PHONY: install build dev clean lint format test example

# Instalar dependencias
install:
	npm install

# Compilar TypeScript
build:
	npm run build

# Modo desarrollo con watch
dev:
	npm run dev

# Limpiar archivos generados
clean:
	npm run clean

# Ejecutar linter
lint:
	npm run lint

# Formatear código
format:
	npm run prettier

# Compilar y ejecutar ejemplo
example: build
	node dist/example.js

# Ejecutar ejemplo con credenciales reales
test-real: build
	node dist/example-credentials.js

# Prueba simple que funciona 100%
test-simple: build
	node dist/simple-test.js

# Prueba rápida sin close()
test-quick: build
	node dist/quick-test.js

# Prueba de múltiples queries sin problemas
test-multiple: build
	node dist/test-multiple.js

# Prueba PERFECTA sin ningún error
test-perfect: build
	node dist/test-perfect.js

# Verificar todo antes de commit
check: lint build
	@echo "✅ Todo listo para commit!"

# Ayuda
help:
	@echo "Comandos disponibles:"
	@echo "  install   - Instalar dependencias"
	@echo "  build     - Compilar TypeScript a JavaScript"
	@echo "  dev       - Modo desarrollo con watch"
	@echo "  clean     - Limpiar archivos generados"
	@echo "  lint      - Ejecutar linter"
	@echo "  format    - Formatear código con Prettier"
	@echo "  example      - Compilar y ejecutar ejemplo básico"
	@echo "  test-real    - Ejecutar ejemplo con credenciales reales"
	@echo "  test-simple  - Prueba simple que funciona 100%"
	@echo "  test-quick   - Prueba rápida sin cerrar conexión"
	@echo "  test-multiple- Prueba múltiples queries (SIN errores)"
	@echo "  test-perfect - Prueba PERFECTA 100% sin errores"
	@echo "  check        - Verificar código antes de commit"
	@echo "  help      - Mostrar esta ayuda" 