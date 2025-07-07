/**
 * Prueba de múltiples queries usando SimpleBurrowClient
 * Esta versión NO tiene problemas de canales
 */

import { SimpleBurrowClient, createSimpleClient } from './simple-client';

const DSN = 'deviceID=fd1825ec5a7b63f3fa2be9e04154d3b16f676663ba38e23d4ffafa7b0df29efb&amqp_uri=amqp://burrowuser:burrowpass123@localhost:5672/&timeout=10s&debug=true';

async function pruebaMultiplesQueries() {
  console.log('🔥 Prueba de Múltiples Queries - Cliente Simple\n');

  try {
    const client = new SimpleBurrowClient(DSN);
    
    console.log('📋 Query 1: Obtener usuarios...');
    const rows1 = await client.query('SELECT id, name FROM users LIMIT 5');
    console.log('✅ Resultado 1:', rows1.getRows());
    console.log('📊 Total filas:', rows1.length());

    console.log('\n📋 Query 2: Query con parámetros...');
    const currentTime = new Date().toISOString();
    const rows2 = await client.query('SELECT ? as timestamp, ? as client_type, ? as status', 
      [currentTime, 'nodejs-simple', 'success']);
    console.log('✅ Resultado 2:', rows2.getRows());

    console.log('\n📋 Query 3: Información del sistema...');
    const rows3 = await client.query('SELECT VERSION() as version, NOW() as current_time, USER() as user');
    console.log('✅ Resultado 3:', rows3.getRows());

    console.log('\n📋 Query 4: Conteo de usuarios...');
    const rows4 = await client.query('SELECT COUNT(*) as total_users FROM users');
    console.log('✅ Resultado 4:', rows4.getRows());

    console.log('\n🎉 ¡TODAS LAS QUERIES EJECUTADAS EXITOSAMENTE!');
    console.log('✅ Sin errores de canales');
    console.log('✅ Sin problemas de conexión');
    console.log('✅ Cliente SimpleBurrowClient funciona PERFECTO');
    
  } catch (error) {
    console.error('❌ Error:', error);
    process.exit(1);
  }
}

// Ejecutar si es el archivo principal
if (require.main === module) {
  pruebaMultiplesQueries();
}

export { pruebaMultiplesQueries }; 