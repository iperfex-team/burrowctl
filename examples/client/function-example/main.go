package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

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

	// Verificar argumentos de línea de comandos
	if len(os.Args) > 1 {
		functionName := os.Args[1]
		args := os.Args[2:] // Resto de argumentos son parámetros

		if functionName == "help" || functionName == "--help" || functionName == "-h" {
			showUsage()
			return
		}

		executeFunctionWithArgs(db, functionName, args)
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
	showUsage()
}

func showUsage() {
	fmt.Println("\n💡 Uso del cliente:")
	fmt.Println("================")
	fmt.Println()
	fmt.Println("🔹 Funciones sin parámetros:")
	fmt.Println("   go run main.go returnString")
	fmt.Println("   go run main.go returnInt")
	fmt.Println("   go run main.go returnBool")
	fmt.Println("   go run main.go returnStruct")
	fmt.Println("   go run main.go returnError")
	fmt.Println()
	fmt.Println("🔹 Funciones con parámetros:")
	fmt.Println("   go run main.go lengthOfString \"Hello World\"")
	fmt.Println("   go run main.go isEven 42")
	fmt.Println("   go run main.go sumArray \"[1,2,3,4,5]\"")
	fmt.Println("   go run main.go greetPerson \"Juan\" 30")
	fmt.Println("   go run main.go validateString \"test\"")
	fmt.Println("   go run main.go modifyJSON '{\"name\":\"Juan\",\"age\":30}'")
	fmt.Println()
	fmt.Println("🔹 Funciones con múltiples valores de retorno:")
	fmt.Println("   go run main.go complexFunction \"Go\" 10")
	fmt.Println()
	fmt.Println("🔹 Ejemplos especiales:")
	fmt.Println("   go run main.go flagToPerson true")
	fmt.Println("   go run main.go validateString \"\"  # Prueba string vacío")
	fmt.Println()
	fmt.Println("📋 Funciones disponibles:")
	fmt.Println("   Sin parámetros: returnError, returnBool, returnInt, returnString,")
	fmt.Println("                   returnStruct, returnIntArray, returnStringArray, returnJSON")
	fmt.Println("   Con parámetros: lengthOfString, isEven, sumArray, greetPerson,")
	fmt.Println("                   validateString, flagToPerson, modifyJSON, complexFunction")
}

func executeFunctionWithArgs(db *sql.DB, functionName string, args []string) {
	fmt.Printf("\n🔧 Executing function: %s", functionName)
	if len(args) > 0 {
		fmt.Printf(" with args: %v", args)
	}
	fmt.Println()

	// Construir parámetros según la función y los argumentos
	params, err := buildFunctionParams(functionName, args)
	if err != nil {
		log.Printf("❌ Error building parameters: %v", err)
		showFunctionHelp(functionName)
		return
	}

	executeFunction(db, functionName, params)
}

