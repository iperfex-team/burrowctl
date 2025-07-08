# syntax=docker/dockerfile:1

# -------------------------
# Etapa de construcción
# -------------------------
FROM golang:1.22-bookworm AS builder

WORKDIR /app

# Copiar archivos de módulos primero (mejor cache)
COPY go.mod go.sum ./
RUN go mod download

# Copiar todo el código fuente
COPY . .

# Compilar binario de servidor avanzado de forma estática
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o burrowctl-server-advanced ./examples/server/advanced/advanced_server_example.go

# -------------------------
# Imagen final mínima
# -------------------------
FROM debian:bookworm-slim 

RUN apt-get update && apt-get install -y procps net-tools curl vim tzdata && rm -rf /var/lib/apt/lists/*

# Copiar binario desde la etapa anterior
COPY --from=builder /app/burrowctl-server-advanced /app/server

# Establecer zona horaria
ENV TZ=America/Argentina/Buenos_Aires

# Crear directorio de trabajo
WORKDIR /app

# Hacer el binario ejecutable
RUN chmod +x /app/server

# El ENTRYPOINT será sobrescrito por docker-compose con los parámetros específicos
ENTRYPOINT ["/app/server"]