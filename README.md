# 🐇 burrowctl

<div align="center">
  <h3>Remote SQL Execution & Device Control via RabbitMQ</h3>
  <p>
    <strong>burrowctl</strong> is a powerful Go library and service that provides a RabbitMQ-based bridge to remotely execute SQL queries, custom functions, and system commands on devices behind NAT or firewalls.
  </p>
  <p>
    <a href="./README.md">🇺🇸 English</a> | 
    <a href="./README.es.md">🇪🇸 Español</a> | 
    <a href="./README.pt.md">🇧🇷 Português</a>
  </p>
</div>

## 🎯 What is burrowctl?

**burrowctl** enables secure remote database access and device control without exposing direct connections. It's perfect for:

- 🏢 **SaaS Platforms**: Manage customer databases behind NAT/firewalls
- 🌐 **IoT Management**: Control distributed devices securely
- 🔐 **Remote Administration**: Execute queries and commands without SSH/direct DB access
- 📊 **Distributed Monitoring**: Collect data from multiple remote sources

## ✨ Key Features

### 🔌 **Multi-Client Support**
- **Go Client**: Native `database/sql` driver compatibility
- **Node.js/TypeScript Client**: Modern async API with full type safety
- **Universal DSN**: Same connection string format across all clients

### 🚀 **Three Execution Types**
- **SQL Queries**: Direct database access with parameter binding
- **Custom Functions**: Extensible function system with 16+ built-in functions
- **System Commands**: Execute OS commands with controlled access

### 🔒 **Enterprise-Ready**
- **Secure Transport**: RabbitMQ AMQP 0-9-1 protocol
- **Connection Pooling**: Configurable database connection pools
- **Error Handling**: Comprehensive error management and debugging
- **Timeout Control**: Configurable query and command timeouts

### 🏗️ **Enterprise Features** (NEW)
- **🔄 Worker Pool**: Concurrent message processing (10-50+ workers)
- **🛡️ Rate Limiting**: Per-client IP protection with token bucket algorithm
- **📝 Prepared Statements**: Client-side statement caching and SQL injection protection
- **🔄 Automatic Reconnection**: Connection recovery with exponential backoff
- **📊 Performance Monitoring**: Real-time metrics and configurable parameters
- **⚙️ Advanced Configuration**: Granular control over all performance aspects

### 📦 **Production Features**
- **Docker Support**: Complete containerized development environment
- **Makefile Automation**: Build, test, and deployment automation
- **Version Control**: Automatic semantic versioning
- **Multiple Examples**: Comprehensive usage examples and documentation

---

## 🚀 Quick Start

### Prerequisites

- **Go 1.22+** for Go client/server
- **Node.js 16+** for TypeScript client
- **RabbitMQ** server running
- **MySQL/MariaDB** database

### Installation

```bash
git clone https://github.com/lordbasex/burrowctl.git
cd burrowctl
go mod tidy
```

### Basic Usage

#### Go Client (SQL)
```go
package main

import (
    "database/sql"
    "fmt"
    "log"
    _ "github.com/lordbasex/burrowctl/client"
)

func main() {
    dsn := "deviceID=my-device&amqp_uri=amqp://user:pass@localhost:5672/&timeout=10s&debug=true"
    
    db, err := sql.Open("rabbitsql", dsn)
    if err != nil {
        log.Fatal("Connection failed:", err)
    }
    defer db.Close()
    
    rows, err := db.Query("SELECT id, name FROM users WHERE active = ?", true)
    if err != nil {
        log.Fatal("Query failed:", err)
    }
    defer rows.Close()
    
    for rows.Next() {
        var id int
        var name string
        rows.Scan(&id, &name)
        fmt.Printf("ID: %d, Name: %s\n", id, name)
    }
}
```

#### Node.js/TypeScript Client
```typescript
import { createClient } from 'burrowctl-client-nodejs';

const client = await createClient(
  'deviceID=my-device&amqp_uri=amqp://user:pass@localhost:5672/&timeout=10s'
);

const rows = await client.query('SELECT * FROM users WHERE active = ?', [true]);
console.log('Results:', rows.getRows());
console.log('Columns:', rows.getColumns());

await client.close();
```

