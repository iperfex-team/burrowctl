# 🐇 burrowctl - Cliente Node.js

<div align="center">
  <h3>Ejecución Remota de SQL y Control de Dispositivos vía RabbitMQ</h3>
  <p>
    Una biblioteca cliente Node.js/TypeScript para <strong>burrowctl</strong> que permite la ejecución remota de consultas SQL, funciones personalizadas y comandos del sistema usando el protocolo RPC de RabbitMQ.
  </p>
  <p>
    <a href="./README.md">🇺🇸 English</a> | 
    <a href="./README.es.md">🇪🇸 Español</a> | 
    <a href="./README.pt.md">🇧🇷 Português</a>
  </p>
</div>

## 🎯 ¿Qué es burrowctl?

**burrowctl** permite el acceso seguro a bases de datos remotas y control de dispositivos sin exponer conexiones directas. Es perfecto para:

- 🏢 **Plataformas SaaS**: Gestionar bases de datos de clientes detrás de NAT/firewalls
- 🌐 **Gestión IoT**: Controlar dispositivos distribuidos de forma segura
- 🔐 **Administración Remota**: Ejecutar consultas y comandos sin acceso SSH/DB directo
- 📊 **Monitoreo Distribuido**: Recopilar datos de múltiples fuentes remotas

## ✨ Características Principales

### 🚀 **Fácil de usar**
- **API Simple**: Interfaz clara y directa para conectar y ejecutar consultas
- **TypeScript Nativo**: Escrito completamente en TypeScript con definiciones de tipos completas
- **Compatibilidad Universal**: Mismo formato de DSN que el cliente Go

### 📡 **Protocolo RabbitMQ RPC**
- **Transporte Seguro**: Usa RabbitMQ AMQP 0-9-1 como protocolo de comunicación
- **Comunicación Asíncrona**: Manejo eficiente de solicitudes/respuestas
- **Gestión de Colas**: Manejo automático de colas de dispositivos

### 🔒 **Seguridad de Tipos**
- **Definiciones Completas**: Tipos TypeScript para todas las interfaces
- **Validación de Parámetros**: Verificación automática de tipos de datos
- **Manejo de Errores**: Gestión de errores integral con información detallada

### ⚡ **Flexible y Configurable**
- **Consultas Parametrizadas**: Soporte completo para parámetros de consulta
- **Timeouts Personalizados**: Control granular de tiempos de espera
- **Opciones de Depuración**: Logging detallado para desarrollo y producción

### 🛠️ **Características Empresariales**
- **Manejo de Errores**: Gestión robusta de errores con información contextual
- **Dependencias Mínimas**: Biblioteca liviana con dependencias optimizadas
- **Compatibilidad**: Node.js 16+ y navegadores modernos

---

## 🚀 Instalación

### Via NPM
```bash
npm install burrowctl-client-nodejs
```

### Via Yarn
```bash
yarn add burrowctl-client-nodejs
```

### Via PNPM
```bash
pnpm add burrowctl-client-nodejs
```

---

## 🔧 Inicio Rápido

### Uso Básico

```typescript
import { BurrowClient, createClient } from 'burrowctl-client-nodejs';

// Método 1: Usando la clase BurrowClient directamente
const client = new BurrowClient();
await client.connect('deviceID=mi-dispositivo&amqp_uri=amqp://localhost&timeout=10s&debug=true');

const rows = await client.query('SELECT * FROM usuarios WHERE id = ?', [123]);
console.log('Filas:', rows.getRows());
console.log('Columnas:', rows.getColumns());

await client.close();

// Método 2: Usando función de conveniencia
const client2 = await createClient('deviceID=mi-dispositivo&amqp_uri=amqp://localhost');
const rows2 = await client2.query('SELECT COUNT(*) as total FROM productos');
console.log('Total de productos:', rows2.getRows()[0].total);
await client2.close();
```

### Configuración Avanzada

```typescript
import { BurrowClient, QueryOptions } from 'burrowctl-client-nodejs';

const client = new BurrowClient();
await client.connect('deviceID=webapp-001&amqp_uri=amqp://usuario:contraseña@localhost:5672&timeout=30s&debug=true');

// Consulta con opciones personalizadas
const options: QueryOptions = {
    timeout: 45000, // 45 segundos
};

const rows = await client.query(
    'SELECT * FROM ventas WHERE fecha_venta BETWEEN ? AND ?',
    ['2024-01-01', '2024-12-31'],
    options
);

console.log(`Encontradas ${rows.length()} ventas`);
await client.close();
```

