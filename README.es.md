# 🐇 burrowctl

<div align="center">
  <h3>Ejecución Remota de SQL y Control de Dispositivos vía RabbitMQ</h3>
  <p>
    <strong>burrowctl</strong> es una potente librería y servicio de Go que proporciona un puente basado en RabbitMQ para ejecutar remotamente consultas SQL, funciones personalizadas y comandos del sistema en dispositivos detrás de NAT o firewalls.
  </p>
  <p>
    <a href="./README.md">🇺🇸 English</a> | 
    <a href="./README.es.md">🇪🇸 Español</a> | 
    <a href="./README.pt.md">🇧🇷 Português</a>
  </p>
</div>

## 🎯 ¿Qué es burrowctl?

**burrowctl** permite acceso seguro a bases de datos remotas y control de dispositivos sin exponer conexiones directas. Es perfecto para:

- 🏢 **Plataformas SaaS**: Gestionar bases de datos de clientes detrás de NAT/firewalls
- 🌐 **Gestión IoT**: Controlar dispositivos distribuidos de forma segura
- 🔐 **Administración Remota**: Ejecutar consultas y comandos sin SSH/acceso directo a BD
- 📊 **Monitoreo Distribuido**: Recopilar datos de múltiples fuentes remotas

## ✨ Características Principales

### 🔌 **Soporte Multi-Cliente**
- **Cliente Go**: Compatibilidad nativa con driver `database/sql`
- **Cliente Node.js/TypeScript**: API async moderna con tipado completo
- **DSN Universal**: Mismo formato de cadena de conexión para todos los clientes

### 🚀 **Tres Tipos de Ejecución**
- **Consultas SQL**: Acceso directo a base de datos con binding de parámetros
- **Funciones Personalizadas**: Sistema de funciones extensible con 16+ funciones incorporadas
- **Comandos del Sistema**: Ejecutar comandos del SO con acceso controlado

### 🔒 **Listo para Empresa**
- **Transporte Seguro**: Protocolo RabbitMQ AMQP 0-9-1
- **Pool de Conexiones**: Pools de conexiones de base de datos configurables
- **Manejo de Errores**: Gestión integral de errores y debugging
- **Control de Timeouts**: Timeouts configurables para consultas y comandos

### 📦 **Características de Producción**
- **Soporte Docker**: Entorno de desarrollo completamente containerizado
- **Automatización Makefile**: Automatización de build, test y despliegue
- **Control de Versiones**: Versionado semántico automático
- **Múltiples Ejemplos**: Ejemplos de uso y documentación comprensiva

---

## 🚀 Inicio Rápido

### Prerrequisitos

- **Go 1.22+** para cliente/servidor Go
- **Node.js 16+** para cliente TypeScript
- Servidor **RabbitMQ** ejecutándose
- Base de datos **MySQL/MariaDB**

### Instalación

```bash
git clone https://github.com/lordbasex/burrowctl.git
cd burrowctl
go mod tidy
```

### Uso Básico

#### Cliente Go (SQL)
```go
package main

import (
    "database/sql"
    "fmt"
    "log"
    _ "github.com/lordbasex/burrowctl/client"
)

func main() {
    dsn := "deviceID=mi-dispositivo&amqp_uri=amqp://user:pass@localhost:5672/&timeout=10s&debug=true"
    
    db, err := sql.Open("rabbitsql", dsn)
    if err != nil {
        log.Fatal("Falló la conexión:", err)
    }
    defer db.Close()
    
    rows, err := db.Query("SELECT id, name FROM users WHERE active = ?", true)
    if err != nil {
        log.Fatal("Falló la consulta:", err)
    }
    defer rows.Close()
    
    for rows.Next() {
        var id int
        var name string
        rows.Scan(&id, &name)
        fmt.Printf("ID: %d, Nombre: %s\n", id, name)
    }
}
```

#### Cliente Node.js/TypeScript
```typescript
import { createClient } from 'burrowctl-client-nodejs';

const client = await createClient(
  'deviceID=mi-dispositivo&amqp_uri=amqp://user:pass@localhost:5672/&timeout=10s'
);

const rows = await client.query('SELECT * FROM users WHERE active = ?', [true]);
console.log('Resultados:', rows.getRows());
console.log('Columnas:', rows.getColumns());

await client.close();
```