#### Server Setup
```go
package main

import (
    "context"
    "log"
    "github.com/lordbasex/burrowctl/server"
)

func main() {
    pool := &server.PoolConfig{
        MaxIdleConns:    10,
        MaxOpenConns:    20,
        ConnMaxLifetime: 5 * time.Minute,
    }
    
    handler := server.NewHandler(
        "my-device",                                    // Device ID
        "amqp://user:pass@localhost:5672/",            // RabbitMQ URI
        "user:pass@tcp(localhost:3306)/dbname",        // MySQL DSN
        "open",                                        // Connection mode
        pool,                                          // Pool config
    )
    
    // Register custom functions
    handler.RegisterFunction("getSystemInfo", getSystemInfo)
    handler.RegisterFunction("processData", processData)
    
    ctx := context.Background()
    log.Println("Starting burrowctl server...")
    if err := handler.Start(ctx); err != nil {
        log.Fatal("Server failed:", err)
    }
}
```

---

## 📖 Execution Types

### 1. 🗃️ SQL Queries (`sql`)

Execute direct SQL queries with parameter binding and full transaction support.

```go
// Go client
rows, err := db.Query("SELECT * FROM products WHERE category = ? AND price > ?", "electronics", 100)

// Node.js client
const rows = await client.query("SELECT * FROM products WHERE category = ? AND price > ?", ["electronics", 100]);
```

**Features:**
- Parameter binding for security
- Transaction support
- Connection pooling
- Type-safe result handling

### 2. ⚙️ Custom Functions (`function`)

Execute server-side functions with typed parameters and multiple return values.

```go
// Go client - using JSON function request
funcReq := FunctionRequest{
    Name: "calculateTax",
    Params: []FunctionParam{
        {Type: "float64", Value: 100.0},
        {Type: "string", Value: "US"},
    },
}
jsonData, _ := json.Marshal(funcReq)
rows, err := db.Query("FUNCTION:" + string(jsonData))
```

```typescript
// Node.js client
const result = await client.query('FUNCTION:{"name":"calculateTax","params":[{"type":"float64","value":100.0},{"type":"string","value":"US"}]}');
```

**Built-in Functions (16+):**
- `lengthOfString`: Get string length
- `addIntegers`: Add two integers
- `getCurrentTimestamp`: Get current timestamp
- `generateUUID`: Generate UUID
- `encodeBase64`: Base64 encoding
- `decodeBase64`: Base64 decoding
- `parseJSON`: Parse JSON string
- `formatJSON`: Format JSON with indentation
- `getSystemInfo`: Get system information
- `listFiles`: List directory contents
- `readFile`: Read file contents
- `writeFile`: Write file contents
- `calculateHash`: Calculate SHA256 hash
- `validateEmail`: Validate email address
- `generateRandomString`: Generate random string
- `convertTimezone`: Convert timezone

### 3. 🖥️ System Commands (`command`)

Execute system commands with controlled access and timeout management.

```go
// Go client
rows, err := db.Query("COMMAND:ps aux | grep mysql")
rows, err := db.Query("COMMAND:df -h")
rows, err := db.Query("COMMAND:systemctl status nginx")
```

```typescript
// Node.js client
const result = await client.query('COMMAND:ps aux | grep mysql');
const diskUsage = await client.query('COMMAND:df -h');
```

**Features:**
- Stdout/stderr capture
- Configurable timeouts
- Line-by-line output preservation
- Error code handling

---

## 🔧 Configuration

### DSN Format
```
deviceID=<device-id>&amqp_uri=<rabbitmq-url>&timeout=<timeout>&debug=<boolean>
```

**Parameters:**
- `deviceID`: Unique device identifier (typically SHA256 hash)
- `amqp_uri`: RabbitMQ connection URL
- `timeout`: Query timeout (e.g., `5s`, `30s`, `2m`)
- `debug`: Enable debug logging (`true`/`false`)