---

## 🔗 Formato de Cadena de Conexión (DSN)

La cadena de conexión sigue este formato:
```
deviceID=<identificador-dispositivo>&amqp_uri=<url-rabbitmq>&timeout=<timeout>&debug=<boolean>
```

### Parámetros Detallados

| Parámetro | Descripción | Ejemplo | Requerido |
|-----------|-------------|---------|-----------|
| `deviceID` | Identificador único del dispositivo/cliente | `webapp-001` | ✅ |
| `amqp_uri` | URL de conexión de RabbitMQ | `amqp://usuario:contraseña@host:puerto/vhost` | ✅ |
| `timeout` | Timeout de consulta | `5s`, `30s`, `2m` | ❌ |
| `debug` | Habilitar logging de depuración | `true`, `false` | ❌ |

### Ejemplos de DSN

```typescript
// Desarrollo local
const dsnDesarrollo = 'deviceID=dev-001&amqp_uri=amqp://localhost:5672&timeout=10s&debug=true';

// Producción con autenticación
const dsnProduccion = 'deviceID=prod-webapp-001&amqp_uri=amqp://usuario:contraseña@rabbitmq.empresa.com:5672&timeout=30s&debug=false';

// Con VHost personalizado
const dsnVhost = 'deviceID=staging-001&amqp_uri=amqp://usuario:contraseña@localhost:5672/staging&timeout=15s';
```

---

## 📖 Tipos de Ejecución

### 1. 🗃️ Consultas SQL (`sql`)

Ejecuta consultas SQL directas con binding de parámetros y soporte completo de transacciones.

```typescript
// Consulta simple
const usuarios = await client.query('SELECT id, nombre, email FROM usuarios');

// Consulta con parámetros
const usuariosActivos = await client.query(
    'SELECT * FROM usuarios WHERE activo = ? AND rol = ?',
    [true, 'admin']
);

// Consulta con JOINS
const pedidosConDetalles = await client.query(`
    SELECT p.id, p.fecha, u.nombre as cliente, COUNT(pd.id) as items
    FROM pedidos p
    JOIN usuarios u ON p.usuario_id = u.id
    LEFT JOIN pedidos_detalles pd ON p.id = pd.pedido_id
    WHERE p.fecha >= ?
    GROUP BY p.id, p.fecha, u.nombre
    ORDER BY p.fecha DESC
`, ['2024-01-01']);
```

**Características:**
- Binding de parámetros para seguridad
- Soporte para transacciones
- Pooling de conexiones
- Manejo seguro de tipos de resultado

### 2. ⚙️ Funciones Personalizadas (`function`)

Ejecuta funciones del lado del servidor con parámetros tipados y múltiples valores de retorno.

```typescript
// Función simple
const longitud = await client.query('FUNCTION:{"name":"lengthOfString","params":[{"type":"string","value":"Hola Mundo"}]}');

// Función con múltiples parámetros
const suma = await client.query('FUNCTION:{"name":"addIntegers","params":[{"type":"int","value":15},{"type":"int","value":25}]}');

// Función de sistema
const infoSistema = await client.query('FUNCTION:{"name":"getSystemInfo","params":[]}');

// Función de procesamiento de datos
const hash = await client.query('FUNCTION:{"name":"calculateHash","params":[{"type":"string","value":"mi-dato-importante"}]}');
```

**Funciones Integradas (16+):**
- `lengthOfString`: Obtener longitud de cadena
- `addIntegers`: Sumar dos enteros
- `getCurrentTimestamp`: Obtener timestamp actual
- `generateUUID`: Generar UUID
- `encodeBase64`: Codificación Base64
- `decodeBase64`: Decodificación Base64
- `parseJSON`: Parsear cadena JSON
- `formatJSON`: Formatear JSON con indentación
- `getSystemInfo`: Obtener información del sistema
- `listFiles`: Listar contenido de directorio
- `readFile`: Leer contenido de archivo
- `writeFile`: Escribir contenido de archivo
- `calculateHash`: Calcular hash SHA256
- `validateEmail`: Validar dirección de email
- `generateRandomString`: Generar cadena aleatoria
- `convertTimezone`: Convertir zona horaria

