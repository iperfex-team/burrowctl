# 🚀 Advanced Server Example

Este ejemplo demuestra las nuevas características empresariales del servidor burrowctl con configuración completamente personalizable.

## 🎯 Características Demostradas

- **🏗️ Worker Pool**: Procesamiento concurrente configurable
- **🛡️ Rate Limiting**: Protección contra abuso por cliente
- **💾 Connection Pooling**: Optimización de conexiones de base de datos
- **⚙️ Configuración Granular**: Control total sobre todos los parámetros

## 🏁 Inicio Rápido

### Local
```bash
# Compilar el ejemplo
go build advanced_server_example.go

# Ejecutar con configuración por defecto
./advanced_server_example

# Ver configuración actual
./advanced_server_example -show-config

# Ver todas las opciones
./advanced_server_example -help
```

### Docker
```bash
# Levantar entorno completo (RabbitMQ + MariaDB + Servidor Avanzado)
docker-compose up -d

# Solo servicios (para ejecutar servidor local)
docker-compose -f docker-compose-basic.yml up -d

# Ver logs del servidor avanzado
docker-compose logs -f app-advanced

# Detener entorno
docker-compose down
```

## 📋 Configuraciones Predefinidas

### 1. Alto Rendimiento
```bash
./advanced_server_example \
  -workers=20 -queue-size=500 \
  -pool-idle=20 -pool-open=50 \
  -rate-limit=50 -burst-size=100
```

### 2. Configuración Conservadora
```bash
./advanced_server_example \
  -workers=5 -queue-size=50 \
  -pool-idle=5 -pool-open=10 \
  -rate-limit=5 -burst-size=10
```

### 3. Desarrollo Local
```bash
./advanced_server_example \
  -workers=3 -queue-size=20 \
  -pool-idle=3 -pool-open=8 \
  -rate-limit=20 -burst-size=40
```

## ⚙️ Opciones de Configuración

### Conexión
- `-device`: Device ID único
- `-amqp`: URL de RabbitMQ  
- `-mysql`: DSN de MySQL/MariaDB
- `-mode`: Modo de conexión (open/close)

### Worker Pool
- `-workers`: Número de workers (default: 10)
- `-queue-size`: Tamaño del buffer (default: 100)
- `-worker-timeout`: Timeout por tarea (default: 30s)

### Database Pool
- `-pool-idle`: Conexiones idle máximas (default: 10)
- `-pool-open`: Conexiones totales máximas (default: 20)
- `-pool-lifetime`: Tiempo de vida de conexión (default: 5m)

### Rate Limiting
- `-rate-limit`: Requests/segundo por cliente (default: 10)
- `-burst-size`: Tokens de ráfaga máximos (default: 20)
- `-rate-cleanup`: Intervalo de limpieza (default: 5m)

## 📊 Ejemplos de Salida

### Configuración Mostrada
```bash
./advanced_server_example -show-config
```

```
🔧 Current Server Configuration
===============================

📡 Connection:
   Device ID: fd1825ec5a7b63f3fa2be9e04154d3b16f676663ba38e23d4ffafa7b0df29efb
   RabbitMQ:  amqp://burrowuser:burrowpass123@rabbitmq:5672/
   MySQL:     burrowuser:burrowpass123@tcp(mariadb:3306)/burrowdb?parseTime=true
   Mode:      open

💾 Database Pool:
   Max Idle:     10 connections
   Max Open:     20 connections
   Lifetime:     5m0s

🏗️  Worker Pool:
   Workers:      10 goroutines
   Queue Size:   100 messages
   Timeout:      30s per task

🛡️  Rate Limiting:
   Rate Limit:   10 req/sec per client
   Burst Size:   20 tokens
   Cleanup:      5m0s interval

📊 Performance Estimates:
   Max Throughput: ~600 req/min (with 10 workers)
   Rate Limit Cap: 1000 req/sec (100 clients)
   DB Concurrency: 20 max connections
```

### Servidor en Ejecución
```
🚀 Advanced burrowctl Server Starting
=====================================

📋 Server Configuration:
   📱 Device ID: fd1825ec...
   🐰 RabbitMQ: amqp://burrowuser:burrowpass123@rabbitmq:5672/
   🗄️  MySQL: burrowuser:burrowpass123@tcp(mariadb:3306)/burrowdb
   🔗 Mode: open

💾 Database Pool:
   ├─ Max Idle: 10
   ├─ Max Open: 20
   └─ Lifetime: 5m0s

🏗️  Worker Pool:
   ├─ Workers: 10
   ├─ Queue Size: 100
   └─ Timeout: 30s

🛡️  Rate Limiting:
   ├─ Rate: 10 req/sec per client
   ├─ Burst: 20 tokens
   └─ Cleanup: 5m0s

🔧 Registered Functions: 18
   ├─ returnError
   ├─ returnBool
   ├─ returnInt
   ├─ returnString
   ├─ returnStruct
   ├─ heavyComputation
   └─ sleepFunction

✅ Server Capabilities:
   📊 SQL Queries - Execute remote SQL with connection pooling
   🔧 Functions - Execute typed functions with worker pool  
   ⚡ Commands - Execute system commands with timeout
   🛡️  Rate Limited - Protected against abuse
   🔄 Auto Reconnect - Client-side automatic reconnection

🎯 Performance Features Active:
   • Worker Pool: 10 concurrent message processors
   • Connection Pool: 10-20 database connections
   • Rate Limiting: 10 req/sec per client (burst: 20)
   • Prepared Statements: Client-side statement caching
   • Auto Reconnection: Client-side connection recovery

⏰ Server started at: 2024-01-15T10:30:00Z
🎯 Server is ready to accept connections...
```

