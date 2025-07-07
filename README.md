
# 🐇 burrowctl

**burrowctl** is a Go library and service that provides a RabbitMQ-based bridge to remotely execute SQL queries, custom functions, and system commands on a remote device.
The client is compatible with Go’s database/sql interface, making it easy to integrate into existing applications.
The server listens to a dedicated queue (one per device) and executes tasks based on their **type (sql, function, command)** and sends the result back.
This is ideal for securely controlling remote devices behind NAT or firewalls without exposing their database or SSH directly, using RabbitMQ as the transport layer.

It provides:  
✅ A client compatible with Go’s `database/sql`.  
✅ A server (device) that listens on RabbitMQ and executes tasks.  
✅ RabbitMQ (AMQP 0-9-1) as the transport.  
✅ Support for 3 task types: `sql`, `function`, `command`.  
✅ Server configurable in `open` and `close` modes with optional connection pool configuration.

## Example Use Case

A SaaS platform needs to manage and query databases installed on customer-premises devices, which are behind NAT and cannot be accessed directly.
The platform uses the burrowctl client to send SQL queries over RabbitMQ to the appropriate device.
The burrowctl server running on the device receives the query, executes it against its local database, and returns the results through RabbitMQ — as if the query was executed locally.

---

## 🚀 Features

- Go client and server libraries.
- Client works seamlessly with Go’s `database/sql` interface.
- Asynchronous transport using RabbitMQ.
- JSON-based responses with results or errors.
- Server supports `open` and `close` modes for database connections.
- Optional connection pool tuning for `open` mode.
- Easily extensible for additional task types.

---

## 📦 Installation

```bash
git clone https://github.com/lordbasex/burrowctl.git
cd burrowctl
go mod tidy
```

---

## 🧪 Requirements

- Running RabbitMQ server.
- MySQL/MariaDB on the device.
- Go >= 1.22.0.

---

## 📄 Client

The client acts as a `database/sql` driver.

### Example

```go
package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lordbasex/burrowctl/client"
)

func main() {
	// DSN con credenciales hardcodeadas para RabbitMQ
	dsn := "deviceID=fd1825ec5a7b63f3fa2be9e04154d3b16f676663ba38e23d4ffafa7b0df29efb&amqp_uri=amqp://burrowuser:burrowpass123@localhost:5672/&timeout=5s&debug=true"

	// Abrir conexión usando el driver rabbitsql
	db, err := sql.Open("rabbitsql", dsn)
	if err != nil {
		log.Fatal("Error connecting:", err)
	}
	defer db.Close()

	log.Println("Executing query SELECT id, name FROM users...")

	// Ejecutar query
	rows, err := db.Query("SELECT id, name FROM users")
	if err != nil {
		log.Fatal("Error executing query:", err)
	}
	defer rows.Close()

	fmt.Println("\n--- Results ---")
	fmt.Printf("%-5s %-30s\n", "ID", "Nombre")
	fmt.Println("------------------------------------")

	// Procesar resultados
	for rows.Next() {
		var id int
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			log.Fatal("Error scanning result:", err)
		}
		fmt.Printf("%-5d %-30s\n", id, name)
	}

	if err := rows.Err(); err != nil {
		log.Fatal("Error iterating results:", err)
	}

	fmt.Println("\n✅ Query completed successfully")
}
```

---

## 📄 Server (device)

The server runs on the remote device.  
It subscribes to its `deviceID` queue and executes incoming tasks.

You can configure the server in two modes:

### 🔗 `open` (default)
Keeps a pool of connections open to the database for better performance.  
You can also customize the pool configuration with:

```go
pool := &server.PoolConfig{
	MaxIdleConns:    5,
	MaxOpenConns:    15,
	ConnMaxLifetime: 5 * time.Minute,
}
```

If no `PoolConfig` is provided or any value is zero, defaults are used:
```
MaxIdleConns:    10
MaxOpenConns:    20
ConnMaxLifetime: 3m
```

### 🔗 `close`
Opens a new connection for each query and closes it after. Safer but slower.

### Example

```go
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/lordbasex/burrowctl/server"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Configurar señales para cerrar gracefully
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c
		log.Println("Closing server...")
		cancel()
	}()

	// Configuración del pool de conexiones
	pool := &server.PoolConfig{
		MaxIdleConns:    5,
		MaxOpenConns:    15,
		ConnMaxLifetime: 5 * time.Minute,
	}

	// Crear el handler con credenciales hardcodeadas
	h := server.NewHandler(
		"fd1825ec5a7b63f3fa2be9e04154d3b16f676663ba38e23d4ffafa7b0df29efb",     // Device ID
		"amqp://burrowuser:burrowpass123@localhost:5672/",                      // RabbitMQ URI
		"burrowuser:burrowpass123@tcp(localhost:3306)/burrowdb?parseTime=true", // MariaDB DSN
		"open", // Modo de conexión: "open" para pool de conexiones
		pool,   // Configuración del pool
	)

	log.Println("Iniciando servidor burrowctl...")
	log.Println("Device ID: fd1825ec5a7b63f3fa2be9e04154d3b16f676663ba38e23d4ffafa7b0df29efb")
	log.Println("RabbitMQ: localhost:5672")
	log.Println("MariaDB: localhost:3306/burrowdb")

	if err := h.Start(ctx); err != nil {
		log.Fatal("Error starting server:", err)
	}

	log.Println("Server closed")
}
```

---

## 📖 Supported task types

### `sql`
Executes the SQL query on the configured database.  
Returns rows and columns as a result.

### `function`
Currently returns a mock response:
```json
{
  "columns": ["message"],
  "rows": [["function executed (mock)"]]
}
```

### `command`
Currently returns a mock response:
```json
{
  "columns": ["message"],
  "rows": [["command executed (mock)"]]
}
```

---

## 📄 Request/Response format

### Request
```json
{
  "type": "sql",
  "deviceID": "<device-id>",
  "query": "SELECT * FROM users",
  "params": []
}
```

### Response
```json
{
  "columns": ["id", "name"],
  "rows": [
    [1, "Alice"],
    [2, "Bob"]
  ],
  "error": ""
}
```

---

## 🐇 Project structure

- Client: `client/`
- Server: `server/`
- Examples: `examples/`

---

## 📋 License

MIT License.