### 3. 🖥️ Comandos del Sistema (`command`)

Ejecuta comandos del sistema con acceso controlado y gestión de timeouts.

```typescript
// Información del sistema
const procesos = await client.query('COMMAND:ps aux | grep node');
const discos = await client.query('COMMAND:df -h');
const memoria = await client.query('COMMAND:free -m');

// Gestión de servicios
const estadoNginx = await client.query('COMMAND:systemctl status nginx');
const logsApp = await client.query('COMMAND:tail -n 50 /var/log/app.log');

// Operaciones de archivos
const archivos = await client.query('COMMAND:ls -la /home/usuario/documentos');
const configuracion = await client.query('COMMAND:cat /etc/config/app.conf');
```

**Características:**
- Captura de stdout y stderr
- Timeouts configurables
- Preservación de salida línea por línea
- Manejo de códigos de error

---

## 🔧 Referencia de la API

### Clase BurrowClient

#### Constructor
```typescript
const client = new BurrowClient();
```

#### Métodos

##### `connect(dsn: string): Promise<void>`
Se conecta al servidor BurrowCtl usando el DSN proporcionado.

```typescript
await client.connect('deviceID=mi-dispositivo&amqp_uri=amqp://localhost:5672&timeout=10s&debug=true');
```

##### `query(sql: string, params?: QueryValue[], options?: QueryOptions): Promise<Rows>`
Ejecuta una consulta SQL con parámetros opcionales.

```typescript
// Consulta simple
const rows = await client.query('SELECT * FROM usuarios');

// Consulta con parámetros
const rows = await client.query('SELECT * FROM usuarios WHERE edad > ?', [18]);

// Consulta con opciones
const rows = await client.query('SELECT * FROM datos_grandes', [], { timeout: 60000 });
```

**Parámetros:**
- `sql`: Cadena de consulta SQL
- `params`: Array de parámetros para consultas parametrizadas
- `options`: Opciones de consulta (timeout, etc.)

##### `close(): Promise<void>`
Cierra la conexión al servidor.

```typescript
await client.close();
```

##### `isConnected(): boolean`
Retorna si el cliente está actualmente conectado.

```typescript
if (client.isConnected()) {
    console.log('Cliente conectado');
} else {
    console.log('Cliente desconectado');
}
```

### Clase Rows

La clase `Rows` proporciona métodos para iterar y acceder a los resultados de las consultas.

#### Métodos

##### `getRows(): any[]`
Retorna todas las filas como un array.

```typescript
const rows = await client.query('SELECT * FROM usuarios');
const filas = rows.getRows();
console.log('Total de filas:', filas.length);
```

##### `getColumns(): string[]`
Retorna los nombres de las columnas.

```typescript
const rows = await client.query('SELECT id, nombre, email FROM usuarios');
const columnas = rows.getColumns();
console.log('Columnas:', columnas); // ['id', 'nombre', 'email']
```

##### `length(): number`
Retorna el número de filas.

```typescript
const rows = await client.query('SELECT * FROM productos');
console.log(`Encontrados ${rows.length()} productos`);
```

##### `hasNext(): boolean`
Verifica si hay más filas para iterar.

```typescript
const rows = await client.query('SELECT * FROM usuarios');
while (rows.hasNext()) {
    const fila = rows.next();
    console.log('Usuario:', fila.nombre);
}
```

##### `next(): any`
Retorna la siguiente fila y avanza el iterador.

```typescript
const rows = await client.query('SELECT * FROM usuarios');
const primeraFila = rows.next();
console.log('Primera fila:', primeraFila);
```

##### `reset(): void`
Reinicia el iterador al principio.

```typescript
const rows = await client.query('SELECT * FROM usuarios');
rows.next(); // Avanza al primer elemento
rows.reset(); // Vuelve al principio
```

### Tipos TypeScript

#### QueryValue
```typescript
type QueryValue = string | number | boolean | null | undefined;
```