func buildFunctionParams(functionName string, args []string) ([]FunctionParam, error) {
	switch functionName {
	// Funciones sin parámetros
	case "returnError", "returnBool", "returnInt", "returnString",
		"returnStruct", "returnIntArray", "returnStringArray", "returnJSON":
		if len(args) > 0 {
			return nil, fmt.Errorf("function '%s' no acepta parámetros", functionName)
		}
		return []FunctionParam{}, nil

	// lengthOfString(string) int
	case "lengthOfString":
		if len(args) != 1 {
			return nil, fmt.Errorf("lengthOfString requiere 1 parámetro: string")
		}
		return []FunctionParam{
			{Type: "string", Value: args[0]},
		}, nil

	// isEven(int) bool
	case "isEven":
		if len(args) != 1 {
			return nil, fmt.Errorf("isEven requiere 1 parámetro: int")
		}
		num, err := strconv.Atoi(args[0])
		if err != nil {
			return nil, fmt.Errorf("isEven requiere un número entero, recibido: %s", args[0])
		}
		return []FunctionParam{
			{Type: "int", Value: num},
		}, nil

	// sumArray([]int) int
	case "sumArray":
		if len(args) != 1 {
			return nil, fmt.Errorf("sumArray requiere 1 parámetro: array de enteros como '[1,2,3]'")
		}
		// Parsear array JSON
		var nums []int
		if err := json.Unmarshal([]byte(args[0]), &nums); err != nil {
			return nil, fmt.Errorf("sumArray requiere array JSON válido, ej: '[1,2,3,4,5]'")
		}
		// Convertir a []interface{}
		intArray := make([]interface{}, len(nums))
		for i, v := range nums {
			intArray[i] = v
		}
		return []FunctionParam{
			{Type: "[]int", Value: intArray},
		}, nil

	// greetPerson(Person) string
	case "greetPerson":
		if len(args) != 2 {
			return nil, fmt.Errorf("greetPerson requiere 2 parámetros: nombre edad")
		}
		age, err := strconv.Atoi(args[1])
		if err != nil {
			return nil, fmt.Errorf("greetPerson: edad debe ser un número entero")
		}
		return []FunctionParam{
			{Type: "Person", Value: map[string]interface{}{
				"name": args[0],
				"age":  age,
			}},
		}, nil

	// validateString(string) error
	case "validateString":
		if len(args) != 1 {
			return nil, fmt.Errorf("validateString requiere 1 parámetro: string")
		}
		return []FunctionParam{
			{Type: "string", Value: args[0]},
		}, nil

	// flagToPerson(bool) Person
	case "flagToPerson":
		if len(args) != 1 {
			return nil, fmt.Errorf("flagToPerson requiere 1 parámetro: bool (true/false)")
		}
		flag, err := strconv.ParseBool(args[0])
		if err != nil {
			return nil, fmt.Errorf("flagToPerson requiere bool válido (true/false)")
		}
		return []FunctionParam{
			{Type: "bool", Value: flag},
		}, nil

	// modifyJSON(string) (string, error)
	case "modifyJSON":
		if len(args) != 1 {
			return nil, fmt.Errorf("modifyJSON requiere 1 parámetro: JSON string")
		}
		// Validar que sea JSON válido
		var temp interface{}
		if err := json.Unmarshal([]byte(args[0]), &temp); err != nil {
			return nil, fmt.Errorf("modifyJSON requiere JSON válido")
		}
		return []FunctionParam{
			{Type: "string", Value: args[0]},
		}, nil

	// complexFunction(string, int) (string, int, error)
	case "complexFunction":
		if len(args) != 2 {
			return nil, fmt.Errorf("complexFunction requiere 2 parámetros: string int")
		}
		num, err := strconv.Atoi(args[1])
		if err != nil {
			return nil, fmt.Errorf("complexFunction: segundo parámetro debe ser entero")
		}
		return []FunctionParam{
			{Type: "string", Value: args[0]},
			{Type: "int", Value: num},
		}, nil

	default:
		return nil, fmt.Errorf("función desconocida: %s", functionName)
	}
}

func showFunctionHelp(functionName string) {
	fmt.Printf("\n💡 Ayuda para función '%s':\n", functionName)

	switch functionName {
	case "lengthOfString":
		fmt.Println("   Uso: go run main.go lengthOfString \"texto\"")
		fmt.Println("   Ejemplo: go run main.go lengthOfString \"Hello World\"")
	case "isEven":
		fmt.Println("   Uso: go run main.go isEven número")
		fmt.Println("   Ejemplo: go run main.go isEven 42")
	case "sumArray":
		fmt.Println("   Uso: go run main.go sumArray \"[num1,num2,num3]\"")
		fmt.Println("   Ejemplo: go run main.go sumArray \"[1,2,3,4,5]\"")
	case "greetPerson":
		fmt.Println("   Uso: go run main.go greetPerson \"nombre\" edad")
		fmt.Println("   Ejemplo: go run main.go greetPerson \"Juan\" 30")
	case "validateString":
		fmt.Println("   Uso: go run main.go validateString \"texto\"")
		fmt.Println("   Ejemplo: go run main.go validateString \"test\"")
	case "flagToPerson":
		fmt.Println("   Uso: go run main.go flagToPerson true/false")
		fmt.Println("   Ejemplo: go run main.go flagToPerson true")
	case "modifyJSON":
		fmt.Println("   Uso: go run main.go modifyJSON '{\"name\":\"Juan\",\"age\":30}'")
		fmt.Println("   Ejemplo: go run main.go modifyJSON '{\"name\":\"María\",\"age\":25}'")
	case "complexFunction":
		fmt.Println("   Uso: go run main.go complexFunction \"texto\" número")
		fmt.Println("   Ejemplo: go run main.go complexFunction \"Go\" 10")
	default:
		fmt.Printf("   Función '%s' sin parámetros\n", functionName)
		fmt.Printf("   Uso: go run main.go %s\n", functionName)
	}
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
