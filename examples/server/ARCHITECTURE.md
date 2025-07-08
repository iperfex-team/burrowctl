# Arquitectura del Sistema de Funciones

## Separación Limpia de Responsabilidades

### 🎯 Problema Solucionado

Anteriormente, las funciones de ejemplo estaban mezcladas con el código core del servidor en `server/server.go`. Esto causaba:

- **Acoplamiento**: El servidor tenía dependencias de funciones de ejemplo
- **Código poco limpio**: Funciones de ejemplo en el core del sistema
- **Mantenimiento complicado**: Cambios en ejemplos afectaban el server core

### ✅ Solución Implementada

Ahora tenemos una **arquitectura limpia** con separación clara:

```
server/
├── server.go              # Core del servidor (LIMPIO)
└── (sin funciones de ejemplo)

examples/server/
├── server_example.go       # Funciones de ejemplo + registro
└── (todas las funciones de ejemplo)
```

## Funcionamiento del Sistema

### 1. Core del Servidor (`server/server.go`)

```go
// El servidor tiene un registry vacío por defecto
type Handler struct {
    // ... otros campos
    functionRegistry map[string]interface{}  // Vacío por defecto
}

// Métodos para registro dinámico
func (h *Handler) RegisterFunction(name string, function interface{})
func (h *Handler) RegisterFunctions(functions map[string]interface{})
func (h *Handler) GetRegisteredFunctions() []string
```

### 2. Ejemplo del Servidor (`examples/server/server_example.go`)

```go
// Todas las funciones de ejemplo están aquí
func returnString() string { return "Hola mundo" }
func lengthOfString(s string) int { return len(s) }
// ... 16 funciones de ejemplo

// Registro de funciones
func registerExampleFunctions(h *server.Handler) {
    functions := map[string]interface{}{
        "returnString": returnString,
        "lengthOfString": lengthOfString,
        // ... todas las funciones
    }
    h.RegisterFunctions(functions)
}

// En main()
func runServer() {
    h := server.NewHandler(...)
    registerExampleFunctions(h)  // Registro dinámico
    h.Start(ctx)
}
```

## Beneficios de la Nueva Arquitectura

### 🏗️ Separación de Responsabilidades
- **Server Core**: Solo maneja infraestructura (AMQP, MySQL, HTTP)
- **Ejemplos**: Contienen lógica de negocio específica

### 🔧 Extensibilidad
- Fácil agregar nuevas funciones sin tocar el core
- Múltiples ejemplos pueden registrar sus propias funciones
- Sistema de plugins natural

### 🧪 Testabilidad
- Server core se puede testear independientemente
- Funciones de ejemplo se pueden testear por separado
- Mocks más fáciles de crear

### 📦 Mantenibilidad
- Código core más limpio y enfocado
- Cambios en ejemplos no afectan el servidor
- Mejor organización del código

## Uso del Sistema

### Para Desarrolladores del Core
```go
// El servidor está limpio - no hay funciones hardcodeadas
h := server.NewHandler(deviceID, amqpURL, mysqlDSN, mode, pool)
// Registry está vacío por defecto
```

### Para Desarrolladores de Ejemplos
```go
// Registrar funciones personalizadas
h.RegisterFunction("myFunction", myFunction)

// O registrar múltiples funciones
functions := map[string]interface{}{
    "func1": func1,
    "func2": func2,
}
h.RegisterFunctions(functions)
```

### Para Usuarios Finales
```bash
# El ejemplo funciona igual que antes
./server_example -list        # Lista funciones disponibles
./server_example -functions   # Documentación completa
./server_example             # Ejecutar servidor
```

## Migración Completada

### ❌ Antes (Problemático)
```go
// En server/server.go
func getFunctionByName(name string) reflect.Value {
    functions := map[string]interface{}{
        "returnString": returnString,  // Función hardcodeada
        "lengthOfString": lengthOfString,  // Función hardcodeada
        // ... más funciones hardcodeadas
    }
}

// Funciones mezcladas con código core
func returnString() string { return "Hola mundo" }
func lengthOfString(s string) int { return len(s) }
```

### ✅ Después (Limpio)
```go
// En server/server.go
func (h *Handler) getFunctionByName(name string) reflect.Value {
    if fn, exists := h.functionRegistry[name]; exists {
        return reflect.ValueOf(fn)
    }
    return reflect.Value{}
}
// NO hay funciones hardcodeadas

// En examples/server/server_example.go
func returnString() string { return "Hola mundo" }
func lengthOfString(s string) int { return len(s) }
// Funciones en su lugar correcto
```

## Conclusión

La nueva arquitectura mantiene el **core del servidor limpio** mientras permite **extensibilidad total**. Los ejemplos pueden registrar sus funciones dinámicamente sin ensuciar el código core.

Esta separación facilita:
- **Mantenimiento** del código core
- **Extensiones** por parte de usuarios
- **Testing** independiente
- **Distribución** de responsabilidades

¡El sistema está ahora bien arquitecturado y listo para producción! 🚀 