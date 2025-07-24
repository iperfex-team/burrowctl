# 🚀 Servidor Enterprise Completo - Burrowctl

Este es un ejemplo completo de un servidor burrowctl con todas las características enterprise habilitadas, usando **configuración programática** directamente en el código.

## ✨ Características

- **🔄 Query Caching**: Cache de alto rendimiento con TTL configurable
- **🛡️ SQL Validation**: Validación de seguridad multicapa con detección de inyección
- **⚡ Worker Pool**: Pool de workers configurable para procesamiento concurrente
- **🚦 Rate Limiting**: Limitación de velocidad por IP de cliente con soporte de burst
- **🔗 Connection Pooling**: Pool de conexiones de base de datos optimizado
- **📊 Comprehensive Monitoring**: Métricas de rendimiento y seguridad en tiempo real
- **⚙️ Heartbeat Management**: Monitoreo de conectividad de clientes
- **🎯 Configuración Programática**: Todas las configuraciones directamente en el código

## 🚀 Uso Rápido

```bash
# Compilar y ejecutar
go build -o server main.go
./server

# O ejecutar directamente
go run main.go
```

## ⚙️ Configuración

El servidor usa **configuración programática** directamente en el código. Todas las configuraciones están definidas en `main.go`:

```go
config := &server.ServerConfig{
    // Configuración básica
    DeviceID: "my-custom-device",
    AMQPURL:  "amqp://burrowuser:burrowpass123@localhost:5672/",
    MySQLDSN: "burrowuser:burrowpass123@tcp(localhost:3306)/burrowdb",
    
    // Configuración de cache
    CacheEnabled: true,
    CacheSize:    3000,
    CacheTTL:     20 * time.Minute,
    CacheCleanup: 8 * time.Minute,
    
    // Configuración de validación SQL
    ValidationEnabled: true,
    StrictMode:        false,
    AllowDDL:          false,
    AllowDML:          true,
    AllowStoredProcs:  false,
    MaxQueryLength:    8000,
    LogViolations:     true,
    
    // Configuración de rendimiento
    Workers:   30,
    QueueSize: 1500,
    RateLimit: 150,
    BurstSize: 300,
    
    // Configuración de base de datos
    PoolIdle:     30,
    PoolOpen:     80,
    ConnLifetime: 12 * time.Minute,
    
    // Configuración de monitoreo
    MonitoringEnabled:  true,
    MonitoringInterval: 45 * time.Second,
    
    // Configuración de heartbeat
    HeartbeatEnabled:      true,
    HeartbeatInterval:     15 * time.Second,
    HeartbeatTimeout:      5 * time.Second,
    HeartbeatMaxMissed:    5,
    HeartbeatCleanup:      1 * time.Minute,
    HeartbeatMaxClientAge: 2 * time.Minute,
    
    // Configuración de reconexión
    ReconnectEnabled:           true,
    ReconnectMaxAttempts:       15,
    ReconnectInitialInterval:   2 * time.Second,
    ReconnectMaxInterval:       120 * time.Second,
    ReconnectBackoffMultiplier: 1.5,
    ReconnectResetInterval:     10 * time.Minute,
}
```

## 📋 Referencia de Configuración

### Configuración de Cache
- `CacheEnabled`: Habilitar/deshabilitar cache de queries
- `CacheSize`: Máximo número de queries cacheadas
- `CacheTTL`: Tiempo de vida del cache
- `CacheCleanup`: Intervalo de limpieza

### Validación SQL
- `ValidationEnabled`: Habilitar validación SQL
- `StrictMode`: Modo de validación estricta
- `AllowDDL`: Permitir Data Definition Language
- `AllowDML`: Permitir Data Manipulation Language
- `AllowStoredProcs`: Permitir stored procedures
- `MaxQueryLength`: Longitud máxima de query
- `LogViolations`: Registrar violaciones

### Rendimiento
- `Workers`: Número de goroutines workers
- `QueueSize`: Tamaño de cola de workers
- `RateLimit`: Límite de velocidad por IP de cliente
- `BurstSize`: Tamaño de burst para rate limiting

### Base de Datos
- `PoolIdle`: Máximo de conexiones idle
- `PoolOpen`: Máximo de conexiones abiertas
- `ConnLifetime`: Tiempo de vida de conexión

### Monitoreo
- `MonitoringEnabled`: Habilitar monitoreo periódico
- `MonitoringInterval`: Intervalo de reporte de monitoreo

