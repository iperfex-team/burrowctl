# 🚀 Server Examples

Este directorio contiene ejemplos del servidor burrowctl organizados por nivel de complejidad.

## 📁 Estructura

```
examples/server/
├── basic/          # Ejemplo básico para empezar
└── advanced/       # Ejemplo empresarial con características avanzadas
```

## 🎯 Ejemplos Disponibles

### 📋 **Basic** (`basic/`)
Ejemplo fundamental del servidor burrowctl con configuración estándar.

**Características:**
- Configuración básica de servidor
- Registro de funciones de ejemplo
- Docker Compose para desarrollo
- Pool de conexiones por defecto

**Archivos:**
- `server_example.go` - Servidor básico con funciones de ejemplo
- `docker-compose.yml` - Entorno completo (RabbitMQ + MariaDB + App)
- `docker-compose-basic.yml` - Solo servicios (RabbitMQ + MariaDB)
- `init.sql` - Inicialización de base de datos
- `README.md` - Documentación del ejemplo básico

### 🚀 **Advanced** (`advanced/`)
Ejemplo empresarial con todas las características de rendimiento y configuración granular.

**Características:**
- Worker Pool configurable (5-50+ workers)
- Rate Limiting por cliente IP
- Connection Pooling optimizado
- Configuración avanzada via flags
- Métricas de rendimiento en tiempo real

**Archivos:**
- `advanced_server_example.go` - Servidor empresarial configurable
- `README.md` - Documentación detallada
- `go.mod` - Dependencias específicas

## 🏁 Inicio Rápido

### Ejemplo Básico
```bash
# Iniciar servicios
cd basic/
docker-compose up -d

# Ejecutar servidor básico
go run server_example.go

# O con Docker completo
docker-compose -f docker-compose.yml up
```

### Ejemplo Avanzado
```bash
# Servidor con configuración por defecto
cd advanced/
go run advanced_server_example.go

# Servidor de alto rendimiento
go run advanced_server_example.go \
  -workers=20 -queue-size=500 \
  -rate-limit=50 -pool-open=50

# Ver configuración actual
go run advanced_server_example.go -show-config

# Ver todas las opciones
go run advanced_server_example.go -help

# Con Docker (servidor optimizado)
docker-compose up -d

# Solo servicios (para desarrollo local)
docker-compose -f docker-compose-basic.yml up -d
```

## 🔄 Migración de Básico a Avanzado

Los ejemplos son completamente compatibles. Para migrar de básico a avanzado:

1. **Mantener** el código del servidor básico
2. **Agregar** configuración avanzada según necesidades
3. **Aprovechar** las nuevas características automáticamente

### Equivalencias:
```bash
# Básico (configuración por defecto)
go run basic/server_example.go

# Avanzado (misma configuración)
go run advanced/advanced_server_example.go \
  -workers=10 -queue-size=100 \
  -rate-limit=10 -pool-open=20
```

## 🐳 Docker Configuration

Cada ejemplo tiene su propia configuración Docker optimizada:

### Básico
```yaml
app:
  build:
    context: ../../..      # Raíz del proyecto
    dockerfile: server/Dockerfile
```

### Avanzado
```yaml
app-advanced:
  build:
    context: ../../..      # Raíz del proyecto
    dockerfile: examples/server/advanced/Dockerfile
  command: ["/app/server", "-workers=20", "-queue-size=500", ...]
```

### Servicios Docker:
- **RabbitMQ**: Puerto 5672 (AMQP) y 15672 (Management UI)
- **MariaDB**: Puerto 3306 con inicialización automática
- **App**: Servidor burrowctl con health checks
- **App-Advanced**: Servidor optimizado con configuración empresarial

### Comandos Make:
```bash
# Básico
make docker-up          # Levantar entorno básico
make docker-down        # Detener entorno básico
make docker-logs        # Ver logs básico

# Avanzado
make docker-up-advanced    # Levantar entorno avanzado
make docker-down-advanced  # Detener entorno avanzado
make docker-logs-advanced  # Ver logs avanzado
```

## 📊 Comparación de Ejemplos

| Característica | Basic | Advanced |
|---------------|--------|----------|
| Worker Pool | ✅ (defecto) | ⚙️ Configurable |
| Rate Limiting | ✅ (defecto) | ⚙️ Configurable |
| Connection Pool | ✅ (defecto) | ⚙️ Configurable |
| Configuración | Hardcoded | 🎛️ Via flags |
| Monitoreo | Logs básicos | 📊 Métricas detalladas |
| Documentación | README básico | 📚 Guía completa |
| Docker | 🐳 Básico | 🐳 Optimizado |
| Dockerfile | Compartido | Específico |
| Comandos Make | docker-up | docker-up-advanced |

## 🎓 Recommended Learning Path

1. **Empezar** con `basic/` para entender conceptos fundamentales
2. **Experimentar** con `advanced/` para características empresariales
3. **Personalizar** configuración según necesidades específicas
4. **Implementar** en producción con parámetros optimizados

## 🔗 Referencias

- **Documentación Completa**: `../ADVANCED_FEATURES.md`
- **Ejemplos de Cliente**: `../client/`
- **Cliente Avanzado**: `../client/advanced/`
- **Proyecto Principal**: `../../`

## 💡 Tips

### Para Desarrollo:
```bash
cd basic/
docker-compose up -d  # Solo servicios
go run server_example.go -debug=true
```

### Para Testing:
```bash
cd advanced/
go run advanced_server_example.go \
  -workers=2 -rate-limit=50 -debug=true
```

### Para Producción:
```bash
cd advanced/
go run advanced_server_example.go \
  -workers=20 -queue-size=500 \
  -pool-open=50 -rate-limit=100
```