#### QueryOptions
```typescript
interface QueryOptions {
    timeout?: number; // Timeout en milisegundos
}
```

#### FunctionParam
```typescript
interface FunctionParam {
    type: 'string' | 'int' | 'float64' | 'bool' | 'array' | 'object';
    value: any;
}
```

#### FunctionRequest
```typescript
interface FunctionRequest {
    name: string;
    params: FunctionParam[];
}
```

---

## 📋 Ejemplos Completos

### 1. Aplicación de Gestión de Usuarios

```typescript
import { createClient } from 'burrowctl-client-nodejs';

class GestorUsuarios {
    private client: any;

    async conectar() {
        this.client = await createClient(
            'deviceID=sistema-usuarios&amqp_uri=amqp://localhost:5672&timeout=30s&debug=true'
        );
    }

    async obtenerUsuarios(filtros: {activo?: boolean, rol?: string} = {}) {
        let sql = 'SELECT id, nombre, email, activo, rol FROM usuarios WHERE 1=1';
        const params: any[] = [];

        if (filtros.activo !== undefined) {
            sql += ' AND activo = ?';
            params.push(filtros.activo);
        }

        if (filtros.rol) {
            sql += ' AND rol = ?';
            params.push(filtros.rol);
        }

        const rows = await this.client.query(sql, params);
        return rows.getRows();
    }

    async crearUsuario(datos: {nombre: string, email: string, rol: string}) {
        const resultado = await this.client.query(
            'INSERT INTO usuarios (nombre, email, rol, activo) VALUES (?, ?, ?, ?)',
            [datos.nombre, datos.email, datos.rol, true]
        );
        return resultado;
    }

    async actualizarUsuario(id: number, datos: any) {
        const campos = Object.keys(datos);
        const valores = Object.values(datos);
        const sql = `UPDATE usuarios SET ${campos.map(c => `${c} = ?`).join(', ')} WHERE id = ?`;
        
        const resultado = await this.client.query(sql, [...valores, id]);
        return resultado;
    }

    async cerrar() {
        await this.client.close();
    }
}

// Uso
const gestor = new GestorUsuarios();
await gestor.conectar();

// Obtener usuarios activos
const usuariosActivos = await gestor.obtenerUsuarios({ activo: true });
console.log('Usuarios activos:', usuariosActivos);

// Crear nuevo usuario
await gestor.crearUsuario({
    nombre: 'Juan Pérez',
    email: 'juan@ejemplo.com',
    rol: 'usuario'
});

await gestor.cerrar();
```

### 2. Monitor de Sistema

```typescript
import { createClient } from 'burrowctl-client-nodejs';

class MonitorSistema {
    private client: any;

    async conectar() {
        this.client = await createClient(
            'deviceID=monitor-sistema&amqp_uri=amqp://localhost:5672&timeout=60s'
        );
    }

    async obtenerInfoSistema() {
        const infoSistema = await this.client.query('FUNCTION:{"name":"getSystemInfo","params":[]}');
        return JSON.parse(infoSistema.getRows()[0].result);
    }

    async verificarProcesos() {
        const procesos = await this.client.query('COMMAND:ps aux | grep -E "(nginx|mysql|node)"');
        return procesos.getRows();
    }

    async obtenerUsoDiscos() {
        const discos = await this.client.query('COMMAND:df -h');
        return discos.getRows();
    }

    async verificarServicios() {
        const servicios = ['nginx', 'mysql', 'rabbitmq-server'];
        const estados = [];

        for (const servicio of servicios) {
            try {
                const estado = await this.client.query(`COMMAND:systemctl is-active ${servicio}`);
                estados.push({
                    servicio,
                    estado: estado.getRows()[0].output.trim(),
                    activo: estado.getRows()[0].output.trim() === 'active'
                });
            } catch (error) {
                estados.push({
                    servicio,
                    estado: 'error',
                    activo: false,
                    error: error.message
                });
            }
        }

        return estados;
    }

    async generarReporte() {
        const info = await this.obtenerInfoSistema();
        const procesos = await this.verificarProcesos();
        const discos = await this.obtenerUsoDiscos();
        const servicios = await this.verificarServicios();

        return {
            timestamp: new Date().toISOString(),
            sistema: info,
            procesos: procesos.length,
            discos: discos,
            servicios: servicios
        };
    }

    async cerrar() {
        await this.client.close();
    }
}

// Uso
const monitor = new MonitorSistema();
await monitor.conectar();

const reporte = await monitor.generarReporte();
console.log('Reporte del sistema:', JSON.stringify(reporte, null, 2));

await monitor.cerrar();
```