#### Configuración del Servidor
```go
package main

import (
    "context"
    "log"
    "time"
    "github.com/lordbasex/burrowctl/server"
)

func main() {
    pool := &server.PoolConfig{
        MaxIdleConns:    10,
        MaxOpenConns:    20,
        ConnMaxLifetime: 5 * time.Minute,
    }
    
    handler := server.NewHandler(
        "mi-dispositivo",                               // ID del dispositivo
        "amqp://user:pass@localhost:5672/",            // URI RabbitMQ
        "user:pass@tcp(localhost:3306)/dbname",        // DSN MySQL
        "open",                                        // Modo de conexión
        pool,                                          // Configuración del pool
    )
    
    // Registrar funciones personalizadas
    handler.RegisterFunction("obtenerInfoSistema", obtenerInfoSistema)
    handler.RegisterFunction("procesarDatos", procesarDatos)
    
    ctx := context.Background()
    log.Println("Iniciando servidor burrowctl...")
    if err := handler.Start(ctx); err != nil {
        log.Fatal("Falló el servidor:", err)
    }
}
```

---

## 📖 Tipos de Ejecución

### 1. 🗃️ Consultas SQL (`sql`)

Ejecuta consultas SQL directas con binding de parámetros y soporte completo de transacciones.

```go
// Cliente Go
rows, err := db.Query("SELECT * FROM products WHERE category = ? AND price > ?", "electronics", 100)

// Cliente Node.js
const rows = await client.query("SELECT * FROM products WHERE category = ? AND price > ?", ["electronics", 100]);
```

**Características:**
- Binding de parámetros para seguridad
- Soporte de transacciones
- Pool de conexiones
- Manejo de resultados con tipado seguro

### 2. ⚙️ Funciones Personalizadas (`function`)

Ejecuta funciones del lado del servidor con parámetros tipados y múltiples valores de retorno.

```go
// Cliente Go - usando solicitud de función JSON
funcReq := FunctionRequest{
    Name: "calcularImpuesto",
    Params: []FunctionParam{
        {Type: "float64", Value: 100.0},
        {Type: "string", Value: "ES"},
    },
}
jsonData, _ := json.Marshal(funcReq)
rows, err := db.Query("FUNCTION:" + string(jsonData))
```

```typescript
// Cliente Node.js
const result = await client.query('FUNCTION:{"name":"calcularImpuesto","params":[{"type":"float64","value":100.0},{"type":"string","value":"ES"}]}');
```

**Funciones Incorporadas (16+):**
- `lengthOfString`: Obtener longitud de cadena
- `addIntegers`: Sumar dos enteros
- `getCurrentTimestamp`: Obtener timestamp actual
- `generateUUID`: Generar UUID
- `encodeBase64`: Codificación Base64
- `decodeBase64`: Decodificación Base64
- `parseJSON`: Parsear cadena JSON
- `formatJSON`: Formatear JSON con indentación
- `getSystemInfo`: Obtener información del sistema
- `listFiles`: Listar contenido de directorio
- `readFile`: Leer contenido de archivo
- `writeFile`: Escribir contenido de archivo
- `calculateHash`: Calcular hash SHA256
- `validateEmail`: Validar dirección de email
- `generateRandomString`: Generar cadena aleatoria
- `convertTimezone`: Convertir zona horaria

### 3. 🖥️ Comandos del Sistema (`command`)

Ejecuta comandos del sistema con acceso controlado y gestión de timeouts.

```go
// Cliente Go
rows, err := db.Query("COMMAND:ps aux | grep mysql")
rows, err := db.Query("COMMAND:df -h")
rows, err := db.Query("COMMAND:systemctl status nginx")
```

```typescript
// Cliente Node.js
const result = await client.query('COMMAND:ps aux | grep mysql');
const diskUsage = await client.query('COMMAND:df -h');
```

