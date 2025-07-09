package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/lordbasex/burrowctl/server"
)

func main() {
	// Command line flags for SQL validation configuration
	var (
		validationEnabled    = flag.Bool("validation-enabled", true, "Enable SQL validation")
		strictMode          = flag.Bool("strict-mode", false, "Enable strict validation mode")
		allowDDL            = flag.Bool("allow-ddl", false, "Allow Data Definition Language commands")
		allowDML            = flag.Bool("allow-dml", true, "Allow Data Manipulation Language commands")
		allowStoredProcs    = flag.Bool("allow-stored-procs", false, "Allow stored procedure calls")
		maxQueryLength      = flag.Int("max-query-length", 10000, "Maximum query length in characters")
		logViolations       = flag.Bool("log-violations", true, "Log validation violations")
		workers             = flag.Int("workers", 20, "Number of worker goroutines")
		queueSize           = flag.Int("queue-size", 500, "Worker queue size")
		rateLimit           = flag.Int("rate-limit", 50, "Rate limit per client IP (requests per second)")
	)
	flag.Parse()

	// Configuration
	deviceID := "my-device"
	amqpURL := "amqp://burrowuser:burrowpass123@localhost:5672/"
	mysqlDSN := "burrowuser:burrowpass123@tcp(localhost:3306)/burrowdb"

	// Create handler with default configuration
	handler := server.NewHandler(deviceID, amqpURL, mysqlDSN, "open", nil)

	// Configure SQL validation with custom settings
	validationConfig := server.SQLValidationConfig{
		Enabled:               *validationEnabled,
		AllowedCommands:       buildAllowedCommands(*allowDDL, *allowDML, *allowStoredProcs),
		BlockedCommands:       buildBlockedCommands(*strictMode),
		AllowDDL:              *allowDDL,
		AllowDML:              *allowDML,
		AllowDQL:              true, // Always allow SELECT queries
		AllowStoredProcedures: *allowStoredProcs,
		MaxQueryLength:        *maxQueryLength,
		StrictMode:            *strictMode,
		LogViolations:         *logViolations,
	}
	handler.SetSQLValidationConfig(validationConfig)

	// Configure worker pool for performance
	workerConfig := &server.WorkerPoolConfig{
		WorkerCount: *workers,
		QueueSize:   *queueSize,
		Timeout:     30 * time.Second,
	}
	handler.SetWorkerPoolConfig(workerConfig)

	// Configure rate limiter
	rateLimiterConfig := &server.RateLimiterConfig{
		RequestsPerSecond: *rateLimit,
		BurstSize:         *rateLimit * 2,
		CleanupInterval:   5 * time.Minute,
	}
	handler.SetRateLimiterConfig(rateLimiterConfig)

	// Register security monitoring functions
	registerSecurityFunctions(handler)

	// Display configuration
	fmt.Printf("🛡️ Starting Security-Focused Server\n")
	fmt.Printf("===================================\n")
	fmt.Printf("Device ID: %s\n", deviceID)
	fmt.Printf("AMQP URL: %s\n", amqpURL)
	fmt.Printf("MySQL DSN: %s\n", mysqlDSN)
	fmt.Printf("\n🔒 SQL Validation Configuration:\n")
	fmt.Printf("  Enabled: %v\n", *validationEnabled)
	fmt.Printf("  Strict Mode: %v\n", *strictMode)
	fmt.Printf("  Allow DDL: %v\n", *allowDDL)
	fmt.Printf("  Allow DML: %v\n", *allowDML)
	fmt.Printf("  Allow Stored Procedures: %v\n", *allowStoredProcs)
	fmt.Printf("  Max Query Length: %d\n", *maxQueryLength)
	fmt.Printf("  Log Violations: %v\n", *logViolations)
	fmt.Printf("  Allowed Commands: %v\n", validationConfig.AllowedCommands)
	fmt.Printf("  Blocked Commands: %v\n", validationConfig.BlockedCommands)
	fmt.Printf("\n⚙️ Performance Configuration:\n")
	fmt.Printf("  Workers: %d\n", *workers)
	fmt.Printf("  Queue Size: %d\n", *queueSize)
	fmt.Printf("  Rate Limit: %d req/s\n", *rateLimit)
	fmt.Printf("\n🔄 Starting server...\n")

	// Start server with security monitoring
	ctx := context.Background()
	go securityMonitoringLoop(handler)

	if err := handler.Start(ctx); err != nil {
		log.Fatal("Server failed:", err)
	}
}

// buildAllowedCommands constructs the list of allowed SQL commands based on flags
func buildAllowedCommands(allowDDL, allowDML, allowStoredProcs bool) []string {
	var commands []string

	// Always allow basic query commands
	commands = append(commands, "SELECT", "SHOW", "DESCRIBE", "EXPLAIN")

	// Add DML commands if allowed
	if allowDML {
		commands = append(commands, "INSERT", "UPDATE", "DELETE")
	}

	// Add DDL commands if allowed
	if allowDDL {
		commands = append(commands, "CREATE", "ALTER", "DROP", "TRUNCATE")
	}

	// Add stored procedure commands if allowed
	if allowStoredProcs {
		commands = append(commands, "CALL", "EXEC", "EXECUTE")
	}

	return commands
}

