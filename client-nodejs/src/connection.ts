import * as amqp from 'amqplib';
import { v4 as uuidv4 } from 'uuid';
import * as os from 'os';
import { DSNConfig, RPCRequest, RPCResponse, QueryValue, QueryOptions } from './types';
import { Rows } from './rows';

/**
 * Clase para manejar la conexión y operaciones de base de datos
 * Equivalente a la estructura Conn en Go
 */
export class BurrowConnection {
  private deviceID: string;
  private amqpConnection: any;
  private config: DSNConfig;

  constructor(deviceID: string, amqpConnection: any, config: DSNConfig) {
    this.deviceID = deviceID;
    this.amqpConnection = amqpConnection;
    this.config = config;
  }

  /**
   * Log con debug si está habilitado
   */
  private logf(format: string, ...args: any[]): void {
    if (this.config.debug) {
      const timestamp = new Date().toISOString();
      console.log(`[${timestamp}] [client debug] ${format}`, ...args);
    }
  }

  /**
   * Cierra la conexión a RabbitMQ
   */
  async close(): Promise<void> {
    this.logf('Cerrando conexión a RabbitMQ');
    await this.amqpConnection.close();
  }

  /**
   * Ejecuta una query SQL
   */
  async query(sql: string, params: QueryValue[] = [], options?: QueryOptions): Promise<Rows> {
    const startTotal = Date.now();
    this.logf('Ejecutando query: %s', sql);

    const timeout = options?.timeout || this.config.timeout;
    const rows = await this.queryRPC(sql, params, timeout);

    const total = Date.now() - startTotal;
    this.logf('tiempo total: %d ms', total);

    return rows;
  }

  /**
   * Obtiene la IP local para incluir en la petición
   */
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
   * Ejecuta una query usando RPC sobre RabbitMQ
   * Equivalente a queryRPC en Go
   */
  private async queryRPC(sql: string, params: QueryValue[], timeout: number): Promise<Rows> {
    let channel: any = null;

    try {
      channel = await this.amqpConnection.createChannel();
      this.logf('Canal RabbitMQ abierto');

      // Declarar cola temporal para recibir respuestas
      const replyQueue = await channel.assertQueue('', {
        exclusive: true,
        autoDelete: true
      });
      this.logf('Cola de respuesta declarada: %s', replyQueue.queue);

      const correlationId = uuidv4();

      // Preparar la petición RPC
      const request: RPCRequest = {
        type: 'sql',
        deviceID: this.deviceID,
        query: sql,
        params: params || [],
        clientIP: this.getOutboundIP()
      };

      const requestBody = JSON.stringify(request);

      const startRT = Date.now();
      this.logf("Publicando query a la cola del dispositivo '%s'", this.deviceID);

      // Publicar la petición
      await channel.publish('', this.deviceID, Buffer.from(requestBody), {
        contentType: 'application/json',
        correlationId: correlationId,
        replyTo: replyQueue.queue
      });

      this.logf('Query publicada, esperando respuesta...');

      // Esperar respuesta con timeout
      const response = await this.waitForResponse(channel, replyQueue.queue, correlationId, timeout);
      
      const rt = Date.now() - startRT;
      this.logf('Tiempo de roundtrip RabbitMQ: %d ms', rt);

      // Validar correlationId
      if (response.properties.correlationId !== correlationId) {
        throw new Error(`Error de correlación: esperado ${correlationId}, recibido ${response.properties.correlationId}`);
      }

      // Parsear respuesta
      const resp: RPCResponse = JSON.parse(response.content.toString());
      
      if (resp.error) {
        throw new Error(`Error del servidor: ${resp.error}`);
      }

      this.logf('Respuesta recibida con %d filas', resp.rows.length);
      return new Rows(resp.columns, resp.rows);

    } catch (error) {
      this.logf('Error en queryRPC: %s', error);
      throw error;
    } finally {
      // Cerrar canal de forma segura
      if (channel) {
        try {
          await channel.close();
          this.logf('Canal cerrado correctamente');
        } catch (closeError) {
          this.logf('Error al cerrar canal (no crítico): %s', closeError);
        }
      }
    }
  }

  /**
   * Espera por la respuesta RPC con timeout
   */
  private async waitForResponse(
    channel: any, 
    replyQueue: string, 
    correlationId: string, 
    timeout: number
  ): Promise<any> {
    return new Promise(async (resolve, reject) => {
      let consumerTag: string | null = null;
      
      const timeoutId = setTimeout(() => {
        if (consumerTag) {
          channel.cancel(consumerTag).catch(() => {}); // Cancelar consumer silenciosamente
        }
        reject(new Error(
          `Timeout (${timeout}ms) esperando respuesta del dispositivo '${this.deviceID}'\n` +
          `Por favor verifica:\n` +
          `- El servidor está ejecutándose y respondiendo\n` +
          `- El Device ID '${this.deviceID}' es correcto\n` +
          `- La base de datos es accesible`
        ));
      }, timeout);

      try {
        // Consumir mensajes de la cola de respuesta
        const result = await channel.consume(replyQueue, (msg: any) => {
          if (msg) {
            clearTimeout(timeoutId);
            if (consumerTag) {
              channel.cancel(consumerTag).catch(() => {}); // Cancelar consumer silenciosamente
            }
            resolve(msg);
          }
        }, { noAck: true, exclusive: true });
        
        consumerTag = result.consumerTag;
        
      } catch (error: any) {
        clearTimeout(timeoutId);
        reject(new Error(`Error al consumir de la cola de respuesta: ${error.message}`));
      }
    });
  }
} 