**Características:**
- Captura de stdout/stderr
- Timeouts configurables
- Preservación de salida línea por línea
- Manejo de códigos de error

---

## 🔧 Configuración

### Formato DSN
```
deviceID=<id-dispositivo>&amqp_uri=<url-rabbitmq>&timeout=<timeout>&debug=<boolean>
```

**Parámetros:**
- `deviceID`: Identificador único del dispositivo (típicamente hash SHA256)
- `amqp_uri`: URL de conexión RabbitMQ
- `timeout`: Timeout de consulta (ej., `5s`, `30s`, `2m`)
- `debug`: Habilitar logging de debug (`true`/`false`)

### Configuración del Pool de Conexiones
```go
pool := &server.PoolConfig{
    MaxIdleConns:    10,          // Máximo de conexiones idle
    MaxOpenConns:    20,          // Máximo de conexiones abiertas
    ConnMaxLifetime: 5 * time.Minute, // Tiempo de vida de conexión
}
```

### Modos de Conexión
- **`open`**: Mantiene pool de conexiones (por defecto, mejor rendimiento)
- **`close`**: Abre/cierra conexiones por consulta (más seguro, más lento)

---

## 🛠️ Desarrollo

### Configuración Rápida de Desarrollo
```bash
# Clonar y configurar
git clone https://github.com/lordbasex/burrowctl.git
cd burrowctl

# Iniciar entorno de desarrollo (Docker)
cd examples/server
docker-compose up -d

# Construir proyecto
make build

# Ejecutar ejemplos
make run-server-example
make run-sql-example
make run-function-example
make run-command-example
```

### Comandos Make Disponibles
```bash
make help                    # Mostrar todos los comandos disponibles
make build                   # Construir todos los componentes
make test                    # Ejecutar pruebas
make clean                   # Limpiar artefactos de build
make docker-up              # Iniciar entorno Docker
make docker-down            # Detener entorno Docker
make run-server-example     # Ejecutar ejemplo del servidor
make run-sql-example        # Ejecutar ejemplo de cliente SQL
make run-function-example   # Ejecutar ejemplo de cliente de funciones
make run-command-example    # Ejecutar ejemplo de cliente de comandos
```

---

## 🏗️ Arquitectura

### Componentes del Sistema

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Cliente Go    │    │   Cliente       │    │   Futuros       │
│   (database/sql)│    │   Node.js       │    │   Clientes      │
└─────────┬───────┘    └─────────┬───────┘    └─────────┬───────┘
          │                      │                      │
          └──────────────────────┼──────────────────────┘
                                 │
                    ┌─────────────▼──────────────┐
                    │       RabbitMQ             │
                    │    (AMQP 0-9-1)           │
                    └─────────────┬──────────────┘
                                 │
                ┌─────────────────▼──────────────────┐
                │       Servidor burrowctl           │
                │  ┌─────────────────────────────┐   │
                │  │    Motor SQL               │   │
                │  │    Registro de Funciones   │   │
                │  │    Ejecutor de Comandos    │   │
                │  └─────────────────────────────┘   │
                └─────────────────┬──────────────────┘
                                 │
                    ┌─────────────▼──────────────┐
                    │       MySQL/MariaDB        │
                    │       Sistema de Archivos  │
                    │       Comandos del Sistema │
                    └────────────────────────────┘