### 3. Procesador de Datos

```typescript
import { createClient } from 'burrowctl-client-nodejs';

class ProcesadorDatos {
    private client: any;

    async conectar() {
        this.client = await createClient(
            'deviceID=procesador-datos&amqp_uri=amqp://localhost:5672&timeout=120s'
        );
    }

    async procesarArchivo(rutaArchivo: string) {
        // Leer archivo
        const contenido = await this.client.query(`FUNCTION:{"name":"readFile","params":[{"type":"string","value":"${rutaArchivo}"}]}`);
        
        // Procesar contenido
        const datos = JSON.parse(contenido.getRows()[0].result);
        
        // Validar datos
        for (const item of datos) {
            if (item.email) {
                const validacion = await this.client.query(`FUNCTION:{"name":"validateEmail","params":[{"type":"string","value":"${item.email}"}]}`);
                item.email_valido = validacion.getRows()[0].result === 'true';
            }
        }

        // Generar hash de los datos
        const hash = await this.client.query(`FUNCTION:{"name":"calculateHash","params":[{"type":"string","value":"${JSON.stringify(datos)}"}]}`);
        
        // Guardar en base de datos
        for (const item of datos) {
            await this.client.query(
                'INSERT INTO datos_procesados (nombre, email, email_valido, hash_proceso) VALUES (?, ?, ?, ?)',
                [item.nombre, item.email, item.email_valido, hash.getRows()[0].result]
            );
        }

        return {
            procesados: datos.length,
            hash: hash.getRows()[0].result
        };
    }

    async obtenerEstadisticas() {
        const total = await this.client.query('SELECT COUNT(*) as total FROM datos_procesados');
        const emailsValidos = await this.client.query('SELECT COUNT(*) as validos FROM datos_procesados WHERE email_valido = true');
        const porHora = await this.client.query(`
            SELECT DATE_FORMAT(created_at, '%Y-%m-%d %H:00:00') as hora, COUNT(*) as cantidad
            FROM datos_procesados 
            WHERE created_at >= DATE_SUB(NOW(), INTERVAL 24 HOUR)
            GROUP BY DATE_FORMAT(created_at, '%Y-%m-%d %H:00:00')
            ORDER BY hora DESC
        `);

        return {
            total: total.getRows()[0].total,
            emailsValidos: emailsValidos.getRows()[0].validos,
            porHora: porHora.getRows()
        };
    }

    async cerrar() {
        await this.client.close();
    }
}

// Uso
const procesador = new ProcesadorDatos();
await procesador.conectar();

const resultado = await procesador.procesarArchivo('/path/to/datos.json');
console.log('Procesamiento completado:', resultado);

const estadisticas = await procesador.obtenerEstadisticas();
console.log('Estadísticas:', estadisticas);

await procesador.cerrar();
```

---

## 🔍 Manejo de Errores

### Errores Comunes

```typescript
import { createClient } from 'burrowctl-client-nodejs';

try {
    const client = await createClient('deviceID=test&amqp_uri=amqp://localhost:5672');
    
    const rows = await client.query('SELECT * FROM usuarios');
    console.log(rows.getRows());
    
} catch (error) {
    // Manejo específico de errores
    if (error.message.includes('timeout')) {
        console.error('Error: La consulta excedió el tiempo límite');
    } else if (error.message.includes('connection refused')) {
        console.error('Error: No se puede conectar a RabbitMQ');
    } else if (error.message.includes('authentication failed')) {
        console.error('Error: Credenciales de RabbitMQ inválidas');
    } else if (error.message.includes('Unknown database')) {
        console.error('Error: Base de datos no encontrada');
    } else if (error.message.includes('Access denied')) {
        console.error('Error: Acceso denegado a la base de datos');
    } else {
        console.error('Error inesperado:', error.message);
    }
} finally {
    await client.close();
}
```

### Manejo Avanzado con Reintentos

