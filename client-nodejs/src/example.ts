/**
 * Ejemplo de uso del cliente BurrowCtl para Node.js
 */

import { BurrowClient, createClient } from './index';

async function exemploBasico() {
  console.log('=== Ejemplo Básico del Cliente BurrowCtl ===\n');

  // Configuración de conexión
  const dsn = 'deviceID=mi-dispositivo&amqp_uri=amqp://guest:guest@localhost:5672&timeout=10s&debug=true';

  try {
    // Método 1: Usando la clase BurrowClient directamente
    console.log('1. Conectando usando BurrowClient...');
    const client = new BurrowClient();
    await client.connect(dsn);
    
    // Ejecutar una query simple
    console.log('2. Ejecutando query simple...');
    const rows1 = await client.query('SELECT 1 as numero, \'Hola Mundo\' as mensaje');
    console.log('Resultados:', rows1.getRows());
    console.log('Columnas:', rows1.getColumns());
    console.log('Número de filas:', rows1.length());

    // Ejecutar query con parámetros
    console.log('\n3. Ejecutando query con parámetros...');
    const rows2 = await client.query('SELECT ? as id, ? as nombre', [123, 'Juan Pérez']);
    console.log('Resultados:', rows2.getRows());

    // Iterar sobre las filas
    console.log('\n4. Iterando sobre las filas...');
    rows2.reset(); // Reiniciar el iterador
    while (rows2.hasNext()) {
      const row = rows2.next();
      console.log('Fila:', row);
    }

    await client.close();

    // Método 2: Usando la función de conveniencia
    console.log('\n5. Usando createClient...');
    const client2 = await createClient(dsn);
    const rows3 = await client2.query('SELECT NOW() as fecha_actual');
    console.log('Fecha actual:', rows3.getRows());
    await client2.close();

    console.log('\n✅ Ejemplo completado exitosamente!');

  } catch (error) {
    console.error('❌ Error:', error);
  }
}

async function ejemploConTimeout() {
  console.log('\n=== Ejemplo con Timeout Personalizado ===\n');

  const dsn = 'deviceID=mi-dispositivo&amqp_uri=amqp://guest:guest@localhost:5672&debug=true';

  try {
    const client = await createClient(dsn);
    
    // Query con timeout personalizado
    const rows = await client.query(
      'SELECT SLEEP(2) as delay, \'Operación lenta\' as mensaje',
      [],
      { timeout: 15000 } // 15 segundos
    );
    
    console.log('Resultado de operación lenta:', rows.getRows());
    await client.close();

  } catch (error) {
    console.error('❌ Error con timeout:', error);
  }
}

async function ejemploManejandoErrores() {
  console.log('\n=== Ejemplo Manejando Errores ===\n');

  try {
    // DSN con configuración incorrecta
    const dsnIncorrecto = 'deviceID=dispositivo-inexistente&amqp_uri=amqp://localhost:9999&timeout=5s';
    const client = new BurrowClient();
    
    console.log('Intentando conectar con configuración incorrecta...');
    await client.connect(dsnIncorrecto);
    
  } catch (error) {
    console.error('✅ Error esperado capturado:', error);
  }

  try {
    // Query sin conexión
    const client2 = new BurrowClient();
    await client2.query('SELECT 1');
    
  } catch (error) {
    console.error('✅ Error esperado capturado:', error);
  }
}

// Ejecutar los ejemplos
async function main() {
  console.log('🚀 Ejecutando ejemplos del Cliente BurrowCtl\n');
  console.log('NOTA: Asegúrate de que el servidor RabbitMQ esté ejecutándose en localhost:5672\n');

  await exemploBasico();
  await ejemploConTimeout();
  await ejemploManejandoErrores();
  
  console.log('\n🎉 Todos los ejemplos completados!');
}

// Ejecutar solo si este archivo es el punto de entrada
if (require.main === module) {
  main().catch(console.error);
}

export { exemploBasico, ejemploConTimeout, ejemploManejandoErrores }; 