```

### Flujo de Mensajes

1. **Cliente**: Envía solicitud a cola RabbitMQ específica del dispositivo
2. **RabbitMQ**: Enruta mensaje a cola apropiada del dispositivo
3. **Servidor**: Procesa solicitud basado en tipo (`sql`, `function`, `command`)
4. **Ejecución**: Ejecuta contra base de datos, registro de funciones, o sistema
5. **Respuesta**: Retorna resultados vía cola de respuesta RabbitMQ
6. **Cliente**: Recibe y procesa respuesta

---

## 📁 Estructura del Proyecto

```
burrowctl/
├── client/                 # Cliente Go (driver database/sql)
│   ├── driver.go          # Implementación del driver SQL
│   ├── conn.go            # Gestión de conexiones
│   ├── rows.go            # Manejo de resultados
│   └── rpc.go             # Cliente RPC RabbitMQ
├── server/                 # Librería del servidor core
│   └── server.go          # Implementación del servidor
├── client-nodejs/          # Cliente Node.js/TypeScript
│   ├── src/               # Código fuente TypeScript
│   ├── dist/              # JavaScript compilado
│   └── package.json       # Configuración del paquete NPM
├── examples/              # Ejemplos de uso
│   ├── client/            # Ejemplos de cliente
│   │   ├── sql-example/   # Uso de SQL
│   │   ├── function-example/ # Uso de funciones
│   │   └── command-example/  # Uso de comandos
│   └── server/            # Ejemplos de servidor
│       ├── server_example.go # Configuración completa del servidor
│       └── docker-compose.yml # Entorno de desarrollo
├── Makefile              # Automatización de build
├── go.mod               # Dependencias del módulo Go
└── version.txt          # Información de versión
```

---

## 🔐 Consideraciones de Seguridad

### Mejores Prácticas

1. **Usar Credenciales Fuertes**: Siempre usar contraseñas fuertes para RabbitMQ y base de datos
2. **Habilitar TLS**: Usar TLS/SSL para conexiones RabbitMQ en producción
3. **Limitar Acceso a Funciones**: Solo registrar funciones necesarias en el servidor
4. **Restricciones de Comandos**: Implementar lista blanca de comandos para seguridad
5. **Aislamiento de Red**: Usar VPNs o redes privadas cuando sea posible
6. **Monitoreo**: Implementar logging y monitoreo para auditoría de seguridad

### Configuración de Producción

```go
// Configuración del servidor para producción
handler := server.NewHandler(
    os.Getenv("DEVICE_ID"),
    os.Getenv("AMQP_URI"),     // Usar TLS: amqps://user:pass@host:5671/
    os.Getenv("MYSQL_DSN"),    // Usar SSL: ?tls=true
    "open",
    &server.PoolConfig{
        MaxIdleConns:    5,
        MaxOpenConns:    10,
        ConnMaxLifetime: 2 * time.Minute,
    },
)
```

---

## 🤝 Contribuir

¡Damos la bienvenida a contribuciones! Por favor vea nuestra [Guía de Contribución](CONTRIBUTING.md) para detalles.

### Configuración de Desarrollo

1. Fork el repositorio
2. Crear una rama de feature: `git checkout -b feature/caracteristica-increible`
3. Realizar cambios
4. Agregar tests para nueva funcionalidad
5. Ejecutar tests: `make test`
6. Commit de cambios: `git commit -m 'Agregar característica increíble'`
7. Push a la rama: `git push origin feature/caracteristica-increible`
8. Abrir un Pull Request

---

## 📜 Licencia

Este proyecto está licenciado bajo la Licencia MIT - vea el archivo [LICENSE](LICENSE) para detalles.

---

## 🆘 Soporte

- **Documentación**: [Documentación completa](./examples/)
- **Ejemplos**: [Ejemplos de uso](./examples/client/)
- **Issues**: [GitHub Issues](https://github.com/lordbasex/burrowctl/issues)
- **Discusiones**: [GitHub Discussions](https://github.com/lordbasex/burrowctl/discussions)

---

## 🙏 Agradecimientos

- [RabbitMQ](https://www.rabbitmq.com/) por el excelente broker de mensajes
- [Go SQL Driver](https://github.com/go-sql-driver/mysql) por la conectividad MySQL
- [AMQP 0-9-1 Go Client](https://github.com/rabbitmq/amqp091-go) por la integración RabbitMQ
- Las comunidades Go y Node.js por sus excelentes ecosistemas

---

<div align="center">
  <p>Hecho con ❤️ por el equipo burrowctl</p>
  <p>
    <a href="https://github.com/lordbasex/burrowctl/stargazers">⭐ Dar estrella a este proyecto</a> | 
    <a href="https://github.com/lordbasex/burrowctl/issues">🐛 Reportar Bug</a> | 
    <a href="https://github.com/lordbasex/burrowctl/issues">💡 Solicitar Feature</a>
  </p>
</div>