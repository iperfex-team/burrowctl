# syntax=docker/dockerfile:1

# -------------------------
# Etapa de construcción
# -------------------------
FROM golang:1.22 AS builder

WORKDIR /app

# Copiar archivos de módulos primero (mejor cache)
COPY go.mod go.sum ./
RUN go mod download

# Copiar todo el código fuente
COPY . .

# Compilar binario de forma estática
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o burrowctl-server ./examples/server/server_example.go

# -------------------------
# Imagen final mínima
# -------------------------
FROM scratch

# Copiar binario desde la etapa anterior
COPY --from=builder /app/burrowctl-server /burrowctl-server

# Configurar zona horaria por defecto (opcional, comenta si no es necesario)
# ENV TZ=America/Argentina/Buenos_Aires

ENTRYPOINT ["/burrowctl-server"]