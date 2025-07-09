package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/lordbasex/burrowctl/server"
)

func main() {
	// Command line flags for complete server configuration
	var (
		// Cache configuration
		cacheEnabled     = flag.Bool("cache-enabled", true, "Enable query caching")
		cacheSize        = flag.Int("cache-size", 2000, "Maximum number of cached queries")
		cacheTTL         = flag.Duration("cache-ttl", 15*time.Minute, "Cache TTL duration")
		cacheCleanup     = flag.Duration("cache-cleanup", 5*time.Minute, "Cache cleanup interval")
		
		// SQL Validation configuration
		validationEnabled  = flag.Bool("validation-enabled", true, "Enable SQL validation")
		strictMode        = flag.Bool("strict-mode", false, "Enable strict validation mode")
		allowDDL          = flag.Bool("allow-ddl", false, "Allow Data Definition Language commands")
		allowDML          = flag.Bool("allow-dml", true, "Allow Data Manipulation Language commands")
		allowStoredProcs  = flag.Bool("allow-stored-procs", false, "Allow stored procedure calls")
		maxQueryLength    = flag.Int("max-query-length", 10000, "Maximum query length in characters")
		logViolations     = flag.Bool("log-violations", true, "Log validation violations")
		
		// Performance configuration
		workers           = flag.Int("workers", 25, "Number of worker goroutines")
		queueSize         = flag.Int("queue-size", 1000, "Worker queue size")
		rateLimit         = flag.Int("rate-limit", 100, "Rate limit per client IP (requests per second)")
		burstSize         = flag.Int("burst-size", 200, "Rate limit burst size")
		
		// Database configuration
		poolIdle          = flag.Int("pool-idle", 25, "Maximum idle database connections")
		poolOpen          = flag.Int("pool-open", 75, "Maximum open database connections")
		connLifetime      = flag.Duration("conn-lifetime", 10*time.Minute, "Database connection lifetime")
		
		// Monitoring configuration
		monitoringEnabled = flag.Bool("monitoring-enabled", true, "Enable periodic monitoring")
		monitoringInterval = flag.Duration("monitoring-interval", 60*time.Second, "Monitoring report interval")
	)
	flag.Parse()

	// Configuration
	deviceID := "my-device"
	amqpURL := "amqp://burrowuser:burrowpass123@localhost:5672/"
	mysqlDSN := "burrowuser:burrowpass123@tcp(localhost:3306)/burrowdb"

	// Advanced pool configuration
	poolConfig := &server.PoolConfig{
		MaxIdleConns:    *poolIdle,
		MaxOpenConns:    *poolOpen,
		ConnMaxLifetime: *connLifetime,
	}

	// Create handler with advanced configuration
	handler := server.NewHandler(deviceID, amqpURL, mysqlDSN, "open", poolConfig)

	// Configure query cache
	cacheConfig := server.QueryCacheConfig{
		MaxSize:         *cacheSize,
		TTL:             *cacheTTL,
		CleanupInterval: *cacheCleanup,
		Enabled:         *cacheEnabled,
	}
	handler.SetCacheConfig(cacheConfig)

	// Configure SQL validation
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

	// Configure worker pool
	workerConfig := &server.WorkerPoolConfig{
		WorkerCount: *workers,
		QueueSize:   *queueSize,
		Timeout:     30 * time.Second,
	}
	handler.SetWorkerPoolConfig(workerConfig)

	// Configure rate limiter
	rateLimiterConfig := &server.RateLimiterConfig{
		RequestsPerSecond: *rateLimit,
		BurstSize:         *burstSize,
		CleanupInterval:   5 * time.Minute,
	}
	handler.SetRateLimiterConfig(rateLimiterConfig)

	// Register comprehensive monitoring functions
	registerMonitoringFunctions(handler)

	// Display complete configuration
	displayConfiguration(deviceID, amqpURL, mysqlDSN, &cacheConfig, &validationConfig, 
		workerConfig, rateLimiterConfig, poolConfig)

	// Start server with comprehensive monitoring
	ctx := context.Background()
	if *monitoringEnabled {
		go comprehensiveMonitoringLoop(handler, *monitoringInterval)
	}

	log.Printf("🚀 Starting Full-Featured Enterprise Server...")
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

// registerMonitoringFunctions registers comprehensive monitoring functions
func registerMonitoringFunctions(handler *server.Handler) {
	// Cache statistics
	handler.RegisterFunction("getCacheStats", func() map[string]interface{} {
		stats := handler.GetCacheStats()
		hitRatio := float64(0)
		if stats.TotalRequests > 0 {
			hitRatio = float64(stats.Hits) / float64(stats.TotalRequests)
		}
		return map[string]interface{}{
			"hits":            stats.Hits,
			"misses":          stats.Misses,
			"hit_ratio":       hitRatio,
			"total_requests":  stats.TotalRequests,
			"current_size":    stats.CurrentSize,
			"evictions":       stats.Evictions,
			"expirations":     stats.Expirations,
			"last_cleanup":    stats.LastCleanup.Format(time.RFC3339),
		}
	})

	// Validation statistics
	handler.RegisterFunction("getValidationStats", func() map[string]interface{} {
		stats := handler.GetSQLValidationStats()
		blockRate := float64(0)
		injectionRate := float64(0)
		if stats.TotalQueries > 0 {
			blockRate = float64(stats.BlockedQueries) / float64(stats.TotalQueries)
			injectionRate = float64(stats.InjectionAttempts) / float64(stats.TotalQueries)
		}
		return map[string]interface{}{
			"total_queries":        stats.TotalQueries,
			"valid_queries":        stats.ValidQueries,
			"blocked_queries":      stats.BlockedQueries,
			"injection_attempts":   stats.InjectionAttempts,
			"command_violations":   stats.CommandViolations,
			"structure_violations": stats.StructureViolations,
			"block_rate":          blockRate,
			"injection_rate":      injectionRate,
			"security_level":      getSecurityLevel(blockRate, injectionRate),
		}
	})

	// Overall system status
	handler.RegisterFunction("getSystemStatus", func() map[string]interface{} {
		cacheStats := handler.GetCacheStats()
		validationStats := handler.GetSQLValidationStats()
		
		cacheHitRatio := float64(0)
		if cacheStats.TotalRequests > 0 {
			cacheHitRatio = float64(cacheStats.Hits) / float64(cacheStats.TotalRequests)
		}
		
		blockRate := float64(0)
		if validationStats.TotalQueries > 0 {
			blockRate = float64(validationStats.BlockedQueries) / float64(validationStats.TotalQueries)
		}
		
		return map[string]interface{}{
			"status":              "healthy",
			"uptime":              time.Since(startTime).String(),
			"cache_hit_ratio":     cacheHitRatio,
			"cache_size":          cacheStats.CurrentSize,
			"validation_enabled":  validationStats.TotalQueries > 0,
			"security_level":      getSecurityLevel(blockRate, 0),
			"total_queries":       validationStats.TotalQueries,
			"blocked_queries":     validationStats.BlockedQueries,
			"injection_attempts":  validationStats.InjectionAttempts,
		}
	})

	// Performance metrics
	handler.RegisterFunction("getPerformanceMetrics", func() map[string]interface{} {
		cacheStats := handler.GetCacheStats()
		validationStats := handler.GetSQLValidationStats()
		
		return map[string]interface{}{
			"cache_performance": map[string]interface{}{
				"hit_ratio":      float64(cacheStats.Hits) / float64(cacheStats.TotalRequests),
				"evictions":      cacheStats.Evictions,
				"current_size":   cacheStats.CurrentSize,
			},
			"validation_performance": map[string]interface{}{
				"total_queries":    validationStats.TotalQueries,
				"blocked_queries":  validationStats.BlockedQueries,
				"processing_rate":  float64(validationStats.ValidQueries) / float64(validationStats.TotalQueries),
			},
		}
	})

	// Clear all caches and stats
	handler.RegisterFunction("clearAllCaches", func() string {
		handler.ClearCache()
		return "All caches cleared successfully"
	})
}

// getSecurityLevel determines the current security threat level
func getSecurityLevel(blockRate, injectionRate float64) string {
	if injectionRate > 0.1 {
		return "HIGH"
	} else if blockRate > 0.2 {
		return "MEDIUM"
	} else if injectionRate > 0.01 {
		return "ELEVATED"
	} else {
		return "LOW"
	}
}

// displayConfiguration shows the complete server configuration
func displayConfiguration(deviceID, amqpURL, mysqlDSN string, cacheConfig *server.QueryCacheConfig, 
	validationConfig *server.SQLValidationConfig, workerConfig *server.WorkerPoolConfig, 
	rateLimiterConfig *server.RateLimiterConfig, poolConfig *server.PoolConfig) {
	
	fmt.Printf("🏢 Full-Featured Enterprise Server Configuration\n")
	fmt.Printf("===============================================\n")
	fmt.Printf("Device ID: %s\n", deviceID)
	fmt.Printf("AMQP URL: %s\n", amqpURL)
	fmt.Printf("MySQL DSN: %s\n", mysqlDSN)
	
	fmt.Printf("\n📊 Cache Configuration:\n")
	fmt.Printf("  Enabled: %v\n", cacheConfig.Enabled)
	fmt.Printf("  Max Size: %d queries\n", cacheConfig.MaxSize)
	fmt.Printf("  TTL: %v\n", cacheConfig.TTL)
	fmt.Printf("  Cleanup Interval: %v\n", cacheConfig.CleanupInterval)
	
	fmt.Printf("\n🔒 SQL Validation Configuration:\n")
	fmt.Printf("  Enabled: %v\n", validationConfig.Enabled)
	fmt.Printf("  Strict Mode: %v\n", validationConfig.StrictMode)
	fmt.Printf("  Allow DDL: %v\n", validationConfig.AllowDDL)
	fmt.Printf("  Allow DML: %v\n", validationConfig.AllowDML)
	fmt.Printf("  Allow Stored Procedures: %v\n", validationConfig.AllowStoredProcedures)
	fmt.Printf("  Max Query Length: %d\n", validationConfig.MaxQueryLength)
	fmt.Printf("  Allowed Commands: %v\n", validationConfig.AllowedCommands)
	
	fmt.Printf("\n⚙️ Performance Configuration:\n")
	fmt.Printf("  Workers: %d\n", workerConfig.WorkerCount)
	fmt.Printf("  Queue Size: %d\n", workerConfig.QueueSize)
	fmt.Printf("  Rate Limit: %d req/s\n", rateLimiterConfig.RequestsPerSecond)
	fmt.Printf("  Burst Size: %d\n", rateLimiterConfig.BurstSize)
	
	fmt.Printf("\n🗄️ Database Configuration:\n")
	fmt.Printf("  Max Idle Connections: %d\n", poolConfig.MaxIdleConns)
	fmt.Printf("  Max Open Connections: %d\n", poolConfig.MaxOpenConns)
	fmt.Printf("  Connection Lifetime: %v\n", poolConfig.ConnMaxLifetime)
	
	fmt.Printf("\n🔄 Starting comprehensive monitoring...\n")
}

// Global start time for uptime calculation
var startTime = time.Now()

// comprehensiveMonitoringLoop provides detailed monitoring of all server components
func comprehensiveMonitoringLoop(handler *server.Handler, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			printComprehensiveStats(handler)
		}
	}
}

