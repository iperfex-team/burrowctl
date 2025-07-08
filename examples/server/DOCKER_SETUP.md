# 🐳 Docker Setup Guide

Esta guía explica cómo usar Docker con los ejemplos de servidor de burrowctl.

## 📁 Estructura Docker

```
examples/server/
├── basic/
│   ├── docker-compose.yml       # Entorno completo básico
│   ├── docker-compose-basic.yml # Solo servicios
│   └── init.sql                 # Inicialización DB
└── advanced/
    ├── docker-compose.yml       # Entorno completo avanzado
    ├── docker-compose-basic.yml # Solo servicios
    ├── Dockerfile               # Dockerfile específico
    └── init.sql                 # Inicialización DB
```

## 🚀 Servidor Básico

### Entorno Completo
```bash
cd examples/server/basic
docker-compose up -d
```

**Servicios:**
- **RabbitMQ**: `localhost:5672` (AMQP) y `localhost:15672` (Management UI)
- **MariaDB**: `localhost:3306`
- **App**: Servidor básico con configuración estándar

### Solo Servicios
```bash
cd examples/server/basic
docker-compose -f docker-compose-basic.yml up -d
go run server_example.go
```

## 🏗️ Servidor Avanzado

### Entorno Completo
```bash
cd examples/server/advanced
docker-compose up -d
```

**Servicios:**
- **RabbitMQ**: `localhost:5672` (AMQP) y `localhost:15672` (Management UI)
- **MariaDB**: `localhost:3306`
- **App-Advanced**: Servidor optimizado con configuración empresarial

**Configuración automática:**
- Workers: 20
- Queue Size: 500
- Rate Limit: 50 req/sec
- DB Pool: 20-50 conexiones

### Solo Servicios
```bash
cd examples/server/advanced
docker-compose -f docker-compose-basic.yml up -d
go run advanced_server_example.go -workers=20 -rate-limit=50
```

## 🎛️ Comandos Make

### Básico
```bash
make docker-up          # Levantar entorno básico
make docker-down        # Detener entorno básico
make docker-logs        # Ver logs básico
```

### Avanzado
```bash
make docker-up-advanced    # Levantar entorno avanzado
make docker-down-advanced  # Detener entorno avanzado
make docker-logs-advanced  # Ver logs avanzado
```

## 🔧 Configuración Docker

### Básico (server/Dockerfile)
```dockerfile
# Compila: server_example.go
RUN go build -o burrowctl-server ./examples/server/basic/server_example.go
ENTRYPOINT ["/burrowctl-server"]
```

### Avanzado (examples/server/advanced/Dockerfile)
```dockerfile
# Compila: advanced_server_example.go
RUN go build -o burrowctl-server-advanced ./examples/server/advanced/advanced_server_example.go
ENTRYPOINT ["/app/server"]
```

## 📊 Diferencias Clave

| Aspecto | Básico | Avanzado |
|---------|--------|----------|
| **Dockerfile** | Compartido (`server/Dockerfile`) | Específico (`examples/server/advanced/Dockerfile`) |
| **Configuración** | Hardcoded en código | Via flags en docker-compose |
| **Contenedor** | `app` | `app-advanced` |
| **Comando** | Default ENTRYPOINT | Custom command con flags |
| **Optimización** | Configuración estándar | Configuración empresarial |

## 🌐 Puertos

- **RabbitMQ AMQP**: 5672
- **RabbitMQ Management**: 15672 (usuario: `burrowuser`, password: `burrowpass123`)
- **MariaDB**: 3306 (usuario: `burrowuser`, password: `burrowpass123`, db: `burrowdb`)

## 📝 Logs y Debugging

### Ver logs en tiempo real
```bash
# Básico
docker-compose logs -f app

# Avanzado
docker-compose logs -f app-advanced
```

### Ver estado de servicios
```bash
docker-compose ps
```

### Acceder a contenedor
```bash
# Básico
docker exec -it app bash

# Avanzado
docker exec -it app-advanced bash
```

## 🔍 Troubleshooting

### Problemas comunes:

**1. Puerto ya en uso**
```bash
# Detener servicios existentes
docker-compose down
# O cambiar puertos en docker-compose.yml
```

**2. Volúmenes corruptos**
```bash
docker-compose down -v  # Elimina volúmenes
docker-compose up -d    # Reinicia con volúmenes limpios
```

**3. Imagen no actualizada**
```bash
docker-compose build --no-cache
docker-compose up -d
```

**4. Conflictos entre básico y avanzado**
```bash
# Usar redes separadas o diferentes puertos
# Los ejemplos ya usan contenedores con nombres únicos
```

## 🚀 Desarrollo con Docker

### Desarrollo local con servicios Docker
```bash
# Levantar solo servicios
docker-compose -f docker-compose-basic.yml up -d

# Ejecutar servidor local
go run server_example.go  # o advanced_server_example.go

# Conectar cliente
cd ../client/sql-example
go run main.go
```

### Desarrollo completo en Docker
```bash
# Levantar entorno completo
docker-compose up -d

# Conectar cliente desde local
cd ../client/sql-example
go run main.go
```

## 📚 Referencias

- **Documentación básica**: `basic/README.md`
- **Documentación avanzada**: `advanced/README.md`
- **Guía general**: `README.md`
- **Makefile**: `../../Makefile`