```typescript
class ClienteRobusto {
    private client: any;
    private maxReintentos = 3;
    private tiempoEspera = 1000;

    async conectarConReintentos(dsn: string) {
        for (let intento = 1; intento <= this.maxReintentos; intento++) {
            try {
                this.client = await createClient(dsn);
                console.log('Conexión exitosa');
                return;
            } catch (error) {
                console.log(`Intento ${intento} fallido:`, error.message);
                
                if (intento < this.maxReintentos) {
                    await new Promise(resolve => setTimeout(resolve, this.tiempoEspera * intento));
                } else {
                    throw new Error(`Falló la conexión después de ${this.maxReintentos} intentos`);
                }
            }
        }
    }

    async queryConReintentos(sql: string, params: any[] = []) {
        for (let intento = 1; intento <= this.maxReintentos; intento++) {
            try {
                return await this.client.query(sql, params);
            } catch (error) {
                console.log(`Query intento ${intento} fallido:`, error.message);
                
                if (intento < this.maxReintentos) {
                    await new Promise(resolve => setTimeout(resolve, this.tiempoEspera * intento));
                } else {
                    throw error;
                }
            }
        }
    }
}
```

---

## 🚀 Optimización de Rendimiento

### Reutilización de Conexiones

```typescript
// ❌ Ineficiente - crear conexión para cada consulta
async function consultaIneficiente() {
    const client = await createClient('deviceID=test&amqp_uri=amqp://localhost:5672');
    const rows = await client.query('SELECT * FROM usuarios');
    await client.close();
    return rows;
}

// ✅ Eficiente - reutilizar conexión
class GestorConexion {
    private client: any;
    private conectado = false;

    async asegurarConexion() {
        if (!this.conectado) {
            this.client = await createClient('deviceID=test&amqp_uri=amqp://localhost:5672');
            this.conectado = true;
        }
    }

    async query(sql: string, params: any[] = []) {
        await this.asegurarConexion();
        return this.client.query(sql, params);
    }

    async cerrar() {
        if (this.conectado) {
            await this.client.close();
            this.conectado = false;
        }
    }
}
```

### Consultas en Lote

```typescript
// Procesamiento en lotes para mejor rendimiento
async function procesarDatosEnLotes(datos: any[], tamañoLote = 100) {
    const client = await createClient('deviceID=batch&amqp_uri=amqp://localhost:5672');
    
    for (let i = 0; i < datos.length; i += tamañoLote) {
        const lote = datos.slice(i, i + tamañoLote);
        
        // Construir consulta de inserción múltiple
        const valores = lote.map(() => '(?, ?, ?)').join(', ');
        const sql = `INSERT INTO datos (nombre, email, fecha) VALUES ${valores}`;
        
        // Aplanar parámetros
        const params = lote.flatMap(item => [item.nombre, item.email, item.fecha]);
        
        await client.query(sql, params);
        console.log(`Procesado lote ${Math.floor(i / tamañoLote) + 1}`);
    }
    
    await client.close();
}
```

---

## 🔒 Seguridad y Mejores Prácticas

### Configuración Segura

```typescript
// ✅ Usar variables de entorno
const dsn = `deviceID=${process.env.DEVICE_ID}&amqp_uri=${process.env.AMQP_URI}&timeout=30s&debug=false`;

// ✅ Validar parámetros de entrada
function validarEntrada(valor: any): boolean {
    if (typeof valor === 'string') {
        // Evitar caracteres peligrosos
        return !/[<>\"'%;)(&+]/.test(valor);
    }
    return true;
}

// ✅ Usar consultas parametrizadas
async function consultaSegura(id: number, nombre: string) {
    if (!validarEntrada(nombre)) {
        throw new Error('Parámetro inválido');
    }
    
    const client = await createClient(dsn);
    const rows = await client.query(
        'SELECT * FROM usuarios WHERE id = ? AND nombre = ?',
        [id, nombre]
    );
    await client.close();
    return rows;
}
```

### Logging y Monitoreo

