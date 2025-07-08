# 🚀 Advanced Client Example

Este ejemplo demuestra las nuevas características empresariales del cliente burrowctl.

## 🎯 Características Demostradas

- **🔄 Reconexión Automática**: Manejo inteligente de desconexiones
- **📝 Prepared Statements**: Mejor rendimiento y seguridad
- **🛡️ Rate Limiting**: Protección contra abuso (demo del lado servidor)
- **⚡ Configuración Avanzada**: Timeouts, debug y parámetros personalizables

## 🏁 Inicio Rápido

```bash
# Compilar el ejemplo
go build advanced-main.go

# Ejecutar con configuración básica
./advanced-main

# Ver todas las opciones
./advanced-main -help
```

## 📋 Modos de Demostración

### 1. Demo Básico (por defecto)
```bash
./advanced-main
```
Ejecuta una query simple mostrando las configuraciones aplicadas.

### 2. Prepared Statements
```bash
./advanced-main -prepared
```
Demuestra el uso de prepared statements con diferentes parámetros.

### 3. Test de Estrés (Rate Limiting)
```bash
./advanced-main -stress -concurrent=10 -requests=100
```
Genera carga para demostrar el rate limiting del servidor.

### 4. Demo de Reconexión
```bash
./advanced-main -reconnect-demo
```
Muestra la reconexión automática en acción.

## ⚙️ Configuración Avanzada

### Opciones de Conexión
```bash
./advanced-main \
  -device=mydevice \
  -amqp=amqp://user:pass@host:5672/ \
  -timeout=30s \
  -debug=true
```

### Test de Rendimiento
```bash
# Test de alta concurrencia
./advanced-main -stress -concurrent=50 -requests=200

# Test conservador  
./advanced-main -stress -concurrent=5 -requests=20
```

## 📊 Ejemplos de Salida

### Demo Básico
```
🗃️  Advanced burrowctl SQL Example
================================================
📱 Device ID: fd1825ec5a7b63f3fa2be9e04154d3b16f676663ba38e23d4ffafa7b0df29efb
🐰 RabbitMQ: amqp://burrowuser:burrowpass123@localhost:5672/
⏱️  Timeout: 10s
🐛 Debug: true
📝 Prepared Statements: false

🎯 Running Basic Demo
---------------------
📊 Executing: SELECT 'Hello' as greeting, 'World' as target, NOW() as timestamp

greeting        | target          | timestamp      
----------------|-----------------|----------------
Hello           | World           | 2024-01-15...
⏱️  Query completed in: 45ms
✅ Basic demo completed!
```

### Prepared Statements Demo
```
🎯 Running Prepared Statements Demo
-----------------------------------
📝 Preparing statement: SELECT ? as message, ? as number, ? as flag

📋 Execution 1 with params: [Hello World 42 true]
message         | number          | flag           
----------------|-----------------|----------------
Hello World     | 42              | true           
⏱️  Execution 1 completed in: 35ms

📋 Execution 2 with params: [Prepared Statement 100 false]
message         | number          | flag           
----------------|-----------------|----------------
Prepared Statement | 100          | false          
⏱️  Execution 2 completed in: 15ms
```

### Stress Test Demo
```
🎯 Running Stress Test (Rate Limiting Demo)
-------------------------------------------
🏗️  Concurrent connections: 10
📊 Requests per connection: 50
📈 Total requests: 500

📊 Live Results:
Worker | Success | Errors | Rate Limited | Avg Time
-------|---------|--------|--------------|----------
   1   | 45      | 5      | 5            |    25ms
   2   | 42      | 8      | 8            |    30ms
   ...

📈 Final Statistics:
✅ Total Successful: 420
❌ Total Errors: 80
🛡️  Rate Limited: 65
⏱️  Average Response Time: 27ms

🎯 Rate Limiting Demonstration Successful!
   Server protected against 65 excessive requests
```

## 🔧 Troubleshooting

### Error: Connection Failed
```bash
# Verificar conectividad con timeout extendido
./advanced-main -timeout=60s -debug=true
```

### Rate Limiting Activado
```bash
# Reducir concurrencia
./advanced-main -stress -concurrent=2 -requests=10
```

### Problemas de Reconexión
```bash
# Demo con logging detallado
./advanced-main -reconnect-demo -debug=true
```

## 💡 Tips de Uso

1. **Producción**: Usar `debug=false` y timeouts apropiados
2. **Desarrollo**: Activar debug para troubleshooting
3. **Testing**: Usar stress test para validar límites
4. **Monitoreo**: Observar logs de reconexión automática

## 🚀 Próximos Pasos

- Ejecutar el ejemplo del servidor avanzado: `../server/advanced/`
- Revisar la documentación: `../../ADVANCED_FEATURES.md`
- Explorar configuraciones de producción
- Implementar métricas personalizadas