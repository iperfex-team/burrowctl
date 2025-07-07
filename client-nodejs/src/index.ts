/**
 * Cliente BurrowCtl para Node.js/TypeScript
 * 
 * Este cliente permite conectarse a un servidor BurrowCtl usando RabbitMQ como protocolo de comunicación
 * para ejecutar queries SQL remotas.
 * 
 * @example
 * ```typescript
 * import { BurrowClient, createClient } from 'burrowctl-client-nodejs';
 * 
 * // Usando la clase directamente
 * const client = new BurrowClient();
 * await client.connect('deviceID=my-device&amqp_uri=amqp://localhost&timeout=10s&debug=true');
 * const rows = await client.query('SELECT * FROM users WHERE id = ?', [123]);
 * console.log(rows.getRows());
 * await client.close();
 * 
 * // Usando la función de conveniencia
 * const client2 = await createClient('deviceID=my-device&amqp_uri=amqp://localhost');
 * const rows2 = await client2.query('SELECT COUNT(*) as total FROM products');
 * console.log(rows2.getRows());
 * await client2.close();
 * ```
 */

// Exportar todas las clases principales
export { BurrowClient, createClient } from './client';
export { BurrowConnection } from './connection';
export { SimpleBurrowClient, createSimpleClient } from './simple-client';
export { Rows } from './rows';

// Exportar funciones de utilidad
export { parseDSN } from './config';

// Exportar tipos e interfaces
export type {
  DSNConfig,
  RPCRequest,
  RPCResponse,
  QueryResult,
  QueryValue,
  QueryOptions
} from './types';

// Información de versión
export const VERSION = '1.0.0'; 