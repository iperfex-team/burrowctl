package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/lordbasex/burrowctl/server"
)

// Person representa el struct utilizado en varios ejemplos de función
type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

// ============================================================================
// FUNCIONES DE EJEMPLO DISPONIBLES PARA EJECUCIÓN REMOTA
// ============================================================================

// 1. Devuelve error
func returnError() error {
	return errors.New("algo salió mal")
}

// 2. Devuelve bool
func returnBool() bool {
	return true
}

// 3. Devuelve int
func returnInt() int {
	return 42
}

// 4. Devuelve string
func returnString() string {
	return "Hola mundo"
}

// 5. Devuelve struct
func returnStruct() Person {
	return Person{Name: "Juan", Age: 30}
}

// 6. Devuelve array de int
func returnIntArray() []int {
	return []int{1, 2, 3, 4, 5}
}

// 7. Devuelve array de string
func returnStringArray() []string {
	return []string{"uno", "dos", "tres"}
}

// 8. Devuelve JSON string
func returnJSON() string {
	p := Person{Name: "Ana", Age: 25}
	data, _ := json.Marshal(p)
	return string(data)
}

// 9. Entrada string y devuelve int
func lengthOfString(s string) int {
	return len(s)
}

// 10. Entrada int y devuelve bool
func isEven(n int) bool {
	return n%2 == 0
}

// 11. Entrada struct y devuelve string
func greetPerson(p Person) string {
	return fmt.Sprintf("Hola, %s. Tienes %d años.", p.Name, p.Age)
}

// 12. Entrada array y devuelve suma
func sumArray(arr []int) int {
	sum := 0
	for _, n := range arr {
		sum += n
	}
	return sum
}

// 13. Entrada string y devuelve error o nil
func validateString(s string) error {
	if s == "" {
		return errors.New("cadena vacía")
	}
	return nil
}

// 14. Entrada y salida múltiples
func complexFunction(s string, n int) (string, int, error) {
	if s == "" {
		return "", 0, errors.New("string vacío")
	}
	return s, n * 2, nil
}

// 15. Entrada bool y devuelve struct
func flagToPerson(flag bool) Person {
	if flag {
		return Person{Name: "Verdadero", Age: 1}
	}
	return Person{Name: "Falso", Age: 0}
}

