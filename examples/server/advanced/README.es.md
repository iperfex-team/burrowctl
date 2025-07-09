# Servidor Avanzado

Implementación mejorada del servidor burrowctl con características empresariales para entornos de alto rendimiento.

## Características

- **Pool de Workers**: Procesamiento concurrente configurable
- **Limitación de Velocidad**: Limitación por IP con soporte de ráfagas
- **Pool de Conexiones**: Gestión optimizada de conexiones DB
- **Monitoreo de Rendimiento**: Métricas en tiempo real
- **Apagado Elegante**: Cierre limpio con drenaje de requests

## Uso

### Ejecución directa
```bash
go run main.go
```

### Usando Makefile
```bash
make run-server-advanced
```

### Docker
```bash
make docker-up-advanced
```

## Configuración

### Opciones de Línea de Comandos

#### Opciones de Rendimiento
- `-workers=20`: Número de goroutines worker (defecto: 10)
- `-queue-size=500`: Tamaño de cola de worker (defecto: 100)
- `-rate-limit=50`: Límite de velocidad por IP (req/s) (defecto: 10)
- `-burst-size=100`: Tamaño de ráfaga (defecto: 20)

#### Opciones de Base de Datos
- `-pool-idle=20`: Máximo de conexiones idle (defecto: 5)
- `-pool-open=50`: Máximo de conexiones abiertas (defecto: 15)

### Configuraciones de Ejemplo

#### Modo Alto Rendimiento
```bash
go run main.go -workers=50 -queue-size=1000 -rate-limit=100 -burst-size=200
```

## Características de Rendimiento

### Pool de Workers
- Número configurable de workers concurrentes
- Cola limitada con protección de overflow
- Balanceo de carga entre workers

### Limitación de Velocidad
- Limitación por IP de cliente
- Algoritmo token bucket con soporte de ráfagas
- Limpieza automática de datos de límite

### Pool de Conexiones
- Gestión optimizada de conexiones DB
- Límites configurables idle/open
- Gestión de tiempo de vida de conexiones

## Siguientes Pasos

Para configuraciones especializadas:
- [Servidor de Cache](cache-server/README.es.md) - Cache de resultados de consultas
- [Servidor de Validación](validation-server/README.es.md) - Validación de seguridad SQL
- [Servidor Completo](full-featured-server/README.es.md) - Todas las características empresariales