### Connection Pool Configuration
```go
pool := &server.PoolConfig{
    MaxIdleConns:    10,          // Maximum idle connections
    MaxOpenConns:    20,          // Maximum open connections
    ConnMaxLifetime: 5 * time.Minute, // Connection lifetime
}
```

### Connection Modes
- **`open`**: Maintains connection pool (default, better performance)
- **`close`**: Opens/closes connections per query (safer, slower)

---

## 🛠️ Development

### Quick Development Setup
```bash
# Clone and setup
git clone https://github.com/lordbasex/burrowctl.git
cd burrowctl

# Start development environment (Docker)
cd examples/server
docker-compose up -d

# Build project
make build

# Run examples
make run-server-example
make run-sql-example
make run-function-example
make run-command-example
```

### Available Make Commands
```bash
make help                    # Show all available commands
make build                   # Build all components
make test                    # Run tests
make clean                   # Clean build artifacts

# Docker environments
make docker-up              # Basic server environment
make docker-up-advanced     # Advanced server environment
make docker-up-cache        # Cache-optimized server
make docker-up-validation   # SQL validation server
make docker-up-full         # Full enterprise server

# Server examples
make run-server-example     # Basic server
make run-server-advanced    # Advanced server
make run-server-cache       # Cache-optimized server
make run-server-validation  # SQL validation server
make run-server-full        # Full enterprise server

# Client examples
make run-sql-example        # SQL client example
make run-function-example   # Function client example
make run-command-example    # Command client example
```

### Docker Environment

The project includes a complete Docker Compose environment:

```yaml
services:
  rabbitmq:
    image: rabbitmq:3-management
    ports:
      - "5672:5672"
      - "15672:15672"
    
  mariadb:
    image: mariadb:10.6
    environment:
      MYSQL_ROOT_PASSWORD: rootpass
      MYSQL_DATABASE: burrowdb
      MYSQL_USER: burrowuser
      MYSQL_PASSWORD: burrowpass123
    
  burrowctl-server:
    build: .
    depends_on:
      - rabbitmq
      - mariadb
```

---

## 🏗️ Architecture

### System Components

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Go Client     │    │   Node.js       │    │   Future        │
│   (database/sql)│    │   Client        │    │   Clients       │
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
                │        burrowctl Server            │
                │  ┌─────────────────────────────┐   │
                │  │    SQL Engine              │   │
                │  │    Function Registry       │   │
                │  │    Command Executor        │   │
                │  └─────────────────────────────┘   │
                └─────────────────┬──────────────────┘
                                 │
                    ┌─────────────▼──────────────┐
                    │       MySQL/MariaDB        │
                    │       File System          │
                    │       System Commands      │
                    └────────────────────────────┘
