package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"

	_ "github.com/lordbasex/burrowctl/client"
)

// Estructuras para construir solicitudes de función
type FunctionParam struct {
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}

type FunctionRequest struct {
	Name   string          `json:"name"`
	Params []FunctionParam `json:"params"`
}

func main() {
	// DSN con credenciales hardcodeadas para RabbitMQ
	dsn := "deviceID=fd1825ec5a7b63f3fa2be9e04154d3b16f676663ba38e23d4ffafa7b0df29efb&amqp_uri=amqp://burrowuser:burrowpass123@localhost:5672/&timeout=5s&debug=true"

	// Abrir conexión usando el driver rabbitsql
	db, err := sql.Open("rabbitsql", dsn)
	if err != nil {
		log.Fatal("Error connecting:", err)
	}
	defer db.Close()

	fmt.Printf("📡 Connecting to device: %s\n", "fd1825ec5a7b63f3fa2be9e04154d3b16f676663ba38e23d4ffafa7b0df29efb")

	// Verificar argumentos de línea de comandos para función específica
	if len(os.Args) > 1 {
		functionName := os.Args[1]
		executeSingleFunction(db, functionName)
		return
	}

	// Ejecutar múltiples ejemplos de funciones
	fmt.Println("\n🔧 --- Function Examples ---")

	// Ejemplos de funciones sin parámetros
	fmt.Println("\n1. Functions without parameters:")
	testFunctions := []string{"returnBool", "returnInt", "returnString", "returnStruct"}
	for _, funcName := range testFunctions {
		executeFunction(db, funcName, []FunctionParam{})
	}

	// Ejemplos de funciones con parámetros
	fmt.Println("\n2. Functions with parameters:")

	// Función que recibe string y devuelve int
	executeFunction(db, "lengthOfString", []FunctionParam{
		{Type: "string", Value: "Hello World"},
	})

	// Función que recibe int y devuelve bool
	executeFunction(db, "isEven", []FunctionParam{
		{Type: "int", Value: 42},
	})

	// Función que recibe array y devuelve suma
	executeFunction(db, "sumArray", []FunctionParam{
		{Type: "[]int", Value: []interface{}{1, 2, 3, 4, 5}},
	})

	// Función que recibe struct y devuelve string
	executeFunction(db, "greetPerson", []FunctionParam{
		{Type: "Person", Value: map[string]interface{}{
			"name": "María",
			"age":  28,
		}},
	})

	// Ejemplos de funciones que devuelven errores
	fmt.Println("\n3. Functions that return errors:")
	executeFunction(db, "returnError", []FunctionParam{})
	executeFunction(db, "validateString", []FunctionParam{
		{Type: "string", Value: ""},
	})
	executeFunction(db, "validateString", []FunctionParam{
		{Type: "string", Value: "valid string"},
	})

	// Función con múltiples valores de retorno
	fmt.Println("\n4. Functions with multiple return values:")
	executeFunction(db, "complexFunction", []FunctionParam{
		{Type: "string", Value: "Go"},
		{Type: "int", Value: 10},
	})

	fmt.Println("\n✅ All function examples executed successfully")
	fmt.Println("\n💡 Tips:")
	fmt.Println("   - You can specify a function name as argument:")
	fmt.Println("     go run main.go returnString")
	fmt.Println("   - Available functions: returnBool, returnInt, returnString, returnStruct")
	fmt.Println("     lengthOfString, isEven, sumArray, greetPerson, validateString, complexFunction")
	fmt.Println("   - The client sends JSON with function name and typed parameters")
}

func executeSingleFunction(db *sql.DB, functionName string) {
	fmt.Printf("\n🔧 Executing single function: %s\n", functionName)

	// Ejemplos de parámetros según la función
	var params []FunctionParam
	switch functionName {
	case "lengthOfString":
		params = []FunctionParam{{Type: "string", Value: "Hello World"}}
	case "isEven":
		params = []FunctionParam{{Type: "int", Value: 42}}
	case "sumArray":
		params = []FunctionParam{{Type: "[]int", Value: []interface{}{1, 2, 3, 4, 5}}}
	case "greetPerson":
		params = []FunctionParam{{Type: "Person", Value: map[string]interface{}{
			"name": "Juan",
			"age":  30,
		}}}
	case "validateString":
		params = []FunctionParam{{Type: "string", Value: "test string"}}
	case "complexFunction":
		params = []FunctionParam{
			{Type: "string", Value: "Go"},
			{Type: "int", Value: 5},
		}
	}

	executeFunction(db, functionName, params)
}

func executeFunction(db *sql.DB, functionName string, params []FunctionParam) {
	// Construir solicitud de función
	funcReq := FunctionRequest{
		Name:   functionName,
		Params: params,
	}

	// Serializar a JSON
	jsonData, err := json.Marshal(funcReq)
	if err != nil {
		log.Printf("Error marshaling function request: %v", err)
		return
	}

	// Ejecutar función remota usando el prefijo "FUNCTION:"
	query := "FUNCTION:" + string(jsonData)
	rows, err := db.Query(query)
	if err != nil {
		log.Printf("Error executing function %s: %v", functionName, err)
		return
	}
	defer rows.Close()

	fmt.Printf("  📋 Function: %s", functionName)
	if len(params) > 0 {
		fmt.Printf(" (with %d params)", len(params))
	}
	fmt.Println()

	// Obtener columnas
	cols, err := rows.Columns()
	if err != nil {
		log.Printf("Error getting columns: %v", err)
		return
	}

	// Mostrar encabezados
	fmt.Printf("     ")
	for i, col := range cols {
		if i > 0 {
			fmt.Printf(" | ")
		}
		fmt.Printf("%-20s", col)
	}
	fmt.Println()

	fmt.Printf("     ")
	for i := range cols {
		if i > 0 {
			fmt.Printf(" | ")
		}
		fmt.Printf("%-20s", "--------------------")
	}
	fmt.Println()

	// Procesar resultados
	for rows.Next() {
		// Crear slice de interface{} para escanear
		values := make([]interface{}, len(cols))
		scanArgs := make([]interface{}, len(cols))
		for i := range values {
			scanArgs[i] = &values[i]
		}

		if err := rows.Scan(scanArgs...); err != nil {
			log.Printf("Error scanning result: %v", err)
			return
		}

		// Mostrar resultados
		fmt.Printf("     ")
		for i, val := range values {
			if i > 0 {
				fmt.Printf(" | ")
			}
			fmt.Printf("%-20v", val)
		}
		fmt.Println()
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error iterating results: %v", err)
		return
	}

	fmt.Println()
}
