# Full-Featured Enterprise Server

This is a comprehensive example of a burrowctl server with all enterprise features enabled:

## Features

- **Query Caching**: High-performance query result caching with configurable TTL
- **SQL Validation**: Multi-layered security validation with injection detection
- **Worker Pool**: Configurable worker pool for concurrent request processing
- **Rate Limiting**: Per-client IP rate limiting with burst support
- **Connection Pooling**: Optimized database connection pooling
- **Comprehensive Monitoring**: Real-time performance and security metrics

## Usage

### Quick Start

```bash
# Run with default settings
make run-server-full

# Or directly with go
go run main.go
```

### Advanced Configuration

```bash
# Custom cache settings
go run main.go -cache-size=5000 -cache-ttl=30m

# Strict security mode
go run main.go -validation-enabled=true -strict-mode=true -allow-ddl=false

# High performance mode
go run main.go -workers=50 -queue-size=2000 -rate-limit=200

# Full enterprise configuration
go run main.go \
  -cache-enabled=true -cache-size=10000 -cache-ttl=1h \
  -validation-enabled=true -strict-mode=false -allow-ddl=false \
  -workers=100 -queue-size=5000 -rate-limit=500 -burst-size=1000 \
  -pool-idle=50 -pool-open=200 -monitoring-enabled=true
```

## Configuration Options

### Cache Configuration
- `-cache-enabled`: Enable/disable query caching (default: true)
- `-cache-size`: Maximum cached queries (default: 2000)
- `-cache-ttl`: Cache time-to-live (default: 15m)
- `-cache-cleanup`: Cleanup interval (default: 5m)

### SQL Validation
- `-validation-enabled`: Enable SQL validation (default: true)
- `-strict-mode`: Enable strict validation mode (default: false)
- `-allow-ddl`: Allow Data Definition Language (default: false)
- `-allow-dml`: Allow Data Manipulation Language (default: true)
- `-allow-stored-procs`: Allow stored procedures (default: false)
- `-max-query-length`: Maximum query length (default: 10000)
- `-log-violations`: Log validation violations (default: true)

### Performance
- `-workers`: Number of worker goroutines (default: 25)
- `-queue-size`: Worker queue size (default: 1000)
- `-rate-limit`: Rate limit per client IP (default: 100 req/s)
- `-burst-size`: Rate limit burst size (default: 200)

### Database
- `-pool-idle`: Maximum idle connections (default: 25)
- `-pool-open`: Maximum open connections (default: 75)
- `-conn-lifetime`: Connection lifetime (default: 10m)

### Monitoring
- `-monitoring-enabled`: Enable periodic monitoring (default: true)
- `-monitoring-interval`: Monitoring report interval (default: 60s)

## Docker Usage

```bash
# Start with Docker
make docker-up-full

# View logs
make docker-logs-full

# Stop
make docker-down-full
```

## Monitoring Functions

The server registers several monitoring functions:

- `getCacheStats()`: Cache performance statistics
- `getValidationStats()`: SQL validation and security metrics
- `getSystemStatus()`: Overall system health
- `getPerformanceMetrics()`: Performance analysis
- `clearAllCaches()`: Administrative cache clearing

## Security Features

### SQL Injection Detection
- Pattern-based injection detection
- Command whitelist/blacklist enforcement
- Structural query validation
- Parameter validation
- Risk level assessment

### Rate Limiting
- Per-client IP rate limiting
- Configurable burst capacity
- Automatic cleanup of rate limit data

### Database Security
- Connection pooling with lifetime management
- Prepared statement support
- Transaction isolation

## Performance Optimizations

### Query Caching
- LRU cache with configurable size
- TTL-based expiration
- Automatic cache cleanup
- Cache hit/miss statistics

### Worker Pool
- Configurable worker count
- Bounded queue with overflow protection
- Graceful shutdown support
- Load balancing across workers

### Connection Management
- Connection pooling with idle/open limits
- Connection lifetime management
- Health checking and reconnection
- Graceful connection draining

## Monitoring Output

The server provides comprehensive monitoring output every 60 seconds (configurable):

```
📊 COMPREHENSIVE SYSTEM REPORT - 14:30:15
============================================================
🏢 System Overview:
  Uptime: 2h45m30s

📈 Cache Performance:
  Total Requests: 15,432
  Cache Hits: 12,345
  Cache Misses: 3,087
  Hit Ratio: 80.00%
  Current Size: 1,234 entries
  Evictions: 45
  Expirations: 123

🛡️ Security & Validation:
  Total Queries: 14,567
  Valid Queries: 14,445
  Blocked Queries: 122
  Injection Attempts: 5
  Command Violations: 87
  Structure Violations: 30
  Block Rate: 0.84%
  Injection Rate: 0.03%
  Security Level: LOW

⚡ Performance Summary:
  Cache Efficiency: 80.00%
  Validation Efficiency: 99.16%
============================================================
```

## Example Client Usage

Once the server is running, you can connect clients to test the features:

```bash
# Test basic functionality
cd ../../../client/command-example
go run main.go

# Test SQL with caching
cd ../sql-example
go run main.go

# Test validation
cd ../validation-example
go run main.go
```