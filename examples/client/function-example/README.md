# Function Execution Example

Este ejemplo demuestra cómo ejecutar funciones remotas usando burrowctl.

## Características

- **Ejecución de funciones remotas** con reflection
- **Soporte para múltiples tipos de parámetros**: string, int, bool, arrays, structs
- **Manejo de múltiples valores de retorno**
- **Conversión automática de tipos**
- **Manejo de errores**

## Funciones Disponibles

### Funciones Sin Parámetros
- `returnBool()` - Devuelve un valor booleano
- `returnInt()` - Devuelve un entero
- `returnString()` - Devuelve un string
- `returnStruct()` - Devuelve un struct Person
- `returnIntArray()` - Devuelve un array de enteros
- `returnStringArray()` - Devuelve un array de strings
- `returnJSON()` - Devuelve JSON como string
- `returnError()` - Devuelve un error

### Funciones Con Parámetros
- `lengthOfString(string)` - Devuelve la longitud de un string
- `isEven(int)` - Verifica si un número es par
- `sumArray([]int)` - Suma todos los elementos de un array
- `greetPerson(Person)` - Saluda a una persona
- `validateString(string)` - Valida que un string no esté vacío
- `flagToPerson(bool)` - Convierte un booleano a Person
- `modifyJSON(string)` - Modifica un JSON y devuelve el resultado

### Funciones Con Múltiples Valores de Retorno
- `complexFunction(string, int)` - Devuelve (string, int, error)

## Formato de Solicitud

Las funciones se ejecutan enviando una solicitud JSON con el formato:

```json
{
  "name": "functionName",
  "params": [
    {
      "type": "string",
      "value": "example"
    },
    {
      "type": "int", 
      "value": 42
    }
  ]
}
```

## Tipos de Parámetros Soportados

- `string` - Cadenas de texto
- `int` - Números enteros
- `bool` - Valores booleanos
- `[]int` - Arrays de enteros
- `[]string` - Arrays de strings
- `Person` - Struct con campos name y age

## Ejemplos de Uso

### Ejecutar Todas las Funciones de Ejemplo
```bash
go run main.go
```

### Ejecutar una Función Específica
```bash
go run main.go returnString
go run main.go lengthOfString
go run main.go isEven
go run main.go sumArray
```

### Usando el Makefile
```bash
make run-function-example
```

## Ejemplos de Solicitudes

### Función Sin Parámetros
```json
{
  "name": "returnString",
  "params": []
}
```

### Función Con String
```json
{
  "name": "lengthOfString",
  "params": [
    {
      "type": "string",
      "value": "Hello World"
    }
  ]
}
```

### Función Con Array
```json
{
  "name": "sumArray",
  "params": [
    {
      "type": "[]int",
      "value": [1, 2, 3, 4, 5]
    }
  ]
}
```

### Función Con Struct
```json
{
  "name": "greetPerson",
  "params": [
    {
      "type": "Person",
      "value": {
        "name": "Juan",
        "age": 30
      }
    }
  ]
}
```

## Salida de Ejemplo

```
🔧 --- Function Examples ---

1. Functions without parameters:
  📋 Function: returnBool
     result              
     --------------------
     true                

  📋 Function: returnInt
     result              
     --------------------
     42                  

  📋 Function: returnString
     result              
     --------------------
     Hola mundo          

2. Functions with parameters:
  📋 Function: lengthOfString (with 1 params)
     result              
     --------------------
     11                  

  📋 Function: isEven (with 1 params)
     result              
     --------------------
     true                

3. Functions that return errors:
  📋 Function: returnError
     error               
     --------------------
     algo salió mal      

4. Functions with multiple return values:
  📋 Function: complexFunction (with 2 params)
     result_1            | result_2            | result_3            
     -------------------- | -------------------- | --------------------
     Go                  | 20                  | success             
```

## Notas Técnicas

- El servidor usa **reflection** para ejecutar funciones dinámicamente
- Los parámetros se convierten automáticamente a los tipos esperados
- Se soportan **múltiples valores de retorno**
- Los errores se manejan automáticamente y se devuelven como columnas separadas
- El timeout por defecto es de **30 segundos**

## Estructura del Código

- `FunctionRequest` - Estructura para la solicitud de función
- `FunctionParam` - Estructura para parámetros con tipo
- `executeFunction()` - Función principal para ejecutar funciones remotas
- `executeSingleFunction()` - Función para ejecutar una función específica

## Extensión

Para agregar nuevas funciones:

1. Agregar la función al mapa `functions` en `server/server.go`
2. Asegurarse de que los tipos de parámetros estén soportados
3. Agregar ejemplos de uso en este archivo

¡La funcionalidad de función está completamente implementada y lista para usar! 