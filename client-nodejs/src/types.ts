/**
 * Tipos e interfaces para el cliente BurrowCtl Node.js
 */

export interface DSNConfig {
  deviceID: string;
  amqpURL: string;
  timeout: number; // en milisegundos
  debug: boolean;
}

export interface RPCRequest {
  type: string;
  deviceID: string;
  query: string;
  params: any[];
  clientIP: string;
}

export interface RPCResponse {
  columns: string[];
  rows: any[][];
  error?: string;
}

export interface QueryResult {
  columns: string[];
  rows: any[][];
}

export type QueryValue = string | number | boolean | null | undefined;

export interface QueryOptions {
  timeout?: number;
} 