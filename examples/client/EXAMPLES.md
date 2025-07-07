# 🐇 burrowctl Client Examples

Este directorio contiene ejemplos de cómo usar el cliente burrowctl con los **tres tipos de comandos** soportados:

## 📊 Tipos de Comandos Soportados

### 1. **SQL Queries** (`type: sql`)
- **Descripción**: Ejecuta consultas SQL directamente en la base de datos remota
- **Sintaxis**: Query normal de SQL
- **Ejemplo**: `SELECT * FROM users`

### 2. **Functions** (`type: function`)
- **Descripción**: Ejecuta funciones personalizadas en el servidor remoto
- **Sintaxis**: `FUNCTION:nombreFuncion`
- **Ejemplo**: `FUNCTION:processData`

### 3. **Commands** (`type: command`)
- **Descripción**: Ejecuta comandos del sistema en el servidor remoto
- **Sintaxis**: `COMMAND:comando`
- **Ejemplo**: `COMMAND:ls -la`

---

## 🚀 Ejemplos Disponibles

### 1. SQL Example
```bash
cd examples/client/sql-example/
go run main.go
# O con query personalizada:
go run main.go "SELECT * FROM products"
```

### 2. Function Example
```bash
cd examples/client/function-example/
go run main.go
# O con función personalizada:
go run main.go myCustomFunction
```

### 3. Command Example
```bash
cd examples/client/command-example/
go run main.go
# O con comando personalizado:
go run main.go "ps aux"
```

---

## 🔧 Cómo Funciona la Detección Automática

El cliente **detecta automáticamente** el tipo de comando basándose en el prefijo:

```go
// El cliente internamente hace esto:
func parseCommand(query string) (cmdType string, actualQuery string) {
    if strings.HasPrefix(query, "FUNCTION:") {
        return "function", query[9:]  // Remueve "FUNCTION:"
    }
    if strings.HasPrefix(query, "COMMAND:") {
        return "command", query[8:]   // Remueve "COMMAND:"
    }
    return "sql", query              // Por defecto es SQL
}
```

---

## 💡 Tips de Uso

### Para SQL:
```go
rows, err := db.Query("SELECT id, name FROM users")
```

### Para Functions:
```go
rows, err := db.Query("FUNCTION:processData")
// O con argumentos (pasados como params):
rows, err := db.Query("FUNCTION:processUser", userID)
```

### Para Commands:
```go
rows, err := db.Query("COMMAND:ls -la")
// Para comandos con espacios, usa comillas:
rows, err := db.Query("COMMAND:ps aux | grep mysql")
```

---

## 🏗️ Estructura del Proyecto

```
examples/client/
├── sql-example/               # Ejemplo SQL básico
│   ├── main.go
│   └── go.mod
├── function-example/          # Ejemplo de funciones
│   ├── main.go
│   └── go.mod
├── command-example/           # Ejemplo de comandos
│   ├── main.go
│   └── go.mod
├── go.mod                     # Módulo principal
└── EXAMPLES.md               # Esta documentación
```

---

## 🛠️ Configuración

Todos los ejemplos usan la misma configuración DSN:

```go
dsn := "deviceID=fd1825ec5a7b63f3fa2be9e04154d3b16f676663ba38e23d4ffafa7b0df29efb&amqp_uri=amqp://burrowuser:burrowpass123@localhost:5672/&timeout=5s&debug=true"
```

### Parámetros DSN:
- `deviceID`: ID único del dispositivo remoto
- `amqp_uri`: URI de conexión a RabbitMQ
- `timeout`: Tiempo de espera para respuestas
- `debug`: Habilitar logs de depuración

---

## 📝 Respuestas del Servidor

### SQL Response:
```json
{
  "columns": ["id", "name"],
  "rows": [[1, "Alice"], [2, "Bob"]],
  "error": ""
}
```

### Function Response (Mock):
```json
{
  "columns": ["message"],
  "rows": [["function executed (mock)"]],
  "error": ""
}
```

### Command Response (Mock):
```json
{
  "columns": ["message"],
  "rows": [["command executed (mock)"]],
  "error": ""
}
```

---

## 🎯 Próximos Pasos

1. **Implementar funciones reales** en el servidor
2. **Implementar comandos reales** en el servidor
3. **Añadir autenticación** para comandos sensibles
4. **Implementar parámetros** para funciones y comandos 