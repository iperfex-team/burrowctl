/**
 * Cliente Simple BurrowCtl - Una conexión por query
 * Esta versión elimina completamente los problemas de canales múltiples
 */

import * as amqp from 'amqplib';
import { v4 as uuidv4 } from 'uuid';
import * as os from 'os';
import { DSNConfig, RPCRequest, RPCResponse, QueryValue } from './types';
import { Rows } from './rows';
import { parseDSN } from './config';

/**
 * Cliente simple que crea una nueva conexión para cada query
 * Más lento pero 100% confiable para queries múltiples
 */
export class SimpleBurrowClient {
  private config: DSNConfig;

  constructor(dsn: string) {
    this.config = parseDSN(dsn);
  }

  private logf(format: string, ...args: any[]): void {
    if (this.config.debug) {
      const timestamp = new Date().toISOString();
      console.log(`[${timestamp}] [simple client] ${format}`, ...args);
    }
  }

  private getOutboundIP(): string {
    try {
      const interfaces = os.networkInterfaces();
      for (const name of Object.keys(interfaces)) {
        const iface = interfaces[name];
        if (iface) {
          for (const alias of iface) {
            if (alias.family === 'IPv4' && !alias.internal) {
              return alias.address;
            }
          }
        }
      }
      return 'unknown';
    } catch (error) {
      return 'unknown';
    }
  }

  /**
   * Ejecuta una query SQL - cada query usa su propia conexión
   */
  async query(sql: string, params: QueryValue[] = []): Promise<Rows> {
    const startTotal = Date.now();
    this.logf('Ejecutando query: %s', sql);

    let connection: any = null;
    let channel: any = null;

    try {
      // Nueva conexión para esta query
      connection = await amqp.connect(this.config.amqpURL);
      this.logf('Nueva conexión establecida');

      channel = await connection.createChannel();
      this.logf('Canal creado');

      // Cola temporal para respuestas
      const replyQueue = await channel.assertQueue('', {
        exclusive: true,
        autoDelete: true
      });
      this.logf('Cola de respuesta: %s', replyQueue.queue);

      const correlationId = uuidv4();

      // Preparar petición RPC
      const request: RPCRequest = {
        type: 'sql',
        deviceID: this.config.deviceID,
        query: sql,
        params: params || [],
        clientIP: this.getOutboundIP()
      };

      const requestBody = JSON.stringify(request);

      const startRT = Date.now();
      this.logf("Enviando query al dispositivo '%s'", this.config.deviceID);

      // Publicar petición
      await channel.publish('', this.config.deviceID, Buffer.from(requestBody), {
        contentType: 'application/json',
        correlationId: correlationId,
        replyTo: replyQueue.queue
      });

      this.logf('Query enviada, esperando respuesta...');

      // Esperar respuesta
      const response = await this.waitForSingleResponse(channel, replyQueue.queue, correlationId);
      
      const rt = Date.now() - startRT;
      this.logf('Roundtrip: %d ms', rt);

      // Validar correlación
      if (response.properties.correlationId !== correlationId) {
        throw new Error(`Error de correlación: esperado ${correlationId}, recibido ${response.properties.correlationId}`);
      }

      // Parsear respuesta
      const resp: RPCResponse = JSON.parse(response.content.toString());
      
      if (resp.error) {
        throw new Error(`Error del servidor: ${resp.error}`);
      }

      const total = Date.now() - startTotal;
      this.logf('Query completada en %d ms con %d filas', total, resp.rows.length);
      
      return new Rows(resp.columns, resp.rows);

    } finally {
      // Limpiar recursos
      if (channel) {
        try {
          await channel.close();
          this.logf('Canal cerrado');
        } catch (e) {
          this.logf('Error cerrando canal: %s', e);
        }
      }
      
      if (connection) {
        try {
          await connection.close();
          this.logf('Conexión cerrada');
        } catch (e) {
          this.logf('Error cerrando conexión: %s', e);
        }
      }
    }
  }

  /**
   * Espera por una respuesta RPC usando Promise simple
   */
  private async waitForSingleResponse(
    channel: any,
    replyQueue: string,
    correlationId: string
  ): Promise<any> {
    return new Promise((resolve, reject) => {
      const timeoutId = setTimeout(() => {
        reject(new Error(`Timeout esperando respuesta del dispositivo '${this.config.deviceID}'`));
      }, this.config.timeout);

      channel.consume(replyQueue, (msg: any) => {
        if (msg) {
          clearTimeout(timeoutId);
          resolve(msg);
        }
      }, { noAck: true })
      .catch((error: any) => {
        clearTimeout(timeoutId);
        reject(error);
      });
    });
  }
}

/**
 * Función de conveniencia para crear un cliente simple
 */
export async function createSimpleClient(dsn: string): Promise<SimpleBurrowClient> {
  return new SimpleBurrowClient(dsn);
} 