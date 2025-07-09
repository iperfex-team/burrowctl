# Ejemplo de Servidor Básico

Una implementación simple del servidor burrowctl que demuestra la funcionalidad básica.

## Características

- Manejo básico de mensajes AMQP
- Conexión a base de datos MySQL
- Ejecución simple de comandos
- Patrón básico de solicitud/respuesta

## Uso

### Ejecución directa
```bash
go run main.go
```

### Usando Makefile
```bash
make run-server-example
```

### Docker
```bash
make docker-up
```

## Configuración

El servidor usa valores de configuración hardcodeados para simplicidad:

- **ID del Dispositivo**: `my-device`
- **URL AMQP**: `amqp://burrowuser:burrowpass123@localhost:5672/`
- **DSN MySQL**: `burrowuser:burrowpass123@tcp(localhost:3306)/burrowdb`

## Dependencias

- Servidor RabbitMQ (puerto 5672)
- Servidor MariaDB/MySQL (puerto 3306)
- Go 1.22 o superior

## Primeros Pasos

1. Iniciar los servicios requeridos:
   ```bash
   make docker-up
   ```

2. Ejecutar el servidor:
   ```bash
   make run-server-example
   ```

3. Probar con un cliente:
   ```bash
   cd ../../client/command-example
   go run main.go "ls -la"
   ```

## Arquitectura

Este servidor básico proporciona:

- **Integración de Cola de Mensajes**: Se conecta a RabbitMQ para recibir comandos
- **Conexión a Base de Datos**: Usa MySQL para almacenamiento de datos
- **Procesamiento de Comandos**: Ejecuta comandos recibidos y devuelve resultados
- **Manejo de Errores**: Manejo básico de errores y logging

## Próximos Pasos

Para características más avanzadas, consulta:
- [Servidor Avanzado](../advanced/README.es.md)
- [Servidor de Cache](../advanced/cache-server/README.es.md)
- [Servidor de Validación](../advanced/validation-server/README.es.md)
- [Servidor Completo](../advanced/full-featured-server/README.es.md)