```

### Message Flow

1. **Client**: Sends request to device-specific RabbitMQ queue
2. **RabbitMQ**: Routes message to appropriate device queue
3. **Server**: Processes request based on type (`sql`, `function`, `command`)
4. **Execution**: Executes against database, function registry, or system
5. **Response**: Returns results via RabbitMQ reply queue
6. **Client**: Receives and processes response

---

## 📁 Project Structure

```
burrowctl/
├── client/                 # Go client (database/sql driver)
│   ├── driver.go          # SQL driver implementation
│   ├── conn.go            # Connection management
│   ├── rows.go            # Result handling
│   └── rpc.go             # RabbitMQ RPC client
├── server/                 # Core server library
│   └── server.go          # Server implementation
├── client-nodejs/          # Node.js/TypeScript client
│   ├── src/               # TypeScript source
│   ├── dist/              # Compiled JavaScript
│   └── package.json       # NPM package configuration
├── examples/              # Usage examples
│   ├── client/            # Client examples
│   │   ├── sql-example/   # SQL usage
│   │   ├── function-example/ # Function usage
│   │   └── command-example/  # Command usage
│   └── server/            # Server examples
│       ├── Dockerfile     # Universal Docker build
│       ├── basic/         # Basic server (main.go)
│       └── advanced/      # Enterprise servers
│           ├── main.go    # Advanced server
│           ├── cache-server/    # Cache-optimized
│           ├── validation-server/ # SQL security
│           └── full-featured-server/ # Complete enterprise
├── Makefile              # Build automation
├── go.mod               # Go module dependencies
└── version.txt          # Version information
```

---

## 🚀 **Server Examples**

burrowctl now provides multiple server configurations for different use cases:

### 📋 **Basic Server** (`examples/server/basic/`)
Simple server implementation for getting started:
```bash
make run-server-example
# or
cd examples/server/basic && go run main.go
```

### 🚀 **Advanced Server** (`examples/server/advanced/`)
Enterprise server with performance features:
```bash
make run-server-advanced
# or
cd examples/server/advanced && go run main.go
```

### 📈 **Cache Server** (`examples/server/advanced/cache-server/`)
Optimized for high-volume query caching:
```bash
make run-server-cache
# or
cd examples/server/advanced/cache-server && go run main.go
```

### 🛡️ **Validation Server** (`examples/server/advanced/validation-server/`)
SQL security and validation focused:
```bash
make run-server-validation
# or
cd examples/server/advanced/validation-server && go run main.go
```

### 🏢 **Full-Featured Server** (`examples/server/advanced/full-featured-server/`)
Complete enterprise server with all features:
```bash
make run-server-full
# or
cd examples/server/advanced/full-featured-server && go run main.go
```

## 🚀 **Enterprise Configuration**

### High-Performance Server
```bash
# Advanced server with optimized settings
cd examples/server/advanced
go run main.go \
  -workers=20 \
  -queue-size=500 \
  -rate-limit=50 \
  -pool-open=50

# Or with Docker (auto-configured)
docker-compose up -d
```

### Advanced Client Features
```bash
# Advanced client with all features
cd examples/client/advanced
go run advanced-main.go -prepared -timeout=30s

# Stress test (rate limiting demo)
go run advanced-main.go -stress -concurrent=10 -requests=100

# Reconnection demo
go run advanced-main.go -reconnect-demo
```

### Performance Comparison
| Feature | Basic | Advanced |
|---------|-------|----------|
| **Throughput** | ~100 req/s | ~1000+ req/s |
| **Concurrency** | Sequential | 10-50+ parallel |
| **Rate Protection** | None | Per-client limiting |
| **Reconnection** | Manual | Automatic |
| **Statements** | One-time | Cached/prepared |

## 🔍 Advanced Usage

### Custom Function Registration

```go
// Define custom function
func calculateDiscount(price float64, percentage float64) (float64, error) {
    if percentage > 100 || percentage < 0 {
        return 0, errors.New("invalid percentage")
    }
    return price * (percentage / 100), nil
}

// Register function
handler.RegisterFunction("calculateDiscount", calculateDiscount)
```

### Transaction Support

```go
// Begin transaction
tx, err := db.Begin()
if err != nil {
    log.Fatal(err)
}

// Execute multiple queries
_, err = tx.Exec("INSERT INTO orders (customer_id, total) VALUES (?, ?)", 1, 100.50)
if err != nil {
    tx.Rollback()
    log.Fatal(err)
}

_, err = tx.Exec("UPDATE inventory SET quantity = quantity - 1 WHERE product_id = ?", 123)
if err != nil {
    tx.Rollback()
    log.Fatal(err)
}

