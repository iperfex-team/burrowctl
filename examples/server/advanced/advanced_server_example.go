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

func main() {
	// Configuración avanzada mediante flags
	var (
		deviceID           = flag.String("device", "fd1825ec5a7b63f3fa2be9e04154d3b16f676663ba38e23d4ffafa7b0df29efb", "Device ID")
		amqpURL            = flag.String("amqp", "amqp://burrowuser:burrowpass123@rabbitmq:5672/", "RabbitMQ URL")
		mysqlDSN           = flag.String("mysql", "burrowuser:burrowpass123@tcp(mariadb:3306)/burrowdb?parseTime=true", "MySQL DSN")
		connectionMode     = flag.String("mode", "open", "Connection mode: 'open' (pooled) or 'close' (per-query)")
		
		// Database pool configuration
		maxIdleConns       = flag.Int("pool-idle", 10, "Maximum idle connections in pool")
		maxOpenConns       = flag.Int("pool-open", 20, "Maximum open connections in pool")
		connMaxLifetime    = flag.Duration("pool-lifetime", 5*time.Minute, "Maximum connection lifetime")
		
		// Worker pool configuration  
		workerCount        = flag.Int("workers", 10, "Number of worker goroutines")
		queueSize          = flag.Int("queue-size", 100, "Worker queue size")
		workerTimeout      = flag.Duration("worker-timeout", 30*time.Second, "Worker task timeout")
		
		// Rate limiter configuration
		rateLimit          = flag.Int("rate-limit", 10, "Requests per second per client")
		burstSize          = flag.Int("burst-size", 20, "Maximum burst size for rate limiting")
		rateLimitCleanup   = flag.Duration("rate-cleanup", 5*time.Minute, "Rate limiter cleanup interval")
		
		showConfig         = flag.Bool("show-config", false, "Show current configuration and exit")
		showHelp           = flag.Bool("help", false, "Show help message")
	)
	flag.Parse()

	if *showHelp {
		showAdvancedHelp()
		return
	}

	if *showConfig {
		showCurrentConfig(*deviceID, *amqpURL, *mysqlDSN, *connectionMode,
			*maxIdleConns, *maxOpenConns, *connMaxLifetime,
			*workerCount, *queueSize, *workerTimeout,
			*rateLimit, *burstSize, *rateLimitCleanup)
		return
	}

	runAdvancedServer(*deviceID, *amqpURL, *mysqlDSN, *connectionMode,
		*maxIdleConns, *maxOpenConns, *connMaxLifetime,
		*workerCount, *queueSize, *workerTimeout,
		*rateLimit, *burstSize, *rateLimitCleanup)
}

func showAdvancedHelp() {
	fmt.Println("🚀 Advanced burrowctl Server")
	fmt.Println("============================")
	fmt.Println()
	fmt.Println("This server demonstrates enterprise-grade features:")
	fmt.Println("• 🏗️  Worker Pool - Configurable concurrent message processing")
	fmt.Println("• 🛡️  Rate Limiting - Per-client request rate control")
	fmt.Println("• 💾 Connection Pooling - Optimized database connections")
	fmt.Println("• ⚡ Performance Tuning - Fine-grained configuration options")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  go run advanced_server_example.go [options]")
	fmt.Println()
	fmt.Println("Connection Options:")
	fmt.Println("  -device string        Device ID (default: SHA256 hash)")
	fmt.Println("  -amqp string          RabbitMQ URL")
	fmt.Println("  -mysql string         MySQL DSN")
	fmt.Println("  -mode string          Connection mode: 'open' or 'close' (default: open)")
	fmt.Println()
	fmt.Println("Database Pool Options:")
	fmt.Println("  -pool-idle int        Max idle connections (default: 10)")
	fmt.Println("  -pool-open int        Max open connections (default: 20)")
	fmt.Println("  -pool-lifetime dur    Connection max lifetime (default: 5m)")
	fmt.Println()
	fmt.Println("Worker Pool Options:")
	fmt.Println("  -workers int          Number of workers (default: 10)")
	fmt.Println("  -queue-size int       Queue buffer size (default: 100)")
	fmt.Println("  -worker-timeout dur   Task timeout (default: 30s)")
	fmt.Println()
	fmt.Println("Rate Limiting Options:")
	fmt.Println("  -rate-limit int       Requests/second per client (default: 10)")
	fmt.Println("  -burst-size int       Max burst tokens (default: 20)")
	fmt.Println("  -rate-cleanup dur     Cleanup interval (default: 5m)")
	fmt.Println()
	fmt.Println("Other Options:")
	fmt.Println("  -show-config          Show current config and exit")
	fmt.Println("  -help                 Show this help")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println()
	fmt.Println("  # High-performance configuration")
	fmt.Println("  go run advanced_server_example.go \\")
	fmt.Println("    -workers=20 -queue-size=500 \\")
	fmt.Println("    -pool-idle=20 -pool-open=50 \\")
	fmt.Println("    -rate-limit=50 -burst-size=100")
	fmt.Println()
	fmt.Println("  # Conservative configuration")
	fmt.Println("  go run advanced_server_example.go \\")
	fmt.Println("    -workers=5 -queue-size=50 \\")
	fmt.Println("    -pool-idle=5 -pool-open=10 \\")
	fmt.Println("    -rate-limit=5 -burst-size=10")
	fmt.Println()
	fmt.Println("  # Show current configuration")
	fmt.Println("  go run advanced_server_example.go -show-config")
}

