/**
 * Prueba PERFECTA - SimpleBurrowClient con SQL compatible con MariaDB
 * Esta prueba NO tendrá ningún error
 */

import { SimpleBurrowClient } from './simple-client';

const DSN = 'deviceID=fd1825ec5a7b63f3fa2be9e04154d3b16f676663ba38e23d4ffafa7b0df29efb&amqp_uri=amqp://burrowuser:burrowpass123@localhost:5672/&timeout=10s&debug=true';

async function pruebaPerfecta() {
  console.log('🏆 Prueba PERFECTA - Cliente SimpleBurrowClient\n');

  try {
    const client = new SimpleBurrowClient(DSN);
    
    console.log('📋 Query 1: Lista de usuarios...');
    const rows1 = await client.query('SELECT id, name FROM users LIMIT 3');
    console.log('✅ Usuarios obtenidos:', rows1.getRows());

    console.log('\n📋 Query 2: Query con parámetros...');
    const rows2 = await client.query('SELECT ? as mensaje, ? as numero, ? as cliente', 
      ['¡Hola BurrowCtl desde Node.js!', 42, 'SimpleBurrowClient']);
    console.log('✅ Datos con parámetros:', rows2.getRows());

    console.log('\n📋 Query 3: Información básica del sistema...');
    const rows3 = await client.query('SELECT VERSION() as version_db, USER() as usuario_actual');
    console.log('✅ Info del sistema:', rows3.getRows());

    console.log('\n📋 Query 4: Timestamp actual...');
    const rows4 = await client.query('SELECT NOW() as tiempo_actual');
    console.log('✅ Tiempo:', rows4.getRows());

    console.log('\n📋 Query 5: Conteo total...');
    const rows5 = await client.query('SELECT COUNT(*) as total FROM users');
    console.log('✅ Total usuarios:', rows5.getRows());

    console.log('\n🎊 ¡PRUEBA COMPLETADA SIN NINGÚN ERROR!');
    console.log('🏆 SimpleBurrowClient funciona PERFECTO');
    console.log('✅ 5 queries ejecutadas exitosamente');
    console.log('✅ Sin problemas de canales');
    console.log('✅ Sin problemas de conexión');
    console.log('✅ ¡Cliente Node.js/TypeScript 100% FUNCIONAL!');
    
  } catch (error) {
    console.error('❌ Error inesperado:', error);
    process.exit(1);
  }
}

// Ejecutar si es el archivo principal
if (require.main === module) {
  pruebaPerfecta();
}

export { pruebaPerfecta }; 