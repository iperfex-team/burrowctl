import * as amqp from 'amqplib';
import { DSNConfig, QueryValue, QueryOptions } from './types';
import { parseDSN } from './config';
import { BurrowConnection } from './connection';
import { Rows } from './rows';

/**
 * Cliente principal para BurrowCtl
 * Equivalente al Driver en Go
 */
export class BurrowClient {
  private connection: BurrowConnection | null = null;

  /**
   * Conecta al servidor usando un DSN
   * Formato: "deviceID=device123&amqp_uri=amqp://user:pass@host:port&timeout=5s&debug=true"
   */
  async connect(dsn: string): Promise<void> {
    const config = parseDSN(dsn);
    
    try {
      // Intentar conectar a RabbitMQ
      const amqpConnection = await amqp.connect(config.amqpURL) as any;
      
      if (config.debug) {
        console.log(`[client debug] Conectado a RabbitMQ ${config.amqpURL} (deviceID=${config.deviceID}, timeout=${config.timeout}ms)`);
      }

      this.connection = new BurrowConnection(config.deviceID, amqpConnection, config);
    } catch (error: any) {
      throw new Error(
        `Conexión a RabbitMQ falló '${config.amqpURL}': ${error.message}\n` +
        `Por favor verifica:\n` +
        `- El servidor RabbitMQ está ejecutándose\n` +
        `- Las credenciales son correctas\n` +
        `- La conectividad de red`
      );
    }
  }

  /**
   * Ejecuta una query SQL
   */
  async query(sql: string, params: QueryValue[] = [], options?: QueryOptions): Promise<Rows> {
    if (!this.connection) {
      throw new Error('No hay conexión activa. Llama a connect() primero.');
    }
    
    return await this.connection.query(sql, params, options);
  }

  /**
   * Cierra la conexión
   */
  async close(): Promise<void> {
    if (this.connection) {
      await this.connection.close();
      this.connection = null;
    }
  }

  /**
   * Verifica si hay una conexión activa
   */
  isConnected(): boolean {
    return this.connection !== null;
  }
}

/**
 * Función de conveniencia para crear una conexión rápida
 */
export async function createClient(dsn: string): Promise<BurrowClient> {
  const client = new BurrowClient();
  await client.connect(dsn);
  return client;
} 