func showCurrentConfig(deviceID, amqpURL, mysqlDSN, mode string,
	idleConns, openConns int, lifetime time.Duration,
	workers, queueSize int, workerTimeout time.Duration,
	rateLimit, burstSize int, cleanup time.Duration) {
	
	fmt.Println("🔧 Current Server Configuration")
	fmt.Println("===============================")
	fmt.Println()
	fmt.Println("📡 Connection:")
	fmt.Printf("   Device ID: %s\n", deviceID)
	fmt.Printf("   RabbitMQ:  %s\n", amqpURL)
	fmt.Printf("   MySQL:     %s\n", mysqlDSN)
	fmt.Printf("   Mode:      %s\n", mode)
	fmt.Println()
	fmt.Println("💾 Database Pool:")
	fmt.Printf("   Max Idle:     %d connections\n", idleConns)
	fmt.Printf("   Max Open:     %d connections\n", openConns)
	fmt.Printf("   Lifetime:     %v\n", lifetime)
	fmt.Println()
	fmt.Println("🏗️  Worker Pool:")
	fmt.Printf("   Workers:      %d goroutines\n", workers)
	fmt.Printf("   Queue Size:   %d messages\n", queueSize)
	fmt.Printf("   Timeout:      %v per task\n", workerTimeout)
	fmt.Println()
	fmt.Println("🛡️  Rate Limiting:")
	fmt.Printf("   Rate Limit:   %d req/sec per client\n", rateLimit)
	fmt.Printf("   Burst Size:   %d tokens\n", burstSize)
	fmt.Printf("   Cleanup:      %v interval\n", cleanup)
	fmt.Println()
	fmt.Println("📊 Performance Estimates:")
	
	// Calculate theoretical throughput
	maxThroughput := workers * 60 // assuming 1 second average per task
	fmt.Printf("   Max Throughput: ~%d req/min (with %d workers)\n", maxThroughput, workers)
	
	totalRateLimit := rateLimit * 100 // assuming 100 concurrent clients
	fmt.Printf("   Rate Limit Cap: %d req/sec (100 clients)\n", totalRateLimit)
	
	fmt.Printf("   DB Concurrency: %d max connections\n", openConns)
	fmt.Println()
	fmt.Println("💡 To start server: go run advanced_server_example.go")
}

func runAdvancedServer(deviceID, amqpURL, mysqlDSN, mode string,
	idleConns, openConns int, lifetime time.Duration,
	workers, queueSize int, workerTimeout time.Duration,
	rateLimit, burstSize int, cleanup time.Duration) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Configure graceful shutdown
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c
		log.Println("🛑 Shutting down server gracefully...")
		cancel()
	}()

	fmt.Println("🚀 Advanced burrowctl Server Starting")
	fmt.Println("=====================================")
	fmt.Println()

	// Create custom configurations
	poolConfig := &server.PoolConfig{
		MaxIdleConns:    idleConns,
		MaxOpenConns:    openConns,
		ConnMaxLifetime: lifetime,
	}

	// Create handler with custom pool config
	h := server.NewHandler(deviceID, amqpURL, mysqlDSN, mode, poolConfig)

	// Register example functions
	registerExampleFunctions(h)

	// Show configuration
	fmt.Println("📋 Server Configuration:")
	fmt.Printf("   📱 Device ID: %s\n", deviceID)
	fmt.Printf("   🐰 RabbitMQ: %s\n", amqpURL)
	fmt.Printf("   🗄️  MySQL: %s\n", mysqlDSN)
	fmt.Printf("   🔗 Mode: %s\n", mode)
	fmt.Println()
	fmt.Println("💾 Database Pool:")
	fmt.Printf("   ├─ Max Idle: %d\n", idleConns)
	fmt.Printf("   ├─ Max Open: %d\n", openConns)
	fmt.Printf("   └─ Lifetime: %v\n", lifetime)
	fmt.Println()
	fmt.Println("🏗️  Worker Pool:")
	fmt.Printf("   ├─ Workers: %d\n", workers)
	fmt.Printf("   ├─ Queue Size: %d\n", queueSize)
	fmt.Printf("   └─ Timeout: %v\n", workerTimeout)
	fmt.Println()
	fmt.Println("🛡️  Rate Limiting:")
	fmt.Printf("   ├─ Rate: %d req/sec per client\n", rateLimit)
	fmt.Printf("   ├─ Burst: %d tokens\n", burstSize)
	fmt.Printf("   └─ Cleanup: %v\n", cleanup)
	fmt.Println()

	// Show registered functions
	functions := h.GetRegisteredFunctions()
	fmt.Printf("🔧 Registered Functions: %d\n", len(functions))
	for i, name := range functions {
		prefix := "├─"
		if i == len(functions)-1 {
			prefix = "└─"
		}
		fmt.Printf("   %s %s\n", prefix, name)
	}
	fmt.Println()

	fmt.Println("✅ Server Capabilities:")
	fmt.Println("   📊 SQL Queries - Execute remote SQL with connection pooling")
	fmt.Println("   🔧 Functions - Execute typed functions with worker pool")
	fmt.Println("   ⚡ Commands - Execute system commands with timeout")
	fmt.Println("   🛡️  Rate Limited - Protected against abuse")
	fmt.Println("   🔄 Auto Reconnect - Client-side automatic reconnection")
	fmt.Println()

	fmt.Println("🎯 Performance Features Active:")
	fmt.Printf("   • Worker Pool: %d concurrent message processors\n", workers)
	fmt.Printf("   • Connection Pool: %d-%d database connections\n", idleConns, openConns)
	fmt.Printf("   • Rate Limiting: %d req/sec per client (burst: %d)\n", rateLimit, burstSize)
	fmt.Println("   • Prepared Statements: Client-side statement caching")
	fmt.Println("   • Auto Reconnection: Client-side connection recovery")
	fmt.Println()

	fmt.Println("📡 Example Usage:")
	fmt.Println("   SQL:      SELECT * FROM users WHERE id = ?")
	fmt.Println("   Command:  COMMAND:ps aux | grep mysql")
	fmt.Println("   Function: FUNCTION:{\"name\":\"returnString\",\"params\":[]}")
	fmt.Println()

	fmt.Println("🔍 Monitor server with:")
	fmt.Println("   • Watch logs for worker pool activity")
	fmt.Println("   • Rate limiting messages for abuse protection")
	fmt.Println("   • Connection pool utilization in debug mode")
	fmt.Println()

	startTime := time.Now()
	fmt.Printf("⏰ Server started at: %s\n", startTime.Format(time.RFC3339))
	fmt.Println("🎯 Server is ready to accept connections...")
	fmt.Println()

	// Start the server
	if err := h.Start(ctx); err != nil {
		log.Fatal("❌ Error starting server:", err)
	}

	uptime := time.Since(startTime)
	fmt.Printf("\n⏰ Server ran for: %v\n", uptime)
	fmt.Println("✅ Server shutdown completed gracefully")
}

