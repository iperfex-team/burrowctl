# Cliente Node.js do BurrowCtl

Uma biblioteca cliente Node.js/TypeScript para BurrowCtl que permite execução remota de consultas SQL usando o protocolo RPC do RabbitMQ.

## Recursos

- 🚀 **Fácil de usar**: API simples para conectar e executar consultas SQL
- 📡 **RabbitMQ RPC**: Usa RabbitMQ como protocolo de comunicação
- 🔒 **Tipagem segura**: Escrito em TypeScript com definições de tipos completas
- ⚡ **Flexível**: Suporte para consultas parametrizadas e timeouts personalizados
- 🛠️ **Tratamento de erros**: Gestão de erros abrangente e opções de debug
- 📦 **Leve**: Dependências mínimas

## Instalação

```bash
npm install burrowctl-client-nodejs
```

## Início Rápido

### Uso Básico

```typescript
import { BurrowClient, createClient } from 'burrowctl-client-nodejs';

// Método 1: Usando a classe BurrowClient diretamente
const client = new BurrowClient();
await client.connect('deviceID=meu-dispositivo&amqp_uri=amqp://localhost&timeout=10s&debug=true');

const rows = await client.query('SELECT * FROM usuarios WHERE id = ?', [123]);
console.log(rows.getRows());

await client.close();

// Método 2: Usando função de conveniência
const client2 = await createClient('deviceID=meu-dispositivo&amqp_uri=amqp://localhost');
const rows2 = await client2.query('SELECT COUNT(*) as total FROM produtos');
console.log(rows2.getRows());
await client2.close();
```

### Formato da String de Conexão

A string de conexão (DSN) segue este formato:
```
deviceID=<identificador-dispositivo>&amqp_uri=<url-rabbitmq>&timeout=<timeout>&debug=<boolean>
```

**Parâmetros:**
- `deviceID`: Identificador único para seu dispositivo/cliente
- `amqp_uri`: URL de conexão do RabbitMQ (ex: `amqp://usuario:senha@host:porta`)
- `timeout`: Timeout da consulta (ex: `5s`, `10s`, `30s`)
- `debug`: Habilitar log de debug (`true` ou `false`)

**Exemplo:**
```
deviceID=webapp-001&amqp_uri=amqp://guest:guest@localhost:5672&timeout=15s&debug=true
```

## Referência da API

### BurrowClient

#### Métodos

##### `connect(dsn: string): Promise<void>`
Conecta ao servidor BurrowCtl usando o DSN fornecido.

##### `query(sql: string, params?: QueryValue[], options?: QueryOptions): Promise<Rows>`
Executa uma consulta SQL com parâmetros opcionais.

- `sql`: String da consulta SQL
- `params`: Array de parâmetros para consultas parametrizadas
- `options`: Opções de consulta (timeout, etc.)

##### `close(): Promise<void>`
Fecha a conexão com o servidor.

##### `isConnected(): boolean`
Retorna se o cliente está atualmente conectado.

### Rows

A classe `Rows` fornece métodos para iterar e acessar resultados de consultas.

#### Métodos

##### `getRows(): any[]`
Retorna todas as linhas como um array.

##### `getColumns(): string[]`
Retorna os nomes das colunas.

##### `length(): number`
Retorna o número de linhas.

##### `hasNext(): boolean`
Verifica se há mais linhas para iterar.

##### `next(): any`
Retorna a próxima linha e avança o iterador.

##### `reset(): void`
Reinicia o iterador para o início.

## Exemplos

### Consulta Simples

```typescript
import { createClient } from 'burrowctl-client-nodejs';

const client = await createClient('deviceID=exemplo&amqp_uri=amqp://localhost');
const rows = await client.query('SELECT 1 as numero, \'Olá Mundo\' as mensagem');

console.log('Resultados:', rows.getRows());
console.log('Colunas:', rows.getColumns());

await client.close();
```

### Consulta Parametrizada

```typescript
const rows = await client.query(
  'SELECT * FROM usuarios WHERE idade > ? AND cidade = ?', 
  [25, 'São Paulo']
);

console.log(`Encontrados ${rows.length()} usuários`);
```

### Timeout Personalizado

```typescript
const rows = await client.query(
  'SELECT * FROM tabela_grande', 
  [], 
  { timeout: 30000 } // 30 segundos
);
```

### Tratamento de Erros

```typescript
try {
  const client = new BurrowClient();
  await client.connect('deviceID=test&amqp_uri=amqp://localhost:5672');
  
  const rows = await client.query('SELECT * FROM usuarios');
  console.log(rows.getRows());
  
} catch (error) {
  console.error('Conexão ou consulta falhou:', error.message);
} finally {
  await client.close();
}
```

## Requisitos

- Node.js 16+
- Servidor RabbitMQ executando e acessível
- Servidor BurrowCtl configurado com RabbitMQ

## Desenvolvimento

### Compilação

```bash
npm run build
```

### Modo Desenvolvimento

```bash
npm run dev
```

### Linting

```bash
npm run lint
```

## Licença

MIT

## Suporte

Para perguntas e suporte, consulte a documentação principal do BurrowCtl. 