```typescript
class ClienteConLogging {
    private client: any;
    private logger: any;

    constructor(logger: any) {
        this.logger = logger;
    }

    async conectar(dsn: string) {
        try {
            this.client = await createClient(dsn);
            this.logger.info('Conexión a burrowctl establecida');
        } catch (error) {
            this.logger.error('Error al conectar:', error);
            throw error;
        }
    }

    async query(sql: string, params: any[] = []) {
        const inicio = Date.now();
        try {
            const resultado = await this.client.query(sql, params);
            const duracion = Date.now() - inicio;
            
            this.logger.info('Query ejecutado', {
                sql: sql.substring(0, 100) + '...',
                parametros: params.length,
                filas: resultado.length(),
                duracion: duracion + 'ms'
            });
            
            return resultado;
        } catch (error) {
            this.logger.error('Error en query:', {
                sql: sql.substring(0, 100) + '...',
                error: error.message
            });
            throw error;
        }
    }
}
```

---

## 🛠️ Desarrollo

### Configuración del Entorno de Desarrollo

```bash
# Clonar el repositorio
git clone https://github.com/lordbasex/burrowctl.git
cd burrowctl/client-nodejs

# Instalar dependencias
npm install

# Compilar TypeScript
npm run build

# Modo desarrollo con watch
npm run dev

# Ejecutar tests
npm test

# Linting
npm run lint
```

### Scripts Disponibles

```json
{
  "scripts": {
    "build": "tsc",
    "dev": "tsc --watch",
    "test": "jest",
    "lint": "eslint src/**/*.ts",
    "lint:fix": "eslint src/**/*.ts --fix",
    "clean": "rm -rf dist",
    "prepublishOnly": "npm run build"
  }
}
```

---

## 📊 Requisitos del Sistema

### Dependencias de Tiempo de Ejecución

- **Node.js**: 16.0.0 o superior
- **RabbitMQ**: 3.8.0 o superior
- **Servidor burrowctl**: Ejecutándose y accesible

### Dependencias NPM

```json
{
  "dependencies": {
    "amqplib": "^0.10.3",
    "uuid": "^9.0.0"
  },
  "devDependencies": {
    "@types/node": "^20.0.0",
    "@types/amqplib": "^0.10.1",
    "@types/uuid": "^9.0.0",
    "typescript": "^5.0.0",
    "jest": "^29.0.0",
    "@types/jest": "^29.0.0"
  }
}
```

---

## 🤝 Contribuir

¡Las contribuciones son bienvenidas! Por favor consulta nuestra [Guía de Contribución](../CONTRIBUTING.md) para más detalles.

### Configuración para Desarrollo

1. Haz fork del repositorio
2. Crea una rama para tu funcionalidad: `git checkout -b feature/funcionalidad-increible`
3. Realiza tus cambios
4. Añade tests para la nueva funcionalidad
5. Ejecuta los tests: `npm test`
6. Confirma tus cambios: `git commit -m 'Agregar funcionalidad increíble'`
7. Push a la rama: `git push origin feature/funcionalidad-increible`
8. Abre un Pull Request

---

## 📜 Licencia

Este proyecto está licenciado bajo la Licencia MIT - consulta el archivo [LICENSE](../LICENSE) para más detalles.

---

## 🆘 Soporte

- **Documentación**: [Documentación completa](../examples/)
- **Ejemplos**: [Ejemplos de uso](../examples/client/)
- **Issues**: [GitHub Issues](https://github.com/lordbasex/burrowctl/issues)
- **Discusiones**: [GitHub Discussions](https://github.com/lordbasex/burrowctl/discussions)

---

## 🙏 Agradecimientos

- [RabbitMQ](https://www.rabbitmq.com/) por el excelente message broker
- [amqplib](https://github.com/amqp-node/amqplib) por la biblioteca de cliente RabbitMQ
- [TypeScript](https://www.typescriptlang.org/) por el sistema de tipos
- La comunidad Node.js por su excelente ecosistema

---

<div align="center">
  <p>Hecho con ❤️ por el equipo de burrowctl</p>
  <p>
    <a href="https://github.com/lordbasex/burrowctl/stargazers">⭐ Marca este proyecto</a> | 
    <a href="https://github.com/lordbasex/burrowctl/issues">🐛 Reportar Bug</a> | 
    <a href="https://github.com/lordbasex/burrowctl/issues">💡 Solicitar Funcionalidad</a>
  </p>
</div>