### Logs en Tiempo Real
```
[server] Database pool initialized: idle=10 open=20 lifetime=5m0s
[server] Worker pool started with 10 workers, queue size: 100
[server] Queue 'fd1825ec...' declared successfully
[server] Listening on queue fd1825ec...
[server] received ip=192.168.1.100 type=sql query=SELECT * FROM users
[server] rate limit exceeded for client 192.168.1.101
[server] received ip=192.168.1.102 type=function query={"name":"returnString","params":[]}
```

## 🔍 Monitoreo y Debugging

### Métricas de Rendimiento
- **Worker Pool**: Utilización de workers y cola
- **Rate Limiting**: Requests bloqueados por cliente
- **DB Pool**: Conexiones activas vs disponibles
- **Throughput**: Mensajes procesados por segundo

### Indicadores Clave
```bash
# Logs importantes a observar:
# - "Worker pool started" - Confirmación de inicialización
# - "rate limit exceeded" - Protección activada
# - "Database pool initialized" - Pool configurado
# - "Server is ready" - Listo para conexiones
```

## 🐳 Docker Configuration

El ejemplo avanzado incluye configuración Docker optimizada:

```yaml
app-advanced:
  build:
    context: ../../..      # Raíz del proyecto
    dockerfile: examples/server/advanced/Dockerfile
  command: >
    /app/server 
      -workers=20 
      -queue-size=500 
      -rate-limit=50 
      -burst-size=100 
      -pool-idle=20 
      -pool-open=50
```

### Servicios Docker:
- **RabbitMQ**: Puerto 5672 (AMQP) y 15672 (Management UI)
- **MariaDB**: Puerto 3306 con inicialización automática
- **App-Advanced**: Servidor burrowctl optimizado con configuración empresarial

### Configuración Automática:
- **Workers**: 20 procesadores concurrentes
- **Queue**: Buffer de 500 mensajes
- **Rate Limit**: 50 requests/segundo por cliente
- **DB Pool**: 20-50 conexiones optimizadas

### Comandos Docker:
```bash
# Levantar entorno completo
docker-compose up -d

# Solo servicios (para desarrollo local)
docker-compose -f docker-compose-basic.yml up -d

# Verificar estado
docker-compose ps

# Ver logs en tiempo real
docker-compose logs -f app-advanced

# Detener servicios
docker-compose down
```

## 🧪 Testing

### Test de Carga
```bash
# Usar cliente avanzado para generar carga
cd ../client/advanced
./advanced-main -stress -concurrent=20 -requests=100
```

### Validación de Rate Limiting
```bash
# Configurar límite bajo para testing
./advanced_server_example -rate-limit=2 -burst-size=5

# Generar carga que exceda límites
cd ../client/advanced  
./advanced-main -stress -concurrent=10 -requests=50
```

### Prueba de Worker Pool
```bash
# Configurar pocos workers para observar queue
./advanced_server_example -workers=2 -queue-size=10

# Generar carga concurrente
./advanced-main -stress -concurrent=15 -requests=20
```

## 🔧 Troubleshooting

### Alta Latencia
- Aumentar `-workers` y `-queue-size`
- Optimizar `-pool-open` para más conexiones DB
- Revisar `-worker-timeout` para tareas lentas

### Memory Usage Alto
- Reducir `-pool-idle` y `-queue-size`
- Ajustar `-rate-cleanup` para limpieza más frecuente
- Usar `-mode=close` para conexiones por query

### Rate Limiting Muy Agresivo
- Aumentar `-rate-limit` y `-burst-size`
- Revisar `-rate-cleanup` para retención de buckets

### Conexiones DB Agotadas
- Aumentar `-pool-open`
- Reducir `-pool-lifetime` para rotación más rápida
- Optimizar queries del cliente

## 🚀 Producción

### Configuración Recomendada (Carga Media)
```bash
./advanced_server_example \
  -workers=15 \
  -queue-size=300 \
  -pool-idle=15 \
  -pool-open=30 \
  -rate-limit=25 \
  -burst-size=50 \
  -worker-timeout=45s \
  -pool-lifetime=10m
```

### Configuración Recomendada (Alta Carga)
```bash  
./advanced_server_example \
  -workers=50 \
  -queue-size=1000 \
  -pool-idle=25 \
  -pool-open=100 \
  -rate-limit=100 \
  -burst-size=200 \
  -worker-timeout=60s \
  -pool-lifetime=15m
```

## 📈 Optimización

1. **Monitorear** logs para identificar cuellos de botella
2. **Ajustar** workers según CPU disponible
3. **Configurar** pool DB según capacidad de base de datos
4. **Calibrar** rate limiting según patrones de uso reales
5. **Probar** bajo carga real antes de producción

## 🔗 Referencias

- Cliente Avanzado: `../client/advanced/`
- Documentación Completa: `./ADVANCED_FEATURES.md`
- Configuración Docker: `../docker-compose.yml`