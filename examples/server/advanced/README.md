# Advanced Server Example

An enhanced burrowctl server implementation with enterprise features for high-performance environments.

## Features

- **Worker Pool**: Configurable concurrent request processing
- **Rate Limiting**: Per-client IP rate limiting with burst support
- **Connection Pooling**: Optimized database connection management
- **Performance Monitoring**: Real-time performance metrics
- **Graceful Shutdown**: Clean server shutdown with request draining

## Usage

### Direct execution
```bash
go run main.go
```

### Using Makefile
```bash
make run-server-advanced
```

### Docker
```bash
make docker-up-advanced
```

## Configuration

### Command Line Options

```bash
go run main.go [options]
```

#### Performance Options
- `-workers=20`: Number of worker goroutines (default: 10)
- `-queue-size=500`: Worker queue size (default: 100)
- `-rate-limit=50`: Rate limit per client IP (req/s) (default: 10)
- `-burst-size=100`: Rate limit burst size (default: 20)

#### Database Options
- `-pool-idle=20`: Maximum idle connections (default: 5)
- `-pool-open=50`: Maximum open connections (default: 15)
- `-device=<id>`: Device ID for identification
- `-amqp=<url>`: AMQP connection URL
- `-mysql=<dsn>`: MySQL connection DSN

### Example Configurations

#### High Performance Mode
```bash
go run main.go -workers=50 -queue-size=1000 -rate-limit=100 -burst-size=200
```

#### Database Intensive Mode
```bash
go run main.go -pool-idle=50 -pool-open=100 -workers=30
```

## Performance Features

### Worker Pool
- Configurable number of concurrent workers
- Bounded queue with overflow protection
- Load balancing across workers
- Graceful shutdown support

### Rate Limiting
- Per-client IP rate limiting
- Token bucket algorithm with burst support
- Configurable rates and burst sizes
- Automatic cleanup of rate limit data

### Connection Pooling
- Optimized database connection management
- Configurable idle/open connection limits
- Connection lifetime management
- Health checking and reconnection

## Monitoring

The server provides real-time performance monitoring:

```
🚀 Starting Advanced Enterprise Server...
📊 Performance Configuration:
  Workers: 20
  Queue Size: 500
  Rate Limit: 50 req/s
  Burst Size: 100
  Pool Idle: 20
  Pool Open: 50

⚡ Performance Metrics (every 30s):
  Active Workers: 15/20
  Queue Usage: 45/500
  Rate Limit Hits: 12
  DB Connections: 18/20 idle, 35/50 open
  Requests/sec: 42.5
  Avg Response Time: 15ms
```

## Architecture

### Request Flow
1. **Rate Limiter**: Checks client IP rate limits
2. **Worker Pool**: Assigns request to available worker
3. **Connection Pool**: Manages database connections
4. **Processing**: Executes the request
5. **Response**: Returns result to client

### Components
- **RateLimiter**: Token bucket rate limiting per client IP
- **WorkerPool**: Concurrent request processing
- **ConnectionPool**: Database connection management
- **PerformanceMonitor**: Real-time metrics collection

## Performance Tuning

### Workers
- **Low traffic**: 5-10 workers
- **Medium traffic**: 20-30 workers
- **High traffic**: 50+ workers

### Queue Size
- Should be 10-50x the number of workers
- Larger queues provide better burst handling
- Monitor queue usage to avoid memory issues

### Rate Limiting
- **Development**: 10-20 req/s per client
- **Production**: 50-100 req/s per client
- **High volume**: 200+ req/s per client

### Database Connections
- **Idle connections**: 20-50% of max open
- **Max open**: Based on database server capacity
- **Connection lifetime**: 5-30 minutes

## Compared to Basic Server

| Feature | Basic Server | Advanced Server |
|---------|-------------|-----------------|
| Concurrency | Single-threaded | Worker pool |
| Rate Limiting | None | Per-client IP |
| Connection Pooling | Basic | Advanced with tuning |
| Monitoring | Basic logs | Performance metrics |
| Graceful Shutdown | None | Full support |
| Configuration | Hardcoded | Command-line flags |

## Next Steps

For specialized configurations:
- [Cache Server](cache-server/README.md) - Query result caching
- [Validation Server](validation-server/README.md) - SQL security validation
- [Full-Featured Server](full-featured-server/README.md) - All enterprise features