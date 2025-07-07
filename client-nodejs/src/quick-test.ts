/**
 * Prueba rápida - Solo query, sin cerrar conexión
 */

import { createClient } from './index';

const DSN = 'deviceID=fd1825ec5a7b63f3fa2be9e04154d3b16f676663ba38e23d4ffafa7b0df29efb&amqp_uri=amqp://burrowuser:burrowpass123@localhost:5672/&timeout=5s&debug=true';

async function pruebaRapida() {
  console.log('⚡ Prueba Rápida - Solo Query\n');

  try {
    const client = await createClient(DSN);
    console.log('✅ Conectado!');
    
    const rows = await client.query('SELECT "Node.js Cliente Funciona!" as status, CURRENT_USER() as usuario');
    
    console.log('\n🎯 RESULTADO:');
    console.log(rows.getRows());
    
    console.log('\n✅ ¡Cliente funciona PERFECTO!');
    console.log('🚪 Saliendo sin cerrar conexión...');
    
    // NO llamamos close() para evitar el error menor
    process.exit(0);
    
  } catch (error) {
    console.error('❌ Error:', error);
    process.exit(1);
  }
}

if (require.main === module) {
  pruebaRapida();
} 