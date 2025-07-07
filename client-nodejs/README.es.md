# Cliente Node.js de BurrowCtl

Una biblioteca cliente de Node.js/TypeScript para BurrowCtl que permite la ejecución remota de consultas SQL usando el protocolo RPC de RabbitMQ.

## Características

- 🚀 **Fácil de usar**: API simple para conectar y ejecutar consultas SQL
- 📡 **RabbitMQ RPC**: Usa RabbitMQ como protocolo de comunicación
- 🔒 **Seguridad de tipos**: Escrito en TypeScript con definiciones de tipos completas
- ⚡ **Flexible**: Soporte para consultas parametrizadas y timeouts personalizados
- 🛠️ **Manejo de errores**: Gestión de errores integral y opciones de depuración
- 📦 **Liviano**: Dependencias mínimas

## Instalación

```bash
npm install burrowctl-client-nodejs
```

## Inicio Rápido

### Uso Básico

```typescript
import { BurrowClient, createClient } from 'burrowctl-client-nodejs';

// Método 1: Usando la clase BurrowClient directamente
const client = new BurrowClient();
await client.connect('deviceID=mi-dispositivo&amqp_uri=amqp://localhost&timeout=10s&debug=true');

const rows = await client.query('SELECT * FROM usuarios WHERE id = ?', [123]);
console.log(rows.getRows());

await client.close();

// Método 2: Usando función de conveniencia
const client2 = await createClient('deviceID=mi-dispositivo&amqp_uri=amqp://localhost');
const rows2 = await client2.query('SELECT COUNT(*) as total FROM productos');
console.log(rows2.getRows());
await client2.close();
```

### Formato de Cadena de Conexión

La cadena de conexión (DSN) sigue este formato:
```
deviceID=<identificador-dispositivo>&amqp_uri=<url-rabbitmq>&timeout=<timeout>&debug=<boolean>
```

**Parámetros:**
- `deviceID`: Identificador único para tu dispositivo/cliente
- `amqp_uri`: URL de conexión de RabbitMQ (ej: `amqp://usuario:contraseña@host:puerto`)
- `timeout`: Timeout de consulta (ej: `5s`, `10s`, `30s`)
- `debug`: Habilitar logging de depuración (`true` o `false`)

**Ejemplo:**
```
deviceID=webapp-001&amqp_uri=amqp://guest:guest@localhost:5672&timeout=15s&debug=true
```

## Referencia de la API

### BurrowClient

#### Métodos

##### `connect(dsn: string): Promise<void>`
Se conecta al servidor BurrowCtl usando el DSN proporcionado.

##### `query(sql: string, params?: QueryValue[], options?: QueryOptions): Promise<Rows>`
Ejecuta una consulta SQL con parámetros opcionales.

- `sql`: Cadena de consulta SQL
- `params`: Array de parámetros para consultas parametrizadas
- `options`: Opciones de consulta (timeout, etc.)

##### `close(): Promise<void>`
Cierra la conexión al servidor.

##### `isConnected(): boolean`
Retorna si el cliente está actualmente conectado.

### Rows

La clase `Rows` proporciona métodos para iterar y acceder a los resultados de las consultas.

#### Métodos

##### `getRows(): any[]`
Retorna todas las filas como un array.

##### `getColumns(): string[]`
Retorna los nombres de las columnas.

##### `length(): number`
Retorna el número de filas.

##### `hasNext(): boolean`
Verifica si hay más filas para iterar.

##### `next(): any`
Retorna la siguiente fila y avanza el iterador.

##### `reset(): void`
Reinicia el iterador al principio.

## Ejemplos

### Consulta Simple

```typescript
import { createClient } from 'burrowctl-client-nodejs';

const client = await createClient('deviceID=ejemplo&amqp_uri=amqp://localhost');
const rows = await client.query('SELECT 1 as numero, \'Hola Mundo\' as mensaje');

console.log('Resultados:', rows.getRows());
console.log('Columnas:', rows.getColumns());

await client.close();
```

### Consulta Parametrizada

```typescript
const rows = await client.query(
  'SELECT * FROM usuarios WHERE edad > ? AND ciudad = ?', 
  [25, 'Madrid']
);

console.log(`Se encontraron ${rows.length()} usuarios`);
```

### Timeout Personalizado

```typescript
const rows = await client.query(
  'SELECT * FROM tabla_grande', 
  [], 
  { timeout: 30000 } // 30 segundos
);
```

### Manejo de Errores

```typescript
try {
  const client = new BurrowClient();
  await client.connect('deviceID=test&amqp_uri=amqp://localhost:5672');
  
  const rows = await client.query('SELECT * FROM usuarios');
  console.log(rows.getRows());
  
} catch (error) {
  console.error('Falló la conexión o consulta:', error.message);
} finally {
  await client.close();
}
```

## Requisitos

- Node.js 16+
- Servidor RabbitMQ ejecutándose y accesible
- Servidor BurrowCtl configurado con RabbitMQ

## Desarrollo

### Compilación

```bash
npm run build
```

### Modo Desarrollo

```bash
npm run dev
```

### Linting

```bash
npm run lint
```

## Licencia

MIT

## Soporte

Para preguntas y soporte, por favor consulta la documentación principal de BurrowCtl. 