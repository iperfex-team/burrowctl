import { DSNConfig } from './types';

/**
 * Parsea un DSN (Data Source Name) con parámetros de conexión
 * Formato esperado: "deviceID=device123&amqp_uri=amqp://user:pass@host:port&timeout=5s&debug=true"
 */
export function parseDSN(dsn: string): DSNConfig {
  // Parsear manualmente los query parameters
  const params = parseQueryParams(dsn);

  // Verificar parámetros requeridos
  const deviceID = params.get('deviceID');
  if (!deviceID) {
    throw new Error("Falta el parámetro requerido 'deviceID' en el DSN");
  }

  const amqpURI = params.get('amqp_uri');
  if (!amqpURI) {
    throw new Error("Falta el parámetro requerido 'amqp_uri' en el DSN");
  }

  // Validar que amqp_uri tenga el formato correcto
  if (!amqpURI.startsWith('amqp://')) {
    throw new Error("Formato inválido de amqp_uri: debe comenzar con 'amqp://'");
  }

  // Parsear timeout (opcional, valor por defecto: 5 segundos)
  const timeoutStr = params.get('timeout');
  let timeout = 5000; // 5 segundos en milisegundos
  if (timeoutStr) {
    timeout = parseDuration(timeoutStr);
  }

  // Parsear debug (opcional)
  const debugStr = params.get('debug')?.toLowerCase();
  const debug = debugStr === 'true' || debugStr === '1';

  return {
    deviceID,
    amqpURL: amqpURI,
    timeout,
    debug
  };
}

/**
 * Convierte una duración en formato string (ej: "5s", "30s", "1m") a milisegundos
 */
function parseDuration(duration: string): number {
  const regex = /^(\d+(?:\.\d+)?)([smh]?)$/;
  const match = duration.match(regex);
  
  if (!match) {
    throw new Error(`Formato inválido de timeout '${duration}': debe ser un número seguido opcionalmente de 's', 'm', o 'h' (ejemplo: '5s', '30s', '1m')`);
  }

  const value = parseFloat(match[1]);
  const unit = match[2] || 's'; // por defecto segundos

  switch (unit) {
    case 's':
      return value * 1000; // segundos a milisegundos
    case 'm':
      return value * 60 * 1000; // minutos a milisegundos
    case 'h':
      return value * 60 * 60 * 1000; // horas a milisegundos
    default:
      throw new Error(`Unidad de tiempo desconocida: ${unit}`);
  }
}

/**
 * Parsea query parameters de un string
 */
function parseQueryParams(queryString: string): Map<string, string> {
  const params = new Map<string, string>();
  const pairs = queryString.split('&');
  
  for (const pair of pairs) {
    const [key, value] = pair.split('=');
    if (key && value !== undefined) {
      params.set(decodeURIComponent(key), decodeURIComponent(value));
    }
  }
  
  return params;
} 