// Commit transaction
err = tx.Commit()
if err != nil {
    log.Fatal(err)
}
```

### Error Handling

```go
// Go client
rows, err := db.Query("SELECT * FROM users")
if err != nil {
    if strings.Contains(err.Error(), "timeout") {
        log.Println("Query timed out")
    } else if strings.Contains(err.Error(), "connection refused") {
        log.Println("Cannot connect to RabbitMQ")
    } else {
        log.Printf("Query error: %v", err)
    }
}
```

```typescript
// Node.js client
try {
    const rows = await client.query('SELECT * FROM users');
    console.log(rows.getRows());
} catch (error) {
    if (error.message.includes('timeout')) {
        console.log('Query timed out');
    } else if (error.message.includes('connection refused')) {
        console.log('Cannot connect to RabbitMQ');
    } else {
        console.error('Query error:', error.message);
    }
}
```

---

## 🔐 Security Considerations

### Best Practices

1. **Use Strong Credentials**: Always use strong passwords for RabbitMQ and database
2. **Enable TLS**: Use TLS/SSL for RabbitMQ connections in production
3. **Limit Function Access**: Only register necessary functions on the server
4. **Command Restrictions**: Implement command whitelisting for security
5. **Network Isolation**: Use VPNs or private networks when possible
6. **Monitoring**: Implement logging and monitoring for security audit

### Production Configuration

```go
// Production server setup
handler := server.NewHandler(
    os.Getenv("DEVICE_ID"),
    os.Getenv("AMQP_URI"),     // Use TLS: amqps://user:pass@host:5671/
    os.Getenv("MYSQL_DSN"),    // Use SSL: ?tls=true
    "open",
    &server.PoolConfig{
        MaxIdleConns:    5,
        MaxOpenConns:    10,
        ConnMaxLifetime: 2 * time.Minute,
    },
)
```

---

## 🚀 Performance Tuning

### Connection Pool Optimization

```go
// High-throughput configuration
pool := &server.PoolConfig{
    MaxIdleConns:    20,
    MaxOpenConns:    50,
    ConnMaxLifetime: 1 * time.Hour,
}
```

### Client-Side Optimization

```go
// Prepare statements for repeated queries
stmt, err := db.Prepare("SELECT * FROM users WHERE department = ?")
if err != nil {
    log.Fatal(err)
}
defer stmt.Close()

// Execute prepared statement multiple times
for _, dept := range departments {
    rows, err := stmt.Query(dept)
    if err != nil {
        log.Printf("Query failed for %s: %v", dept, err)
        continue
    }
    // Process results...
    rows.Close()
}
```

---

## 📊 Monitoring & Debugging

### Enable Debug Logging

```go
// DSN with debug enabled
dsn := "deviceID=my-device&amqp_uri=amqp://localhost:5672/&debug=true"
```

### Performance Metrics

```go
// Add metrics to custom functions
func monitoredFunction(data string) (string, error) {
    start := time.Now()
    defer func() {
        duration := time.Since(start)
        log.Printf("Function executed in %v", duration)
    }()
    
    // Function logic here
    return processData(data)
}
```

---

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

### Development Setup

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/amazing-feature`
3. Make your changes
4. Add tests for new functionality
5. Run tests: `make test`
6. Commit your changes: `git commit -m 'Add amazing feature'`
7. Push to the branch: `git push origin feature/amazing-feature`
8. Open a Pull Request

---

## 📜 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

## 🆘 Support

- **Documentation**: [Full documentation](./examples/)
- **Examples**: [Usage examples](./examples/client/)
- **Issues**: [GitHub Issues](https://github.com/lordbasex/burrowctl/issues)
- **Discussions**: [GitHub Discussions](https://github.com/lordbasex/burrowctl/discussions)

---

## 🙏 Acknowledgments

- [RabbitMQ](https://www.rabbitmq.com/) for the excellent message broker
- [Go SQL Driver](https://github.com/go-sql-driver/mysql) for MySQL connectivity
- [AMQP 0-9-1 Go Client](https://github.com/rabbitmq/amqp091-go) for RabbitMQ integration
- The Go and Node.js communities for their excellent ecosystems

---

## 📄 License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

Copyright (c) 2024 Federico Pereira <lord.basex@gmail.com>

---

<div align="center">
  <p>Made with ❤️ by the burrowctl team</p>
  <p>
    <a href="https://github.com/lordbasex/burrowctl/stargazers">⭐ Star this project</a> | 
    <a href="https://github.com/lordbasex/burrowctl/issues">🐛 Report Bug</a> | 
    <a href="https://github.com/lordbasex/burrowctl/issues">💡 Request Feature</a>
  </p>
</div>