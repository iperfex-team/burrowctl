/**
 * Prueba simple del cliente BurrowCtl - Una sola query
 */

import { createClient } from './index';

const DSN = 'deviceID=fd1825ec5a7b63f3fa2be9e04154d3b16f676663ba38e23d4ffafa7b0df29efb&amqp_uri=amqp://burrowuser:burrowpass123@localhost:5672/&timeout=5s&debug=true';

async function pruebaSimple() {
  console.log('🧪 Prueba Simple del Cliente BurrowCtl\n');

  try {
    console.log('🔌 Conectando...');
    const client = await createClient(DSN);
    
    console.log('✅ Conectado exitosamente!');
    
    console.log('\n📋 Ejecutando query de prueba...');
    const rows = await client.query('SELECT id, name FROM users');
    
    console.log('\n📊 Resultados:');
    console.log('Columnas:', rows.getColumns());
    console.log('Datos:', rows.getRows());
    console.log('Total filas:', rows.length());
    
    console.log('\n🔌 Cerrando conexión...');
    await client.close();
    
    console.log('✅ ¡Prueba completada exitosamente!');
    
  } catch (error) {
    console.error('❌ Error:', error);
    process.exit(1);
  }
}

// Ejecutar si es el archivo principal
if (require.main === module) {
  pruebaSimple();
}

export { pruebaSimple }; 