### Heartbeat
- `HeartbeatEnabled`: Habilitar sistema de heartbeat
- `HeartbeatInterval`: Intervalo entre heartbeats
- `HeartbeatTimeout`: Timeout para respuesta
- `HeartbeatMaxMissed`: Máximo de heartbeats perdidos
- `HeartbeatCleanup`: Intervalo de limpieza
- `HeartbeatMaxClientAge`: Edad máxima del cliente

### Configuración de Reconexión
- `ReconnectEnabled`: Habilitar reconexión automática del cliente
- `ReconnectMaxAttempts`: Máximo número de intentos de reconexión
- `ReconnectInitialInterval`: Intervalo inicial entre intentos
- `ReconnectMaxInterval`: Intervalo máximo entre intentos
- `ReconnectBackoffMultiplier`: Multiplicador para backoff exponencial
- `ReconnectResetInterval`: Intervalo para resetear el backoff

## 🔧 Funciones de Monitoreo

El servidor registra automáticamente las siguientes funciones de monitoreo:

- `getCacheStats()`: Estadísticas de rendimiento del cache
- `getValidationStats()`: Métricas de validación SQL y seguridad
- `getSystemStatus()`: Estado general del sistema
- `getPerformanceMetrics()`: Análisis de rendimiento
- `clearAllCaches()`: Limpieza administrativa de caches

## 🛡️ Características de Seguridad

### Detección de Inyección SQL
- Detección basada en patrones
- Lista blanca/negra de comandos
- Validación estructural de queries
- Validación de parámetros
- Evaluación de nivel de riesgo

### Rate Limiting
- Limitación de velocidad por IP de cliente
- Capacidad de burst configurable
- Limpieza automática de datos de rate limit

### Seguridad de Base de Datos
- Pool de conexiones con gestión de tiempo de vida
- Soporte para prepared statements
- Aislamiento de transacciones

## ⚡ Optimizaciones de Rendimiento

### Query Caching
- Cache LRU con tamaño configurable
- Expiración basada en TTL
- Limpieza automática del cache
- Estadísticas de hit/miss del cache

### Worker Pool
- Número de workers configurable
- Cola acotada con protección de overflow
- Soporte para shutdown graceful
- Balanceo de carga entre workers

### Gestión de Conexiones
- Pool de conexiones con límites idle/open
- Gestión de tiempo de vida de conexiones
- Health checking y reconexión
- Drenado graceful de conexiones

## 📊 Salida de Monitoreo

El servidor proporciona monitoreo completo cada 45 segundos (configurable):

```
📊 COMPREHENSIVE SYSTEM REPORT - 14:30:15
============================================================
🏢 System Overview:
  Uptime: 2h45m30s

📈 Cache Performance:
  Total Requests: 15,432
  Cache Hits: 12,345
  Cache Misses: 3,087
  Hit Ratio: 80.00%
  Current Size: 1,234 entries
  Evictions: 45
  Expirations: 123

🛡️ Security & Validation:
  Total Queries: 14,567
  Valid Queries: 14,445
  Blocked Queries: 122
  Injection Attempts: 5
  Command Violations: 87
  Structure Violations: 30
  Block Rate: 0.84%
  Injection Rate: 0.03%
  Security Level: LOW

⚡ Performance Summary:
  Cache Efficiency: 80.00%
  Validation Efficiency: 99.16%
============================================================
```

## 🐳 Uso con Docker

Puedes ejecutar el servidor completo con Docker usando el archivo `docker-compose-full.yml`:

```bash
# Iniciar todos los servicios
docker-compose -f docker-compose-full.yml up -d

# Ver logs
docker-compose -f docker-compose-full.yml logs -f

# Detener servicios
docker-compose -f docker-compose-full.yml down
```

El archivo Docker Compose incluye:
- **MariaDB**: Base de datos MySQL
- **RabbitMQ**: Message broker con interfaz de administración
- **Servidor Burrowctl**: Con todas las configuraciones enterprise

## 🎯 Uso de Clientes de Ejemplo

Una vez que el servidor esté ejecutándose, puedes conectar clientes para probar las características:

```bash
# Probar funcionalidad básica
cd ../../../client/command-example
go run main.go

# Probar SQL con cache
cd ../sql-example
go run main.go

# Probar validación
cd ../validation-example
go run main.go
```

## ✅ Ventajas de la Configuración Programática

1. **🎯 Control Total**: Tienes control completo sobre cada aspecto de la configuración
2. **📝 Versionado**: La configuración está versionada con tu código
3. **🔒 Type Safety**: Validación de valores de configuración en tiempo de compilación
4. **💻 Soporte de IDE**: Autocompletado completo y detección de errores
5. **⚙️ Fácil Personalización**: Modifica cualquier configuración directamente en el código
6. **🌍 Consistencia**: Misma configuración en todos los despliegues

