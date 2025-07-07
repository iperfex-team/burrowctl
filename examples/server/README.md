# Burrowctl Server Examples

Este directorio contiene ejemplos para ejecutar el servidor burrowctl con todas sus funcionalidades.

## Archivos Incluidos

### `server_example.go`
Servidor principal que implementa todas las funcionalidades:
- **📊 SQL Queries** - Ejecuta consultas SQL remotas
- **🔧 Functions** - Ejecuta funciones remotas con parámetros tipados
- **⚡ Commands** - Ejecuta comandos del sistema

#### Funciones Disponibles (16 total):
- **Sin parámetros:** `returnBool`, `returnInt`, `returnString`, `returnStruct`, `returnIntArray`, `returnStringArray`, `returnJSON`, `returnError`
- **Con parámetros:** `lengthOfString`, `isEven`, `sumArray`, `greetPerson`, `validateString`, `flagToPerson`, `modifyJSON`
- **Múltiples valores:** `complexFunction`

### `demo-func.go`
Archivo de documentación que muestra ejemplos de todas las funciones disponibles y cómo construir las solicitudes JSON.

### Archivos Docker
- `docker-compose.yml` - Configuración completa con RabbitMQ y MariaDB
- `docker-compose-basic.yml` - Configuración básica
- `init.sql` - Script de inicialización de la base de datos

## Configuración

### Variables de Entorno (Opcionales)
```bash
export BURROWCTL_DEVICE_ID="tu-device-id"
export BURROWCTL_AMQP_URL="amqp://user:pass@host:port/"
export BURROWCTL_MYSQL_DSN="user:pass@tcp(host:port)/db?parseTime=true"
export BURROWCTL_CONNECTION_MODE="open"  # o "close"
```

Si no se configuran, se usan valores por defecto para Docker.

## Ejecución

### 1. Con Docker (Recomendado)
```bash
# Levantar servicios (RabbitMQ + MariaDB)
docker-compose up -d

# Ejecutar servidor
go run server_example.go
```

### 2. Ejecutar Servidor
```bash
# Compilar
go build -o server_example server_example.go

# Ejecutar
./server_example
```

### 3. Ver Documentación de Funciones
```bash
# Mostrar todas las funciones disponibles y ejemplos de uso
go run demo-func.go
```

## Salida del Servidor

Al iniciar, el servidor muestra:

```
🚀 Starting burrowctl server...
📱 Device ID: fd1825ec5a7b63f3fa2be9e04154d3b16f676663ba38e23d4ffafa7b0df29efb
🐰 RabbitMQ: amqp://burrowuser:burrowpass123@rabbitmq:5672/
🗄️  MariaDB: burrowuser:burrowpass123@tcp(mariadb:3306)/burrowdb?parseTime=true
🔗 Connection mode: open

✅ Server capabilities:
   📊 SQL Queries - Execute remote SQL queries
   🔧 Functions - Execute remote functions with typed parameters
   ⚡ Commands - Execute system commands

🔧 Available Functions:
   • returnBool, returnInt, returnString, returnStruct
   • returnIntArray, returnStringArray, returnJSON, returnError
   • lengthOfString, isEven, sumArray, greetPerson
   • validateString, complexFunction, flagToPerson, modifyJSON

🎯 Example usage:
   SQL:      SELECT * FROM users
   Command:  COMMAND:ps aux
   Function: FUNCTION:{"name":"returnString","params":[]}
```

## Tipos de Solicitudes Soportadas

### 1. SQL Queries
```sql
SELECT * FROM users WHERE id = 1
```

### 2. System Commands
```bash
COMMAND:ps aux
COMMAND:ls -la
COMMAND:whoami
```

### 3. Function Execution
```json
FUNCTION:{"name":"returnString","params":[]}
FUNCTION:{"name":"lengthOfString","params":[{"type":"string","value":"Hello"}]}
FUNCTION:{"name":"sumArray","params":[{"type":"[]int","value":[1,2,3,4,5]}]}
```

## Pool de Conexiones

El servidor está configurado con:
- **MaxIdleConns:** 5 conexiones idle
- **MaxOpenConns:** 15 conexiones máximas
- **ConnMaxLifetime:** 5 minutos

## Clientes de Ejemplo

Para probar el servidor, usa los clientes en:
- `examples/client/sql-example/` - Para consultas SQL
- `examples/client/command-example/` - Para comandos del sistema
- `examples/client/function-example/` - Para ejecución de funciones

## Logs

El servidor registra todas las operaciones:
```
[server] received ip=192.168.1.100 type=function query={"name":"returnString","params":[]}
[server] executing function: {"name":"returnString","params":[]}
[server] function executed successfully
```

## Troubleshooting

### Problemas Comunes

1. **Error de conexión a RabbitMQ**
   - Verificar que RabbitMQ esté ejecutándose
   - Verificar credenciales y puerto

2. **Error de conexión a MariaDB**
   - Verificar que MariaDB esté ejecutándose
   - Verificar credenciales, host y puerto
   - Verificar que la base de datos existe

3. **Función no encontrada**
   - Verificar que el nombre de la función sea exacto
   - Revisar la lista de funciones disponibles en `demo-func.go`

### Verificar Servicios Docker
```bash
# Ver logs de RabbitMQ
docker-compose logs rabbitmq

# Ver logs de MariaDB
docker-compose logs mariadb

# Ver estado de servicios
docker-compose ps
```

## Extensión

Para agregar nuevas funciones:

1. Implementar la función en `server/server.go`
2. Registrarla en el mapa `functions` de `getFunctionByName()`
3. Agregar soporte de tipos si es necesario en `convertToType()`
4. Actualizar la documentación en `demo-func.go`

¡El servidor está listo para manejar todas las capacidades de burrowctl! 