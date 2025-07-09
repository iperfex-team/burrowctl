# Servidor Avançado

Implementação aprimorada do servidor burrowctl com recursos empresariais para ambientes de alto desempenho.

## Características

- **Pool de Workers**: Processamento concorrente configurável
- **Limitação de Taxa**: Limitação por IP com suporte a rajadas
- **Pool de Conexões**: Gestão otimizada de conexões DB
- **Monitoramento de Desempenho**: Métricas em tempo real
- **Desligamento Elegante**: Encerramento limpo com drenagem de requests

## Uso

### Execução direta
```bash
go run main.go
```

### Usando Makefile
```bash
make run-server-advanced
```

### Docker
```bash
make docker-up-advanced
```

## Configuração

### Opções de Linha de Comando

#### Opções de Desempenho
- `-workers=20`: Número de goroutines worker (padrão: 10)
- `-queue-size=500`: Tamanho da fila de worker (padrão: 100)
- `-rate-limit=50`: Limite de taxa por IP (req/s) (padrão: 10)
- `-burst-size=100`: Tamanho da rajada (padrão: 20)

#### Opções de Banco de Dados
- `-pool-idle=20`: Máximo de conexões idle (padrão: 5)
- `-pool-open=50`: Máximo de conexões abertas (padrão: 15)

### Configurações de Exemplo

#### Modo Alto Desempenho
```bash
go run main.go -workers=50 -queue-size=1000 -rate-limit=100 -burst-size=200
```

## Características de Desempenho

### Pool de Workers
- Número configurável de workers concorrentes
- Fila limitada com proteção de overflow
- Balanceamento de carga entre workers

### Limitação de Taxa
- Limitação por IP de cliente
- Algoritmo token bucket com suporte a rajadas
- Limpeza automática de dados de limite

### Pool de Conexões
- Gestão otimizada de conexões DB
- Limites configuráveis idle/open
- Gestão de tempo de vida das conexões

## Próximos Passos

Para configurações especializadas:
- [Servidor de Cache](cache-server/README.pt.md) - Cache de resultados de consultas
- [Servidor de Validação](validation-server/README.pt.md) - Validação de segurança SQL
- [Servidor Completo](full-featured-server/README.pt.md) - Todas as características empresariais