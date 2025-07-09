# BurrowCtl Server Examples

This directory contains server examples demonstrating different capabilities and configurations of burrowctl.

## Structure

```
examples/server/
├── Dockerfile              # Universal Dockerfile for all server examples
├── basic/                  # Simple server implementation
│   ├── main.go            # Main server file
│   ├── README.md          # English documentation
│   ├── README.es.md       # Spanish documentation
│   └── README.pt.md       # Portuguese documentation
└── advanced/              # Enterprise server with performance features
    ├── main.go            # Advanced server implementation
    ├── README.md, README.es.md, README.pt.md
    ├── cache-server/      # Specialized cache server
    │   ├── main.go
    │   └── README.{md,es.md,pt.md}
    ├── validation-server/ # SQL security validation server
    │   ├── main.go
    │   └── README.{md,es.md,pt.md}
    └── full-featured-server/ # Complete enterprise server
        ├── main.go
        └── README.{md,es.md,pt.md}
```

## Quick Start

### Basic Server
```bash
make run-server-example
# or
cd basic && go run main.go
```

### Advanced Server
```bash
make run-server-advanced
# or  
cd advanced && go run main.go
```

### Specialized Servers
```bash
make run-server-cache        # Cache-optimized server
make run-server-validation   # SQL validation server
make run-server-full         # Full enterprise server
```

## Docker Usage

Each example includes Docker support:

```bash
make docker-up              # Basic server
make docker-up-advanced     # Advanced server
make docker-up-cache        # Cache server
make docker-up-validation   # Validation server
make docker-up-full         # Full-featured server
```

## Universal Dockerfile

All server examples use a single `Dockerfile` located in this directory. It accepts a build argument to specify which example to build:

```bash
# Build basic server
docker build --build-arg EXAMPLE_DIR=basic -t burrowctl-basic .

# Build advanced server
docker build --build-arg EXAMPLE_DIR=advanced -t burrowctl-advanced .

# Build specialized servers
docker build --build-arg EXAMPLE_DIR=advanced/cache-server -t burrowctl-cache .
docker build --build-arg EXAMPLE_DIR=advanced/validation-server -t burrowctl-validation .
docker build --build-arg EXAMPLE_DIR=advanced/full-featured-server -t burrowctl-full .
```

## Server Comparison

| Feature | Basic | Advanced | Cache | Validation | Full |
|---------|-------|----------|-------|------------|------|
| Worker Pool | ❌ | ✅ | ✅ | ✅ | ✅ |
| Rate Limiting | ❌ | ✅ | ✅ | ✅ | ✅ |
| Connection Pooling | Basic | Advanced | Advanced | Advanced | Advanced |
| Query Caching | ❌ | ❌ | ✅ | ❌ | ✅ |
| SQL Validation | ❌ | ❌ | ❌ | ✅ | ✅ |
| Monitoring | Basic | Metrics | Cache Stats | Security Stats | Comprehensive |
| Configuration | Hardcoded | CLI Flags | CLI Flags | CLI Flags | CLI Flags |

## Development

### File Naming Convention
- All main files are named `main.go` for consistency
- Binaries are automatically excluded from git via `.gitignore`
- Use `make clean-examples` to remove compiled binaries

### Adding New Examples
1. Create new directory under `basic/` or `advanced/`
2. Add `main.go` file
3. Create trilingual documentation (README.md, README.es.md, README.pt.md)
4. Update Makefile with new build and run targets
5. Add binary names to `.gitignore`

### Building and Testing
```bash
make build-examples    # Build all examples
make test-examples     # Test all examples
make clean-examples    # Clean binaries
```

## Documentation Languages

Each example includes documentation in three languages:
- **English**: `README.md`
- **Spanish**: `README.es.md` 
- **Portuguese**: `README.pt.md`

## Next Steps

Choose the appropriate server example based on your needs:

1. **Learning/Testing**: Start with [Basic Server](basic/README.md)
2. **Production**: Use [Advanced Server](advanced/README.md)
3. **High Query Volume**: Consider [Cache Server](advanced/cache-server/README.md)
4. **Security Critical**: Use [Validation Server](advanced/validation-server/README.md)
5. **Enterprise**: Deploy [Full-Featured Server](advanced/full-featured-server/README.md)