// buildBlockedCommands constructs the list of explicitly blocked commands
func buildBlockedCommands(strictMode bool) []string {
	blocked := []string{
		"SHUTDOWN", "RESTART", "RESET",
		"CREATE USER", "DROP USER", "ALTER USER",
		"GRANT", "REVOKE", "FLUSH",
		"LOAD DATA", "LOAD XML",
		"INTO OUTFILE", "INTO DUMPFILE",
	}

	// In strict mode, add more restricted commands
	if strictMode {
		blocked = append(blocked,
			"TRUNCATE", "DELETE",  // Restrict bulk operations
			"ALTER", "CREATE",     // Restrict schema changes
			"CALL", "EXECUTE",     // Restrict stored procedures
		)
	}

	return blocked
}

// registerSecurityFunctions registers functions for monitoring security
func registerSecurityFunctions(handler *server.Handler) {
	// Get validation statistics
	handler.RegisterFunction("getValidationStats", func() map[string]interface{} {
		stats := handler.GetSQLValidationStats()
		return map[string]interface{}{
			"total_queries":        stats.TotalQueries,
			"valid_queries":        stats.ValidQueries,
			"blocked_queries":      stats.BlockedQueries,
			"injection_attempts":   stats.InjectionAttempts,
			"command_violations":   stats.CommandViolations,
			"structure_violations": stats.StructureViolations,
			"block_rate":          float64(stats.BlockedQueries) / float64(stats.TotalQueries),
		}
	})

	// Get system security status
	handler.RegisterFunction("getSecurityStatus", func() map[string]interface{} {
		validationStats := handler.GetSQLValidationStats()
		cacheStats := handler.GetCacheStats()
		
		// Calculate security metrics
		blockRate := float64(0)
		if validationStats.TotalQueries > 0 {
			blockRate = float64(validationStats.BlockedQueries) / float64(validationStats.TotalQueries)
		}

		injectionRate := float64(0)
		if validationStats.TotalQueries > 0 {
			injectionRate = float64(validationStats.InjectionAttempts) / float64(validationStats.TotalQueries)
		}

		return map[string]interface{}{
			"security_level":    getSecurityLevel(blockRate, injectionRate),
			"block_rate":        blockRate,
			"injection_rate":    injectionRate,
			"cache_hit_ratio":   float64(cacheStats.Hits) / float64(cacheStats.TotalRequests),
			"total_requests":    validationStats.TotalQueries,
			"threat_detected":   validationStats.InjectionAttempts > 0,
		}
	})

	// Test validation with custom query
	handler.RegisterFunction("testValidation", func(query string) map[string]interface{} {
		// This is for testing purposes - normally you wouldn't expose validation internals
		return map[string]interface{}{
			"test_query": query,
			"message":    "Use validation example client to test queries",
		}
	})
}

// getSecurityLevel determines the current security threat level
func getSecurityLevel(blockRate, injectionRate float64) string {
	if injectionRate > 0.1 { // More than 10% injection attempts
		return "HIGH"
	} else if blockRate > 0.2 { // More than 20% blocked queries
		return "MEDIUM"
	} else if injectionRate > 0.01 { // More than 1% injection attempts
		return "ELEVATED"
	} else {
		return "LOW"
	}
}

// securityMonitoringLoop prints periodic security statistics
func securityMonitoringLoop(handler *server.Handler) {
	ticker := time.NewTicker(60 * time.Second) // Every minute
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			printSecurityStats(handler)
		}
	}
}

// printSecurityStats prints current security statistics
func printSecurityStats(handler *server.Handler) {
	stats := handler.GetSQLValidationStats()
	
	fmt.Printf("\n🛡️ Security Statistics (Last 60s):\n")
	fmt.Printf("  Total Queries: %d\n", stats.TotalQueries)
	fmt.Printf("  Valid Queries: %d\n", stats.ValidQueries)
	fmt.Printf("  Blocked Queries: %d\n", stats.BlockedQueries)
	fmt.Printf("  Injection Attempts: %d\n", stats.InjectionAttempts)
	fmt.Printf("  Command Violations: %d\n", stats.CommandViolations)
	fmt.Printf("  Structure Violations: %d\n", stats.StructureViolations)
	
	if stats.TotalQueries > 0 {
		blockRate := float64(stats.BlockedQueries) / float64(stats.TotalQueries) * 100
		injectionRate := float64(stats.InjectionAttempts) / float64(stats.TotalQueries) * 100
		
		fmt.Printf("  Block Rate: %.2f%%\n", blockRate)
		fmt.Printf("  Injection Rate: %.2f%%\n", injectionRate)
		fmt.Printf("  Security Level: %s\n", getSecurityLevel(blockRate/100, injectionRate/100))
		
		// Security alerts
		if injectionRate > 5 {
			fmt.Printf("  🚨 HIGH INJECTION RATE DETECTED!\n")
		}
		if blockRate > 30 {
			fmt.Printf("  ⚠️  High block rate - consider policy adjustment\n")
		}
	}
}