## 🎓 Aprendizaje

Este ejemplo te enseña:

1. **Patrones de Diseño**: Factory Pattern, Configuration Pattern
2. **Separación de Responsabilidades**: Cada componente tiene una función específica
3. **Encapsulación**: La lógica compleja está oculta en la librería
4. **Buenas Prácticas**: Código limpio y mantenible
5. **Arquitectura Profesional**: Estructura escalable y profesional

¡Ahora tienes control total sobre la configuración del servidor con una arquitectura moderna y profesional! 🎉

## 🔄 **Configuración de Reconexión del Cliente**

### **DSN del Cliente con Reconexión Personalizada**

El cliente ahora soporta configuración de reconexión a través del DSN:

```go
// DSN con configuración de reconexión personalizada
dsn := "deviceID=my-device&amqp_uri=amqp://user:pass@localhost:5672/&timeout=30s&debug=true&reconnect_enabled=true&reconnect_max_attempts=0&reconnect_initial_interval=2s&reconnect_max_interval=120s&reconnect_backoff_multiplier=1.5&reconnect_reset_interval=10m"

db, err := sql.Open("rabbitsql", dsn)
```

### **Parámetros de Reconexión del Cliente**

- `reconnect_enabled`: Habilitar reconexión automática (default: true)
- `reconnect_max_attempts`: Máximo número de intentos (default: 10, 0 = infinito)
- `reconnect_initial_interval`: Intervalo inicial entre intentos (default: 1s)
- `reconnect_max_interval`: Intervalo máximo entre intentos (default: 60s)
- `reconnect_backoff_multiplier`: Multiplicador de backoff (default: 2.0)
- `reconnect_reset_interval`: Intervalo para resetear backoff (default: 5m)

**⚠️ IMPORTANTE**: Para reconexión infinita (recomendado para producción), usa `reconnect_max_attempts=0`

### **Configuración por Defecto del Cliente**

```go
// Configuración automática que ya viene incluida
ReconnectConfig{
    Enabled:           true,           // Reconexión habilitada
    MaxAttempts:       10,             // Hasta 10 intentos
    InitialInterval:   1 * time.Second, // Empieza con 1 segundo
    MaxInterval:       60 * time.Second, // Máximo 60 segundos
    BackoffMultiplier: 2.0,            // Duplica el tiempo cada intento
    ResetInterval:     5 * time.Minute, // Resetea después de 5 min de éxito
}
```

### **Ejemplo Completo de Cliente con Reconexión**

