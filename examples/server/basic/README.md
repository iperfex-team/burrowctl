# Basic Server Example

A simple burrowctl server implementation that demonstrates the core functionality.

## Features

- Basic AMQP message handling
- MySQL database connection
- Simple command execution
- Basic request/response pattern

## Usage

### Direct execution
```bash
go run main.go
```

### Using Makefile
```bash
make run-server-example
```

### Docker
```bash
make docker-up
```

## Configuration

The server uses hardcoded configuration values for simplicity:

- **Device ID**: `my-device`
- **AMQP URL**: `amqp://burrowuser:burrowpass123@localhost:5672/`
- **MySQL DSN**: `burrowuser:burrowpass123@tcp(localhost:3306)/burrowdb`

## Dependencies

- RabbitMQ server (port 5672)
- MariaDB/MySQL server (port 3306)
- Go 1.22 or higher

## Getting Started

1. Start the required services:
   ```bash
   make docker-up
   ```

2. Run the server:
   ```bash
   make run-server-example
   ```

3. Test with a client:
   ```bash
   cd ../../client/command-example
   go run main.go "ls -la"
   ```

## Architecture

This basic server provides:

- **Message Queue Integration**: Connects to RabbitMQ for receiving commands
- **Database Connection**: Uses MySQL for data storage
- **Command Processing**: Executes received commands and returns results
- **Error Handling**: Basic error handling and logging

## Next Steps

For more advanced features, check out:
- [Advanced Server](../advanced/README.md)
- [Cache Server](../advanced/cache-server/README.md)
- [Validation Server](../advanced/validation-server/README.md)
- [Full-Featured Server](../advanced/full-featured-server/README.md)