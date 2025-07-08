# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**burrowctl** is a Go library and service that provides a RabbitMQ-based bridge for remote SQL execution, custom functions, and system command execution on devices behind NAT or firewalls. It enables secure remote database access and device control without requiring direct connections.

## Development Commands

Use the comprehensive Makefile for all development tasks:

### Core Commands
- `make help` - Show all available commands
- `make build` - Build project and all examples
- `make test` - Run all tests
- `make clean` - Clean build artifacts
- `make install` - Install dependencies with `go mod tidy`

### Development Environment
- `make docker-up` - Start development environment (RabbitMQ + MariaDB)
- `make docker-down` - Stop development environment
- `make fmt` - Format code
- `make lint` - Run linter
- `make vet` - Run go vet

### Example Commands
- `make run-server-example` - Run server example
- `make run-sql-example` - Run SQL client example
- `make run-function-example` - Run function client example
- `make run-command-example` - Run command client example

### Release Commands
- `make tag VERSION=vX.Y.Z` - Create version tag
- `make release` - Create tag and push
- `make quick-release` - Auto-version release

## Architecture Overview

The project follows **clean architecture** with clear separation of concerns:

### Core Components
1. **Server Core** (`/server/server.go`): Clean infrastructure handling AMQP, MySQL, and HTTP
2. **Go Client** (`/client/`): Native `database/sql` driver implementation  
3. **Node.js Client** (`/client-nodejs/`): TypeScript client with async API
4. **Examples** (`/examples/`): Function registration and usage patterns

### Execution Flow
```
Client (Go/Node.js) → RabbitMQ (AMQP 0-9-1) → Server → Database/System/Functions
```

### Three Execution Types
- **SQL Queries**: Direct database access with parameter binding
- **Custom Functions**: Extensible function system with 16+ built-in functions
- **System Commands**: Execute OS commands with controlled access

## Key Directories

- `/server/` - Core server library (clean, no hardcoded functions)
- `/client/` - Go client implementing `database/sql` driver interface
- `/client-nodejs/` - TypeScript client with type definitions
- `/examples/server/` - Server setup with function registration patterns
- `/examples/client/` - Client usage examples for SQL, functions, and commands

## Development Patterns

### Clean Architecture
- Core server contains NO hardcoded functions
- Functions are registered dynamically via `RegisterFunction()` and `RegisterFunctions()`
- Examples demonstrate proper function registration patterns

### Connection Management
Two modes available:
- **"open"**: Connection pooling (default, better performance)
- **"close"**: Per-query connections (safer, slower)

### DSN Format
Universal DSN format across all clients:
```
deviceID=<device-id>&amqp_uri=<rabbitmq-url>&timeout=<timeout>&debug=<boolean>
```

## Technology Stack

- **Go 1.22+** with MySQL driver and RabbitMQ client
- **Node.js 16+** with TypeScript support
- **RabbitMQ** for message queuing
- **MySQL/MariaDB** for database operations
- **Docker** for development environment

## Development Environment

The project includes complete Docker Compose setup:
- RabbitMQ with management UI at `localhost:15672`
- MariaDB with automatic initialization
- Health checks and proper startup sequencing

## Version Management

The Makefile includes automatic versioning based on git commits, updating `version.txt` automatically. Current version is tracked in `version.txt`.