// registerExampleFunctions registers all example functions
func registerExampleFunctions(h *server.Handler) {
	functions := map[string]interface{}{
		// Basic functions
		"returnError":       returnError,
		"returnBool":        returnBool,
		"returnInt":         returnInt,
		"returnString":      returnString,
		"returnStruct":      returnStruct,
		"returnIntArray":    returnIntArray,
		"returnStringArray": returnStringArray,
		"returnJSON":        returnJSON,

		// Functions with parameters
		"lengthOfString": lengthOfString,
		"isEven":         isEven,
		"greetPerson":    greetPerson,
		"sumArray":       sumArray,
		"validateString": validateString,
		"flagToPerson":   flagToPerson,
		"modifyJSON":     modifyJSON,

		// Multi-return functions
		"complexFunction": complexFunction,

		// Performance test functions
		"heavyComputation": heavyComputation,
		"sleepFunction":    sleepFunction,
	}

	h.RegisterFunctions(functions)
}

// ============================================================================
// EXAMPLE FUNCTIONS
// ============================================================================

func returnError() error {
	return errors.New("ejemplo de error")
}

func returnBool() bool {
	return true
}

func returnInt() int {
	return 42
}

func returnString() string {
	return "Hola desde servidor avanzado"
}

func returnStruct() Person {
	return Person{Name: "Usuario", Age: 30}
}

func returnIntArray() []int {
	return []int{1, 2, 3, 4, 5}
}

func returnStringArray() []string {
	return []string{"uno", "dos", "tres"}
}

func returnJSON() string {
	p := Person{Name: "JSON", Age: 25}
	data, _ := json.Marshal(p)
	return string(data)
}

func lengthOfString(s string) int {
	return len(s)
}

func isEven(n int) bool {
	return n%2 == 0
}

func greetPerson(p Person) string {
	return fmt.Sprintf("Hola %s, tienes %d años", p.Name, p.Age)
}

func sumArray(arr []int) int {
	sum := 0
	for _, n := range arr {
		sum += n
	}
	return sum
}

func validateString(s string) error {
	if s == "" {
		return errors.New("cadena vacía")
	}
	return nil
}

func complexFunction(s string, n int) (string, int, error) {
	if s == "" {
		return "", 0, errors.New("string vacío")
	}
	return s, n * 2, nil
}

func flagToPerson(flag bool) Person {
	if flag {
		return Person{Name: "Verdadero", Age: 1}
	}
	return Person{Name: "Falso", Age: 0}
}

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

// Performance test functions
func heavyComputation(iterations int) int {
	result := 0
	for i := 0; i < iterations; i++ {
		result += i * i
	}
	return result
}

func sleepFunction(seconds int) string {
	time.Sleep(time.Duration(seconds) * time.Second)
	return fmt.Sprintf("Slept for %d seconds", seconds)
}