```go
package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lordbasex/burrowctl/client" // Registrar el driver
)

func main() {
	// DSN con configuración de reconexión personalizada
	dsn := "deviceID=my-custom-device&amqp_uri=amqp://burrowuser:burrowpass123@localhost:5672/&timeout=30s&debug=true&reconnect_enabled=true&reconnect_max_attempts=0&reconnect_initial_interval=2s&reconnect_max_interval=120s&reconnect_backoff_multiplier=1.5&reconnect_reset_interval=10m"

	// Abrir conexión con reconexión automática
	db, err := sql.Open("rabbitsql", dsn)
	if err != nil {
		log.Fatal("Error al abrir conexión:", err)
	}
	defer db.Close()

	// Probar conexión
	if err := db.Ping(); err != nil {
		log.Fatal("Error al hacer ping:", err)
	}

	fmt.Println("✅ Conexión establecida con reconexión automática configurada")

	// Ejecutar queries - la reconexión será automática si se pierde la conexión
	rows, err := db.Query("SELECT 1 as numero, 'Hola Mundo' as mensaje")
	if err != nil {
		log.Fatal("Error en query:", err)
	}
	defer rows.Close()

	// Procesar resultados
	for rows.Next() {
		var numero int
		var mensaje string
		if err := rows.Scan(&numero, &mensaje); err != nil {
			log.Fatal("Error al escanear:", err)
		}
		fmt.Printf("Resultado: %d - %s\n", numero, mensaje)
	}

	fmt.Println("🎉 Cliente funcionando con reconexión automática!")
}

## 🎉 **Resumen de la Implementación de Reconexión**

### ✅ **Lo que se ha implementado:**

1. **🔧 Configuración del Servidor (`ServerConfig`)**:
   - Campos de reconexión agregados al `ServerConfig`
   - Variables de entorno para configuración
   - Valores por defecto sensatos
   - Método `ToReconnectConfig()` para conversión

2. **🔌 Configuración del Cliente (DSN)**:
   - Parámetros de reconexión en el DSN del cliente
   - Configuración personalizable por conexión
   - Valores por defecto automáticos
   - Documentación completa

3. **📚 Documentación**:
   - README actualizado con ejemplos
   - Variables de entorno documentadas
   - Ejemplos de uso completos
   - Configuración por defecto explicada

### 🚀 **Cómo usar la configuración:**

#### **Servidor (Programático)**:
```go
config := &server.ServerConfig{
    // ... otras configuraciones ...
    
    // Configuración de reconexión
    ReconnectEnabled:           true,
    ReconnectMaxAttempts:       15,
    ReconnectInitialInterval:   2 * time.Second,
    ReconnectMaxInterval:       120 * time.Second,
    ReconnectBackoffMultiplier: 1.5,
    ReconnectResetInterval:     10 * time.Minute,
}
```

#### **Servidor (Variables de Entorno)**:
```bash
export RECONNECT_ENABLED=true
export RECONNECT_MAX_ATTEMPTS=15
export RECONNECT_INITIAL_INTERVAL=2s
export RECONNECT_MAX_INTERVAL=120s
export RECONNECT_BACKOFF_MULTIPLIER=1.5
export RECONNECT_RESET_INTERVAL=10m
```

#### **Cliente (DSN)**:
```go
dsn := "deviceID=my-device&amqp_uri=amqp://user:pass@localhost:5672/&reconnect_max_attempts=15&reconnect_initial_interval=2s&reconnect_max_interval=120s&reconnect_backoff_multiplier=1.5&reconnect_reset_interval=10m"
```

### 🎯 **Beneficios:**

- **✅ Control Total**: Configuración completa de reconexión
- **✅ Flexibilidad**: Múltiples métodos de configuración
- **✅ Valores por Defecto**: Funciona sin configuración
- **✅ Retrocompatibilidad**: Código existente sigue funcionando
- **✅ Documentación**: Guías completas de uso

¡Ahora tienes control total sobre la reconexión tanto en el servidor como en el cliente! 🎉

## ✅ **Capacidades de Reconexión Automática**

### 🔄 **Cliente (Lado del Nodo)**

La librería del cliente tiene un **`ConnectionManager`** muy sofisticado que maneja reconexión automática:

#### **Características de Reconexión:**
- **✅ Reconexión automática** con backoff exponencial
- **✅ Configuración personalizable** de reintentos
- **✅ Monitoreo de salud** de conexión
- **✅ Callbacks** para eventos de conexión/desconexión
- **✅ Estadísticas** en tiempo real

#### **Configuración por Defecto:**
```go
// Configuración automática que ya viene incluida
ReconnectConfig{
    Enabled:           true,           // Reconexión habilitada
    MaxAttempts:       10,             // Hasta 10 intentos
    InitialInterval:   1 * time.Second, // Empieza con 1 segundo
    MaxInterval:       60 * time.Second, // Máximo 60 segundos
    BackoffMultiplier: 2.0,            // Duplica el tiempo cada intento
    ResetInterval:     5 * time.Minute, // Resetea después de 5 min de éxito
}
```

### 🛡️ **Servidor (Lado del Datacenter)**

El servidor también tiene capacidades robustas:

#### **Gestión de Conexiones:**
- **✅ Manejo de desconexiones** de RabbitMQ
- **✅ Reconexión automática** del servidor
- **✅ Heartbeat management** para detectar problemas
- **✅ Worker pool** que maneja reconexiones

## 🎓 Aprendizaje

Este ejemplo te enseña:

1. **Patrones de Diseño**: Factory Pattern, Configuration Pattern
2. **Separación de Responsabilidades**: Cada componente tiene una función específica
3. **Encapsulación**: La lógica compleja está oculta en la librería
4. **Buenas Prácticas**: Código limpio y mantenible
5. **Arquitectura Profesional**: Estructura escalable y profesional

¡La librería está preparada para manejar exactamente este escenario! 🎉

# Heartbeat configuration
export HEARTBEAT_ENABLED=true
export HEARTBEAT_INTERVAL=15s
export HEARTBEAT_TIMEOUT=5s
export HEARTBEAT_MAX_MISSED=5
export HEARTBEAT_CLEANUP=1m
export HEARTBEAT_MAX_CLIENT_AGE=2m

# Reconnection configuration
export RECONNECT_ENABLED=true
export RECONNECT_MAX_ATTEMPTS=15
export RECONNECT_INITIAL_INTERVAL=2s
export RECONNECT_MAX_INTERVAL=120s
export RECONNECT_BACKOFF_MULTIPLIER=1.5
export RECONNECT_RESET_INTERVAL=10m