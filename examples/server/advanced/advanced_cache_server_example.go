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
	// Command line flags for cache configuration
	var (
		cacheSize      = flag.Int("cache-size", 1000, "Maximum number of cached queries")
		cacheTTL       = flag.Duration("cache-ttl", 15*time.Minute, "Cache TTL duration")
		cacheCleanup   = flag.Duration("cache-cleanup", 5*time.Minute, "Cache cleanup interval")
		cacheEnabled   = flag.Bool("cache-enabled", true, "Enable query caching")
		workers        = flag.Int("workers", 20, "Number of worker goroutines")
		queueSize      = flag.Int("queue-size", 500, "Worker queue size")
		rateLimit      = flag.Int("rate-limit", 50, "Rate limit per client IP (requests per second)")
		burstSize      = flag.Int("burst-size", 100, "Rate limit burst size")
		poolIdle       = flag.Int("pool-idle", 20, "Maximum idle database connections")
		poolOpen       = flag.Int("pool-open", 50, "Maximum open database connections")
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
		ConnMaxLifetime: 5 * time.Minute,
	}

	// Create handler with advanced configuration
	handler := server.NewHandler(deviceID, amqpURL, mysqlDSN, "open", poolConfig)

	// Configure custom cache settings
	cacheConfig := server.QueryCacheConfig{
		MaxSize:         *cacheSize,
		TTL:             *cacheTTL,
		CleanupInterval: *cacheCleanup,
		Enabled:         *cacheEnabled,
	}
	handler.SetCacheConfig(cacheConfig)

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

	// Register custom functions for testing
	registerCustomFunctions(handler)

	// Display configuration
	fmt.Printf("🚀 Starting Advanced Cache-Enabled Server\n")
	fmt.Printf("========================================\n")
	fmt.Printf("Device ID: %s\n", deviceID)
	fmt.Printf("AMQP URL: %s\n", amqpURL)
	fmt.Printf("MySQL DSN: %s\n", mysqlDSN)
	fmt.Printf("\n📊 Cache Configuration:\n")
	fmt.Printf("  Max Size: %d queries\n", *cacheSize)
	fmt.Printf("  TTL: %v\n", *cacheTTL)
	fmt.Printf("  Cleanup Interval: %v\n", *cacheCleanup)
	fmt.Printf("  Enabled: %v\n", *cacheEnabled)
	fmt.Printf("\n⚙️ Performance Configuration:\n")
	fmt.Printf("  Workers: %d\n", *workers)
	fmt.Printf("  Queue Size: %d\n", *queueSize)
	fmt.Printf("  Rate Limit: %d req/s\n", *rateLimit)
	fmt.Printf("  Burst Size: %d\n", *burstSize)
	fmt.Printf("  DB Pool: %d idle, %d max\n", *poolIdle, *poolOpen)
	fmt.Printf("\n🔄 Starting server...\n")

	// Start server with monitoring
	ctx := context.Background()
	go monitoringLoop(handler)

	if err := handler.Start(ctx); err != nil {
		log.Fatal("Server failed:", err)
	}
}

// registerCustomFunctions registers sample functions for testing
func registerCustomFunctions(handler *server.Handler) {
	// Cache test function
	handler.RegisterFunction("getCacheStats", func() map[string]interface{} {
		stats := handler.GetCacheStats()
		return map[string]interface{}{
			"hits":           stats.Hits,
			"misses":         stats.Misses,
			"hit_ratio":      float64(stats.Hits) / float64(stats.TotalRequests),
			"total_requests": stats.TotalRequests,
			"current_size":   stats.CurrentSize,
			"evictions":      stats.Evictions,
			"expirations":    stats.Expirations,
		}
	})

	// Clear cache function
	handler.RegisterFunction("clearCache", func() string {
		handler.ClearCache()
		return "Cache cleared successfully"
	})

	// System info function
	handler.RegisterFunction("getSystemTime", func() string {
		return time.Now().Format(time.RFC3339)
	})
}

// monitoringLoop prints periodic statistics
func monitoringLoop(handler *server.Handler) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			printStats(handler)
		}
	}
}

// printStats prints current cache and performance statistics
func printStats(handler *server.Handler) {
	cacheStats := handler.GetCacheStats()
	
	fmt.Printf("\n📊 Cache Statistics (Last 30s):\n")
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
	fmt.Printf("  Last Cleanup: %s\n", cacheStats.LastCleanup.Format("15:04:05"))
}