// 16. Entrada y salida JSON
func modifyJSON(jsonStr string) (string, error) {
	var p Person
	err := json.Unmarshal([]byte(jsonStr), &p)
	if err != nil {
		return "", err
	}
	p.Age += 1
	data, err := json.Marshal(p)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// ============================================================================
// REGISTRO DE FUNCIONES DISPONIBLES
// ============================================================================

// getExampleFunctions retorna un mapa con todas las funciones de ejemplo
// Esto evita los warnings de "unused" y sirve como documentación
func getExampleFunctions() map[string]interface{} {
	return map[string]interface{}{
		// Funciones sin parámetros
		"returnError":       returnError,
		"returnBool":        returnBool,
		"returnInt":         returnInt,
		"returnString":      returnString,
		"returnStruct":      returnStruct,
		"returnIntArray":    returnIntArray,
		"returnStringArray": returnStringArray,
		"returnJSON":        returnJSON,

		// Funciones con parámetros
		"lengthOfString": lengthOfString,
		"isEven":         isEven,
		"greetPerson":    greetPerson,
		"sumArray":       sumArray,
		"validateString": validateString,
		"flagToPerson":   flagToPerson,
		"modifyJSON":     modifyJSON,

		// Funciones con múltiples valores de retorno
		"complexFunction": complexFunction,
	}
}

// registerExampleFunctions registra todas las funciones de ejemplo
func registerExampleFunctions(h *server.Handler) {
	functions := map[string]interface{}{
		// Funciones sin parámetros
		"returnError":       returnError,
		"returnBool":        returnBool,
		"returnInt":         returnInt,
		"returnString":      returnString,
		"returnStruct":      returnStruct,
		"returnIntArray":    returnIntArray,
		"returnStringArray": returnStringArray,
		"returnJSON":        returnJSON,

		// Funciones con parámetros
		"lengthOfString": lengthOfString,
		"isEven":         isEven,
		"greetPerson":    greetPerson,
		"sumArray":       sumArray,
		"validateString": validateString,
		"flagToPerson":   flagToPerson,
		"modifyJSON":     modifyJSON,

		// Funciones con múltiples valores de retorno
		"complexFunction": complexFunction,
	}

	h.RegisterFunctions(functions)
}

// validateFunctionRegistry verifica que todas las funciones estén disponibles
func validateFunctionRegistry(h *server.Handler) {
	registeredFunctions := h.GetRegisteredFunctions()
	log.Printf("📋 Function registry initialized with %d functions:", len(registeredFunctions))

	for _, name := range registeredFunctions {
		log.Printf("   • %s", name)
	}

	log.Printf("✅ All functions are available for remote execution")
}

// getAvailableFunctions devuelve una lista de nombres de funciones disponibles
func getAvailableFunctions() []string {
	registry := getExampleFunctions()
	var functions []string

	for name := range registry {
		functions = append(functions, name)
	}

	return functions
}

func main() {
	// Definir flags de línea de comandos
	showFunctions := flag.Bool("functions", false, "Mostrar documentación de funciones disponibles")
	showHelp := flag.Bool("help", false, "Mostrar ayuda")
	showList := flag.Bool("list", false, "Listar funciones disponibles")
	flag.Parse()

	// Si se solicita ayuda
	if *showHelp {
		showUsage()
		return
	}

	// Si se solicita mostrar funciones
	if *showFunctions {
		showFunctionDocumentation()
		return
	}

	// Si se solicita listar funciones
	if *showList {
		showFunctionList()
		return
	}

	// Ejecutar servidor principal
	runServer()
}

func showUsage() {
	fmt.Println("🔧 Burrowctl Server")
	fmt.Println("==================")
	fmt.Println()
	fmt.Println("Uso:")
	fmt.Println("  go run server_example.go [opciones]")
	fmt.Println()
	fmt.Println("Opciones:")
	fmt.Println("  -functions    Mostrar documentación de funciones disponibles")
	fmt.Println("  -list         Listar funciones disponibles")
	fmt.Println("  -help         Mostrar esta ayuda")
	fmt.Println()
	fmt.Println("Sin opciones: Ejecuta el servidor burrowctl")
	fmt.Println()
	fmt.Println("Variables de entorno (opcionales):")
	fmt.Println("  BURROWCTL_DEVICE_ID      - ID del dispositivo")
	fmt.Println("  BURROWCTL_AMQP_URL       - URL de RabbitMQ")
	fmt.Println("  BURROWCTL_MYSQL_DSN      - DSN de MariaDB/MySQL")
	fmt.Println("  BURROWCTL_CONNECTION_MODE - Modo de conexión (open/close)")
	fmt.Println()
	fmt.Println("Ejemplos:")
	fmt.Println("  go run server_example.go                # Ejecutar servidor")
	fmt.Println("  go run server_example.go -functions     # Ver funciones")
	fmt.Println("  go run server_example.go -list          # Listar funciones")
	fmt.Println("  go run server_example.go -help          # Ver ayuda")
}

func showFunctionList() {
	fmt.Println("📋 Funciones Disponibles")
	fmt.Println("========================")
	fmt.Println()

	functions := getAvailableFunctions()

	fmt.Printf("Total: %d funciones\n\n", len(functions))

	// Agrupar funciones por tipo
	noParams := []string{}
	withParams := []string{}
	multiReturn := []string{}

	for _, name := range functions {
		switch name {
		case "returnError", "returnBool", "returnInt", "returnString",
			"returnStruct", "returnIntArray", "returnStringArray", "returnJSON":
			noParams = append(noParams, name)
		case "complexFunction":
			multiReturn = append(multiReturn, name)
		default:
			withParams = append(withParams, name)
		}
	}

	if len(noParams) > 0 {
		fmt.Println("🔹 Sin parámetros:")
		for _, name := range noParams {
			fmt.Printf("   • %s\n", name)
		}
		fmt.Println()
	}

	if len(withParams) > 0 {
		fmt.Println("🔹 Con parámetros:")
		for _, name := range withParams {
			fmt.Printf("   • %s\n", name)
		}
		fmt.Println()
	}

	if len(multiReturn) > 0 {
		fmt.Println("🔹 Múltiples valores de retorno:")
		for _, name := range multiReturn {
			fmt.Printf("   • %s\n", name)
		}
		fmt.Println()
	}

	fmt.Println("💡 Para ver documentación detallada: go run server_example.go -functions")
	fmt.Println("⚡ Para ejecutar el servidor: go run server_example.go")
}

func showFunctionDocumentation() {
	fmt.Println("🔧 Funciones Disponibles en Burrowctl Server")
	fmt.Println("============================================")
	fmt.Println()
	fmt.Println("Las siguientes funciones están disponibles para ejecución remota:")
	fmt.Println()
	fmt.Println("📋 1. FUNCIONES SIN PARÁMETROS")
	fmt.Println("   • returnError() error                    - Devuelve un error de ejemplo")
	fmt.Println("   • returnBool() bool                      - Devuelve true")
	fmt.Println("   • returnInt() int                        - Devuelve 42")
	fmt.Println("   • returnString() string                  - Devuelve 'Hola mundo'")
	fmt.Println("   • returnStruct() Person                  - Devuelve struct Person{Name: 'Juan', Age: 30}")
	fmt.Println("   • returnIntArray() []int                 - Devuelve [1, 2, 3, 4, 5]")
	fmt.Println("   • returnStringArray() []string           - Devuelve ['uno', 'dos', 'tres']")
	fmt.Println("   • returnJSON() string                    - Devuelve JSON serializado de Person")
	fmt.Println()
	fmt.Println("📋 2. FUNCIONES CON PARÁMETROS")
	fmt.Println("   • lengthOfString(s string) int           - Devuelve la longitud del string")
	fmt.Println("   • isEven(n int) bool                     - Verifica si el número es par")
	fmt.Println("   • sumArray(arr []int) int                - Suma todos los elementos del array")
	fmt.Println("   • greetPerson(p Person) string           - Saluda a la persona")
	fmt.Println("   • validateString(s string) error         - Valida que el string no esté vacío")
	fmt.Println("   • flagToPerson(flag bool) Person         - Convierte booleano a Person")
	fmt.Println("   • modifyJSON(jsonStr string) (string, error) - Modifica JSON y lo devuelve")
	fmt.Println()
	fmt.Println("📋 3. FUNCIONES CON MÚLTIPLES VALORES DE RETORNO")
	fmt.Println("   • complexFunction(s string, n int) (string, int, error) - Múltiples valores")
	fmt.Println()

	// Ejemplos de solicitudes JSON para cada tipo de función
	examples := []struct {
		name        string
		description string
		request     map[string]interface{}
	}{
		{
			name:        "Función sin parámetros",
			description: "Ejecuta returnString() que devuelve 'Hola mundo'",
			request: map[string]interface{}{
				"name":   "returnString",
				"params": []interface{}{},
			},
		},
		{
			name:        "Función con string",
			description: "Ejecuta lengthOfString('Hello World') que devuelve 11",
			request: map[string]interface{}{
				"name": "lengthOfString",
				"params": []interface{}{
					map[string]interface{}{
						"type":  "string",
						"value": "Hello World",
					},
				},
			},
		},
		{
			name:        "Función con entero",
			description: "Ejecuta isEven(42) que devuelve true",
			request: map[string]interface{}{
				"name": "isEven",
				"params": []interface{}{
					map[string]interface{}{
						"type":  "int",
						"value": 42,
					},
				},
			},
		},
		{
			name:        "Función con array",
			description: "Ejecuta sumArray([1,2,3,4,5]) que devuelve 15",
			request: map[string]interface{}{
				"name": "sumArray",
				"params": []interface{}{
					map[string]interface{}{
						"type":  "[]int",
						"value": []interface{}{1, 2, 3, 4, 5},
					},
				},
			},
		},
		{
			name:        "Función con struct",
			description: "Ejecuta greetPerson(Person) que devuelve saludo personalizado",
			request: map[string]interface{}{
				"name": "greetPerson",
				"params": []interface{}{
					map[string]interface{}{
						"type": "Person",
						"value": map[string]interface{}{
							"name": "María",
							"age":  28,
						},
					},
				},
			},
		},
		{
			name:        "Función con múltiples valores de retorno",
			description: "Ejecuta complexFunction('Go', 10) que devuelve (string, int, error)",
			request: map[string]interface{}{
				"name": "complexFunction",
				"params": []interface{}{
					map[string]interface{}{
						"type":  "string",
						"value": "Go",
					},
					map[string]interface{}{
						"type":  "int",
						"value": 10,
					},
				},
			},
		},
	}

	fmt.Println("🎯 EJEMPLOS DE USO:")
	fmt.Println()

	for i, example := range examples {
		fmt.Printf("%d. %s\n", i+1, example.name)
		fmt.Printf("   Descripción: %s\n", example.description)

		jsonData, err := json.MarshalIndent(example.request, "   ", "  ")
		if err != nil {
			log.Printf("Error marshaling example: %v", err)
			continue
		}

		fmt.Printf("   Solicitud JSON:\n   %s\n", string(jsonData))
		fmt.Printf("   Comando cliente: FUNCTION:%s\n", string(jsonData))
		fmt.Println()
	}

	fmt.Println("💡 CÓMO USAR DESDE EL CLIENTE:")
	fmt.Println("   1. Construye la solicitud JSON con 'name' y 'params'")
	fmt.Println("   2. Envía con prefijo 'FUNCTION:' + JSON")
	fmt.Println("   3. El servidor ejecutará la función y devolverá el resultado")
	fmt.Println()
	fmt.Println("📋 TIPOS DE PARÁMETROS SOPORTADOS:")
	fmt.Println("   • string - Cadenas de texto")
	fmt.Println("   • int - Números enteros")
	fmt.Println("   • bool - Valores booleanos")
	fmt.Println("   • []int - Arrays de enteros")
	fmt.Println("   • []string - Arrays de strings")
	fmt.Println("   • Person - Struct con campos name y age")
	fmt.Println()
	fmt.Println("🔗 PARA MÁS EJEMPLOS:")
	fmt.Println("   • examples/client/function-example/main.go")
	fmt.Println("   • examples/client/function-example/README.md")
	fmt.Println()
	fmt.Println("⚡ COMANDOS ÚTILES:")
	fmt.Println("   go run server_example.go          # Ejecutar servidor")
	fmt.Println("   go run server_example.go -help    # Ver ayuda")
	fmt.Println("   make run-function-example         # Probar cliente")
}

func runServer() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Configure signals to close gracefully
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c
		log.Println("🛑 Closing server...")
		cancel()
	}()

	// Configuration of the connection pool
	pool := &server.PoolConfig{
		MaxIdleConns:    5,
		MaxOpenConns:    15,
		ConnMaxLifetime: 5 * time.Minute,
	}

	// Get connection URLs from environment variables or use defaults for Docker
	deviceID := getEnv("BURROWCTL_DEVICE_ID", "fd1825ec5a7b63f3fa2be9e04154d3b16f676663ba38e23d4ffafa7b0df29efb")
	amqpURL := getEnv("BURROWCTL_AMQP_URL", "amqp://burrowuser:burrowpass123@rabbitmq:5672/")
	mysqlDSN := getEnv("BURROWCTL_MYSQL_DSN", "burrowuser:burrowpass123@tcp(mariadb:3306)/burrowdb?parseTime=true")
	connectionMode := getEnv("BURROWCTL_CONNECTION_MODE", "open")

	// Create the handler with configuration
	h := server.NewHandler(
		deviceID,       // Device ID
		amqpURL,        // RabbitMQ URI
		mysqlDSN,       // MariaDB DSN
		connectionMode, // Connection mode: "open" for connection pool, "close" for per-query connections
		pool,           // Configuration of the pool
	)

	// Register example functions
	registerExampleFunctions(h)

	// Validate function registry
	validateFunctionRegistry(h)

	log.Println("🚀 Starting burrowctl server...")
	log.Printf("📱 Device ID: %s", deviceID)
	log.Printf("🐰 RabbitMQ: %s", amqpURL)
	log.Printf("🗄️  MariaDB: %s", mysqlDSN)
	log.Printf("🔗 Connection mode: %s", connectionMode)
	log.Println("")
	log.Println("✅ Server capabilities:")
	log.Println("   📊 SQL Queries - Execute remote SQL queries")
	log.Println("   🔧 Functions - Execute remote functions with typed parameters")
	log.Println("   ⚡ Commands - Execute system commands")
	log.Println("")
	log.Println("🎯 Example usage:")
	log.Println("   SQL:      SELECT * FROM users")
	log.Println("   Command:  COMMAND:ps aux")
	log.Println("   Function: FUNCTION:{\"name\":\"returnString\",\"params\":[]}")
	log.Println("")
	log.Println("💡 To see function documentation: go run server_example.go -functions")
	log.Println("")

	if err := h.Start(ctx); err != nil {
		log.Fatal("❌ Error starting server:", err)
	}

	log.Println("✅ Server closed gracefully")
}

// getEnv returns environment variable value or default if not set
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