// printComprehensiveStats prints detailed statistics for all components
func printComprehensiveStats(handler *server.Handler) {
	cacheStats := handler.GetCacheStats()
	validationStats := handler.GetSQLValidationStats()
	
	fmt.Printf("\n" + strings.Repeat("=", 60) + "\n")
	fmt.Printf("📊 COMPREHENSIVE SYSTEM REPORT - %s\n", time.Now().Format("15:04:05"))
	fmt.Printf(strings.Repeat("=", 60) + "\n")
	
	// System Overview
	fmt.Printf("🏢 System Overview:\n")
	fmt.Printf("  Uptime: %v\n", time.Since(startTime).Round(time.Second))
	
	// Cache Statistics
	fmt.Printf("\n📈 Cache Performance:\n")
	fmt.Printf("  Total Requests: %d\n", cacheStats.TotalRequests)
	fmt.Printf("  Cache Hits: %d\n", cacheStats.Hits)
	fmt.Printf("  Cache Misses: %d\n", cacheStats.Misses)
	if cacheStats.TotalRequests > 0 {
		hitRatio := float64(cacheStats.Hits) / float64(cacheStats.TotalRequests) * 100
		fmt.Printf("  Hit Ratio: %.2f%%\n", hitRatio)
	}
	fmt.Printf("  Current Size: %d entries\n", cacheStats.CurrentSize)
	fmt.Printf("  Evictions: %d\n", cacheStats.Evictions)
	fmt.Printf("  Expirations: %d\n", cacheStats.Expirations)
	
	// Validation Statistics
	fmt.Printf("\n🛡️ Security & Validation:\n")
	fmt.Printf("  Total Queries: %d\n", validationStats.TotalQueries)
	fmt.Printf("  Valid Queries: %d\n", validationStats.ValidQueries)
	fmt.Printf("  Blocked Queries: %d\n", validationStats.BlockedQueries)
	fmt.Printf("  Injection Attempts: %d\n", validationStats.InjectionAttempts)
	fmt.Printf("  Command Violations: %d\n", validationStats.CommandViolations)
	fmt.Printf("  Structure Violations: %d\n", validationStats.StructureViolations)
	
	if validationStats.TotalQueries > 0 {
		blockRate := float64(validationStats.BlockedQueries) / float64(validationStats.TotalQueries) * 100
		injectionRate := float64(validationStats.InjectionAttempts) / float64(validationStats.TotalQueries) * 100
		
		fmt.Printf("  Block Rate: %.2f%%\n", blockRate)
		fmt.Printf("  Injection Rate: %.2f%%\n", injectionRate)
		fmt.Printf("  Security Level: %s\n", getSecurityLevel(blockRate/100, injectionRate/100))
		
		// Security alerts
		if injectionRate > 5 {
			fmt.Printf("  🚨 HIGH INJECTION RATE DETECTED!\n")
		}
		if blockRate > 30 {
			fmt.Printf("  ⚠️  High block rate - review policies\n")
		}
	}
	
	// Performance Summary
	fmt.Printf("\n⚡ Performance Summary:\n")
	if cacheStats.TotalRequests > 0 && validationStats.TotalQueries > 0 {
		cacheEfficiency := float64(cacheStats.Hits) / float64(cacheStats.TotalRequests) * 100
		validationEfficiency := float64(validationStats.ValidQueries) / float64(validationStats.TotalQueries) * 100
		
		fmt.Printf("  Cache Efficiency: %.2f%%\n", cacheEfficiency)
		fmt.Printf("  Validation Efficiency: %.2f%%\n", validationEfficiency)
		
		// Performance recommendations
		if cacheEfficiency < 50 {
			fmt.Printf("  💡 Consider increasing cache size or TTL\n")
		}
		if validationEfficiency < 80 {
			fmt.Printf("  💡 Review validation policies - high block rate\n")
		}
	}
	
	fmt.Printf(strings.Repeat("=", 60) + "\n")
}