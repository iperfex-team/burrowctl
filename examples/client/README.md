# 🐇 burrowctl Client Examples

Esta carpeta contiene ejemplos de uso del cliente burrowctl para los **tres tipos de comandos** soportados.

## 🚀 Quick Start

### 1. SQL Example
```bash
cd sql-example/
go run main.go
# Con query personalizada:
go run main.go "SELECT * FROM products"
```

### 2. Function Example
```bash
cd function-example/
go run main.go
# Con función personalizada:
go run main.go myCustomFunction
```

### 3. Command Example
```bash
cd command-example/
go run main.go
# Con comando personalizado:
go run main.go "ps aux"
```

## 📁 Estructura Consistente

Cada ejemplo sigue la misma estructura:

```
example-directory/
├── main.go        # Código del ejemplo
├── go.mod         # Dependencias Go
└── go.sum         # Checksums de dependencias
```

## 🔧 Configuración

Todos los ejemplos usan las mismas credenciales:
- **Device ID**: `fd1825ec5a7b63f3fa2be9e04154d3b16f676663ba38e23d4ffafa7b0df29efb`
- **RabbitMQ**: `amqp://burrowuser:burrowpass123@localhost:5672/`
- **Timeout**: `5s`
- **Debug**: `true`

## 📚 Documentación Completa

Ver `EXAMPLES.md` para documentación detallada sobre:
- Tipos de comandos soportados
- Sintaxis de cada tipo
- Ejemplos de uso avanzados
- Estructura del proyecto
- Tips y mejores prácticas

