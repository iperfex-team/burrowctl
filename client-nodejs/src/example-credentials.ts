/**
 * Ejemplo con las credenciales reales de BurrowCtl
 * Este archivo muestra cómo conectarse usando las credenciales específicas
 */

import { BurrowClient, createClient } from './index';

// Las credenciales reales del usuario
const REAL_DSN = 'deviceID=fd1825ec5a7b63f3fa2be9e04154d3b16f676663ba38e23d4ffafa7b0df29efb&amqp_uri=amqp://burrowuser:burrowpass123@localhost:5672/&timeout=5s&debug=true';

/**
 * Ejemplo usando las credenciales reales del sistema
 */
async function ejemploConCredencialesReales() {
  console.log('🔑 Ejemplo con Credenciales Reales de BurrowCtl\n');

  try {
    // Método 1: Usando createClient (más simple)
    console.log('1. Conectando al servidor BurrowCtl...');
    const client = await createClient(REAL_DSN);
    
    console.log('✅ Conexión establecida exitosamente!');
    
    // Ejemplo de query simple
    console.log('\n2. Ejecutando query de prueba...');
    const rows = await client.query('SELECT id, name FROM users');
    
    console.log('📊 Resultados:');
    console.log('Columnas:', rows.getColumns());
    console.log('Datos:', rows.getRows());
    console.log('Total de filas:', rows.length());
    
    // Ejemplo con parámetros
    console.log('\n3. Ejecutando query con parámetros...');
    const currentTime = new Date().toISOString();
    const rows2 = await client.query('SELECT ? as timestamp, ? as client_type', [currentTime, 'nodejs']);
    
    console.log('📊 Resultados con parámetros:');
    console.log('Datos:', rows2.getRows());
    
    await client.close();
    console.log('\n✅ Ejemplo completado exitosamente!');
    
  } catch (error) {
    console.error('\n❌ Error durante la conexión o ejecución:');
    console.error(error);
    
    // Sugerencias de troubleshooting
    console.log('\n🔧 Troubleshooting:');
    console.log('1. ¿Está RabbitMQ ejecutándose en localhost:5672?');
    console.log('2. ¿Las credenciales burrowuser:burrowpass123 son correctas?');
    console.log('3. ¿El servidor BurrowCtl está conectado a la cola?');
    console.log('4. ¿El deviceID es válido y está activo?');
  }
}

/**
 * Ejemplo usando variables de entorno
 */
async function ejemploConVariablesDeEntorno() {
  console.log('\n🌍 Ejemplo usando Variables de Entorno\n');
  
  // Construir DSN desde variables de entorno
  const deviceID = process.env.BURROW_DEVICE_ID || 'fd1825ec5a7b63f3fa2be9e04154d3b16f676663ba38e23d4ffafa7b0df29efb';
  const amqpURI = process.env.BURROW_AMQP_URI || 'amqp://burrowuser:burrowpass123@localhost:5672/';
  const timeout = process.env.BURROW_TIMEOUT || '5s';
  const debug = process.env.BURROW_DEBUG || 'true';
  
  const envDSN = `deviceID=${deviceID}&amqp_uri=${amqpURI}&timeout=${timeout}&debug=${debug}`;
  
  try {
    const client = new BurrowClient();
    await client.connect(envDSN);
    
    console.log('✅ Conectado usando variables de entorno!');
    
    // Query de información del sistema
    const rows = await client.query('SELECT VERSION() as db_version, NOW() as tiempo_actual');
    console.log('📊 Información del sistema:', rows.getRows());
    
    await client.close();
    
  } catch (error) {
    console.error('❌ Error con variables de entorno:', error);
  }
}

/**
 * Función principal
 */
async function main() {
  console.log('🚀 Ejemplos del Cliente BurrowCtl con Credenciales Reales\n');
  
  await ejemploConCredencialesReales();
  await ejemploConVariablesDeEntorno();
  
  console.log('\n🎉 ¡Todos los ejemplos completados!');
}

// Ejecutar si este es el archivo principal
if (require.main === module) {
  main().catch((error) => {
    console.error('❌ Error fatal:', error);
    process.exit(1);
  });
}

export { ejemploConCredencialesReales, ejemploConVariablesDeEntorno }; 