# BurrowCtl Node.js Client

<div align="right">
  Leia em outros idiomas: 
  <a title="Spanish" href="./README.es.md">🇦🇷</a>
  <a title="Português" href="./README.pt.md">🇧🇷</a>
</div>

A Node.js/TypeScript client library for BurrowCtl that enables remote SQL query execution using RabbitMQ RPC protocol.

## Features

- 🚀 **Easy to use**: Simple API for connecting and executing SQL queries
- 📡 **RabbitMQ RPC**: Uses RabbitMQ as the communication protocol
- 🔒 **Type-safe**: Written in TypeScript with full type definitions
- ⚡ **Flexible**: Support for parameterized queries and custom timeouts
- 🛠️ **Error handling**: Comprehensive error management and debugging options
- 📦 **Lightweight**: Minimal dependencies

## Installation

```bash
npm install burrowctl-client-nodejs
```

## Quick Start

### Basic Usage

```typescript
import { BurrowClient, createClient } from 'burrowctl-client-nodejs';

// Method 1: Using BurrowClient class directly
const client = new BurrowClient();
await client.connect('deviceID=my-device&amqp_uri=amqp://localhost&timeout=10s&debug=true');

const rows = await client.query('SELECT * FROM users WHERE id = ?', [123]);
console.log(rows.getRows());

await client.close();

// Method 2: Using convenience function
const client2 = await createClient('deviceID=my-device&amqp_uri=amqp://localhost');
const rows2 = await client2.query('SELECT COUNT(*) as total FROM products');
console.log(rows2.getRows());
await client2.close();
```

### Connection String Format

The connection string (DSN) follows this format:
```
deviceID=<device-identifier>&amqp_uri=<rabbitmq-url>&timeout=<timeout>&debug=<boolean>
```

**Parameters:**
- `deviceID`: Unique identifier for your device/client
- `amqp_uri`: RabbitMQ connection URL (e.g., `amqp://user:pass@host:port`)
- `timeout`: Query timeout (e.g., `5s`, `10s`, `30s`)
- `debug`: Enable debug logging (`true` or `false`)

**Example:**
```
deviceID=webapp-001&amqp_uri=amqp://guest:guest@localhost:5672&timeout=15s&debug=true
```

## API Reference

### BurrowClient

#### Methods

##### `connect(dsn: string): Promise<void>`
Connects to the BurrowCtl server using the provided DSN.

##### `query(sql: string, params?: QueryValue[], options?: QueryOptions): Promise<Rows>`
Executes a SQL query with optional parameters.

- `sql`: SQL query string
- `params`: Array of parameters for parameterized queries
- `options`: Query options (timeout, etc.)

##### `close(): Promise<void>`
Closes the connection to the server.

##### `isConnected(): boolean`
Returns whether the client is currently connected.

### Rows

The `Rows` class provides methods to iterate and access query results.

#### Methods

##### `getRows(): any[]`
Returns all rows as an array.

##### `getColumns(): string[]`
Returns column names.

##### `length(): number`
Returns the number of rows.

##### `hasNext(): boolean`
Checks if there are more rows to iterate.

##### `next(): any`
Returns the next row and advances the iterator.

##### `reset(): void`
Resets the iterator to the beginning.

## Examples

### Simple Query

```typescript
import { createClient } from 'burrowctl-client-nodejs';

const client = await createClient('deviceID=example&amqp_uri=amqp://localhost');
const rows = await client.query('SELECT 1 as number, \'Hello World\' as message');

console.log('Results:', rows.getRows());
console.log('Columns:', rows.getColumns());

await client.close();
```

### Parameterized Query

```typescript
const rows = await client.query(
  'SELECT * FROM users WHERE age > ? AND city = ?', 
  [25, 'New York']
);

console.log(`Found ${rows.length()} users`);
```

### Custom Timeout

```typescript
const rows = await client.query(
  'SELECT * FROM large_table', 
  [], 
  { timeout: 30000 } // 30 seconds
);
```

### Error Handling

```typescript
try {
  const client = new BurrowClient();
  await client.connect('deviceID=test&amqp_uri=amqp://localhost:5672');
  
  const rows = await client.query('SELECT * FROM users');
  console.log(rows.getRows());
  
} catch (error) {
  console.error('Connection or query failed:', error.message);
} finally {
  await client.close();
}
```

## Requirements

- Node.js 16+ 
- RabbitMQ server running and accessible
- BurrowCtl server configured with RabbitMQ

## Development

### Building

```bash
npm run build
```

### Development Mode

```bash
npm run dev
```

### Linting

```bash
npm run lint
```

## License

MIT

## Support

For questions and support, please refer to the main BurrowCtl documentation. 