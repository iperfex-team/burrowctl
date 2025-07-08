# 🐇 burrowctl

<div align="center">
  <h3>Execução Remota de SQL e Controle de Dispositivos via RabbitMQ</h3>
  <p>
    <strong>burrowctl</strong> é uma poderosa biblioteca e serviço Go que fornece uma ponte baseada em RabbitMQ para executar remotamente consultas SQL, funções personalizadas e comandos do sistema em dispositivos atrás de NAT ou firewalls.
  </p>
  <p>
    <a href="./README.md">🇺🇸 English</a> | 
    <a href="./README.es.md">🇪🇸 Español</a> | 
    <a href="./README.pt.md">🇧🇷 Português</a>
  </p>
</div>

## 🎯 O que é burrowctl?

**burrowctl** permite acesso seguro a bancos de dados remotos e controle de dispositivos sem expor conexões diretas. É perfeito para:

- 🏢 **Plataformas SaaS**: Gerenciar bancos de dados de clientes atrás de NAT/firewalls
- 🌐 **Gestão IoT**: Controlar dispositivos distribuídos com segurança
- 🔐 **Administração Remota**: Executar consultas e comandos sem SSH/acesso direto ao BD
- 📊 **Monitoramento Distribuído**: Coletar dados de múltiplas fontes remotas

## ✨ Características Principais

### 🔌 **Suporte Multi-Cliente**
- **Cliente Go**: Compatibilidade nativa com driver `database/sql`
- **Cliente Node.js/TypeScript**: API async moderna com tipagem completa
- **DSN Universal**: Mesmo formato de string de conexão para todos os clientes

### 🚀 **Três Tipos de Execução**
- **Consultas SQL**: Acesso direto ao banco de dados com binding de parâmetros
- **Funções Personalizadas**: Sistema de funções extensível com 16+ funções incorporadas
- **Comandos do Sistema**: Executar comandos do SO com acesso controlado

### 🔒 **Pronto para Empresa**
- **Transporte Seguro**: Protocolo RabbitMQ AMQP 0-9-1
- **Pool de Conexões**: Pools de conexões de banco de dados configuráveis
- **Tratamento de Erros**: Gestão abrangente de erros e debugging
- **Controle de Timeouts**: Timeouts configuráveis para consultas e comandos

### 📦 **Recursos de Produção**
- **Suporte Docker**: Ambiente de desenvolvimento completamente containerizado
- **Automação Makefile**: Automação de build, teste e deployment
- **Controle de Versão**: Versionamento semântico automático
- **Múltiplos Exemplos**: Exemplos de uso e documentação abrangente

---

## 🚀 Início Rápido

### Pré-requisitos

- **Go 1.22+** para cliente/servidor Go
- **Node.js 16+** para cliente TypeScript
- Servidor **RabbitMQ** rodando
- Banco de dados **MySQL/MariaDB**

### Instalação

```bash
git clone https://github.com/lordbasex/burrowctl.git
cd burrowctl
go mod tidy
```

### Uso Básico

#### Cliente Go (SQL)
```go
package main

import (
    "database/sql"
    "fmt"
    "log"
    _ "github.com/lordbasex/burrowctl/client"
)

func main() {
    dsn := "deviceID=meu-dispositivo&amqp_uri=amqp://user:pass@localhost:5672/&timeout=10s&debug=true"
    
    db, err := sql.Open("rabbitsql", dsn)
    if err != nil {
        log.Fatal("Falha na conexão:", err)
    }
    defer db.Close()
    
    rows, err := db.Query("SELECT id, name FROM users WHERE active = ?", true)
    if err != nil {
        log.Fatal("Falha na consulta:", err)
    }
    defer rows.Close()
    
    for rows.Next() {
        var id int
        var name string
        rows.Scan(&id, &name)
        fmt.Printf("ID: %d, Nome: %s\n", id, name)
    }
}
```

#### Cliente Node.js/TypeScript
```typescript
import { createClient } from 'burrowctl-client-nodejs';

const client = await createClient(
  'deviceID=meu-dispositivo&amqp_uri=amqp://user:pass@localhost:5672/&timeout=10s'
);

const rows = await client.query('SELECT * FROM users WHERE active = ?', [true]);
console.log('Resultados:', rows.getRows());
console.log('Colunas:', rows.getColumns());

await client.close();
```

#### Configuração do Servidor
```go
package main

import (
    "context"
    "log"
    "time"
    "github.com/lordbasex/burrowctl/server"
)

func main() {
    pool := &server.PoolConfig{
        MaxIdleConns:    10,
        MaxOpenConns:    20,
        ConnMaxLifetime: 5 * time.Minute,
    }
    
    handler := server.NewHandler(
        "meu-dispositivo",                              // ID do dispositivo
        "amqp://user:pass@localhost:5672/",            // URI RabbitMQ
        "user:pass@tcp(localhost:3306)/dbname",        // DSN MySQL
        "open",                                        // Modo de conexão
        pool,                                          // Configuração do pool
    )
    
    // Registrar funções personalizadas
    handler.RegisterFunction("obterInfoSistema", obterInfoSistema)
    handler.RegisterFunction("processarDados", processarDados)
    
    ctx := context.Background()
    log.Println("Iniciando servidor burrowctl...")
    if err := handler.Start(ctx); err != nil {
        log.Fatal("Falha no servidor:", err)
    }
}
```

---

## 📖 Tipos de Execução

### 1. 🗃️ Consultas SQL (`sql`)

Executa consultas SQL diretas com binding de parâmetros e suporte completo a transações.

```go
// Cliente Go
rows, err := db.Query("SELECT * FROM products WHERE category = ? AND price > ?", "electronics", 100)

// Cliente Node.js
const rows = await client.query("SELECT * FROM products WHERE category = ? AND price > ?", ["electronics", 100]);
```

**Características:**
- Binding de parâmetros para segurança
- Suporte a transações
- Pool de conexões
- Tratamento de resultados com tipagem segura

### 2. ⚙️ Funções Personalizadas (`function`)

Executa funções do lado do servidor com parâmetros tipados e múltiplos valores de retorno.

```go
// Cliente Go - usando solicitação de função JSON
funcReq := FunctionRequest{
    Name: "calcularImposto",
    Params: []FunctionParam{
        {Type: "float64", Value: 100.0},
        {Type: "string", Value: "BR"},
    },
}
jsonData, _ := json.Marshal(funcReq)
rows, err := db.Query("FUNCTION:" + string(jsonData))
```

```typescript
// Cliente Node.js
const result = await client.query('FUNCTION:{"name":"calcularImposto","params":[{"type":"float64","value":100.0},{"type":"string","value":"BR"}]}');
```

**Funções Incorporadas (16+):**
- `lengthOfString`: Obter comprimento da string
- `addIntegers`: Somar dois inteiros
- `getCurrentTimestamp`: Obter timestamp atual
- `generateUUID`: Gerar UUID
- `encodeBase64`: Codificação Base64
- `decodeBase64`: Decodificação Base64
- `parseJSON`: Parsear string JSON
- `formatJSON`: Formatar JSON com indentação
- `getSystemInfo`: Obter informações do sistema
- `listFiles`: Listar conteúdo do diretório
- `readFile`: Ler conteúdo do arquivo
- `writeFile`: Escrever conteúdo do arquivo
- `calculateHash`: Calcular hash SHA256
- `validateEmail`: Validar endereço de email
- `generateRandomString`: Gerar string aleatória
- `convertTimezone`: Converter fuso horário

### 3. 🖥️ Comandos do Sistema (`command`)

Executa comandos do sistema com acesso controlado e gestão de timeouts.

```go
// Cliente Go
rows, err := db.Query("COMMAND:ps aux | grep mysql")
rows, err := db.Query("COMMAND:df -h")
rows, err := db.Query("COMMAND:systemctl status nginx")
```

```typescript
// Cliente Node.js
const result = await client.query('COMMAND:ps aux | grep mysql');
const diskUsage = await client.query('COMMAND:df -h');
```

**Características:**
- Captura de stdout/stderr
- Timeouts configuráveis
- Preservação de saída linha por linha
- Tratamento de códigos de erro

---

## 🔧 Configuração

### Formato DSN
```
deviceID=<id-dispositivo>&amqp_uri=<url-rabbitmq>&timeout=<timeout>&debug=<boolean>
```

**Parâmetros:**
- `deviceID`: Identificador único do dispositivo (tipicamente hash SHA256)
- `amqp_uri`: URL de conexão RabbitMQ
- `timeout`: Timeout da consulta (ex., `5s`, `30s`, `2m`)
- `debug`: Habilitar logging de debug (`true`/`false`)

### Configuração do Pool de Conexões
```go
pool := &server.PoolConfig{
    MaxIdleConns:    10,          // Máximo de conexões ociosas
    MaxOpenConns:    20,          // Máximo de conexões abertas
    ConnMaxLifetime: 5 * time.Minute, // Tempo de vida da conexão
}
```

### Modos de Conexão
- **`open`**: Mantém pool de conexões (padrão, melhor performance)
- **`close`**: Abre/fecha conexões por consulta (mais seguro, mais lento)

---

## 🛠️ Desenvolvimento

### Configuração Rápida de Desenvolvimento
```bash
# Clonar e configurar
git clone https://github.com/lordbasex/burrowctl.git
cd burrowctl

# Iniciar ambiente de desenvolvimento (Docker)
cd examples/server
docker-compose up -d

# Construir projeto
make build

# Executar exemplos
make run-server-example
make run-sql-example
make run-function-example
make run-command-example
```

### Comandos Make Disponíveis
```bash
make help                    # Mostrar todos os comandos disponíveis
make build                   # Construir todos os componentes
make test                    # Executar testes
make clean                   # Limpar artefatos de build
make docker-up              # Iniciar ambiente Docker
make docker-down            # Parar ambiente Docker
make run-server-example     # Executar exemplo do servidor
make run-sql-example        # Executar exemplo de cliente SQL
make run-function-example   # Executar exemplo de cliente de funções
make run-command-example    # Executar exemplo de cliente de comandos
```

---

## 🏗️ Arquitetura

### Componentes do Sistema

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Cliente Go    │    │   Cliente       │    │   Futuros       │
│   (database/sql)│    │   Node.js       │    │   Clientes      │
└─────────┬───────┘    └─────────┬───────┘    └─────────┬───────┘
          │                      │                      │
          └──────────────────────┼──────────────────────┘
                                 │
                    ┌─────────────▼──────────────┐
                    │       RabbitMQ             │
                    │    (AMQP 0-9-1)           │
                    └─────────────┬──────────────┘
                                 │
                ┌─────────────────▼──────────────────┐
                │       Servidor burrowctl           │
                │  ┌─────────────────────────────┐   │
                │  │    Motor SQL               │   │
                │  │    Registro de Funções     │   │
                │  │    Executor de Comandos    │   │
                │  └─────────────────────────────┘   │
                └─────────────────┬──────────────────┘
                                 │
                    ┌─────────────▼──────────────┐
                    │       MySQL/MariaDB        │
                    │       Sistema de Arquivos  │
                    │       Comandos do Sistema  │
                    └────────────────────────────┘
```

### Fluxo de Mensagens

1. **Cliente**: Envia solicitação para fila RabbitMQ específica do dispositivo
2. **RabbitMQ**: Roteia mensagem para fila apropriada do dispositivo
3. **Servidor**: Processa solicitação baseada no tipo (`sql`, `function`, `command`)
4. **Execução**: Executa contra banco de dados, registro de funções, ou sistema
5. **Resposta**: Retorna resultados via fila de resposta RabbitMQ
6. **Cliente**: Recebe e processa resposta

---

## 📁 Estrutura do Projeto

```
burrowctl/
├── client/                 # Cliente Go (driver database/sql)
│   ├── driver.go          # Implementação do driver SQL
│   ├── conn.go            # Gestão de conexões
│   ├── rows.go            # Tratamento de resultados
│   └── rpc.go             # Cliente RPC RabbitMQ
├── server/                 # Biblioteca do servidor core
│   └── server.go          # Implementação do servidor
├── client-nodejs/          # Cliente Node.js/TypeScript
│   ├── src/               # Código fonte TypeScript
│   ├── dist/              # JavaScript compilado
│   └── package.json       # Configuração do pacote NPM
├── examples/              # Exemplos de uso
│   ├── client/            # Exemplos de cliente
│   │   ├── sql-example/   # Uso de SQL
│   │   ├── function-example/ # Uso de funções
│   │   └── command-example/  # Uso de comandos
│   └── server/            # Exemplos de servidor
│       ├── server_example.go # Configuração completa do servidor
│       └── docker-compose.yml # Ambiente de desenvolvimento
├── Makefile              # Automação de build
├── go.mod               # Dependências do módulo Go
└── version.txt          # Informações de versão
```

---

## 🔐 Considerações de Segurança

### Melhores Práticas

1. **Usar Credenciais Fortes**: Sempre usar senhas fortes para RabbitMQ e banco de dados
2. **Habilitar TLS**: Usar TLS/SSL para conexões RabbitMQ em produção
3. **Limitar Acesso a Funções**: Registrar apenas funções necessárias no servidor
4. **Restrições de Comandos**: Implementar lista branca de comandos para segurança
5. **Isolamento de Rede**: Usar VPNs ou redes privadas quando possível
6. **Monitoramento**: Implementar logging e monitoramento para auditoria de segurança

### Configuração de Produção

```go
// Configuração do servidor para produção
handler := server.NewHandler(
    os.Getenv("DEVICE_ID"),
    os.Getenv("AMQP_URI"),     // Usar TLS: amqps://user:pass@host:5671/
    os.Getenv("MYSQL_DSN"),    // Usar SSL: ?tls=true
    "open",
    &server.PoolConfig{
        MaxIdleConns:    5,
        MaxOpenConns:    10,
        ConnMaxLifetime: 2 * time.Minute,
    },
)
```

---

## 🤝 Contribuir

Recebemos contribuições! Por favor, veja nosso [Guia de Contribuição](CONTRIBUTING.md) para detalhes.

### Configuração de Desenvolvimento

1. Fazer fork do repositório
2. Criar uma branch de feature: `git checkout -b feature/recurso-incrivel`
3. Fazer alterações
4. Adicionar testes para nova funcionalidade
5. Executar testes: `make test`
6. Fazer commit das alterações: `git commit -m 'Adicionar recurso incrível'`
7. Fazer push para a branch: `git push origin feature/recurso-incrivel`
8. Abrir um Pull Request

---

## 📜 Licença

Este projeto é licenciado sob a Licença MIT - veja o arquivo [LICENSE](LICENSE) para detalhes.

---

## 🆘 Suporte

- **Documentação**: [Documentação completa](./examples/)
- **Exemplos**: [Exemplos de uso](./examples/client/)
- **Issues**: [GitHub Issues](https://github.com/lordbasex/burrowctl/issues)
- **Discussões**: [GitHub Discussions](https://github.com/lordbasex/burrowctl/discussions)

---

## 🙏 Agradecimentos

- [RabbitMQ](https://www.rabbitmq.com/) pelo excelente message broker
- [Go SQL Driver](https://github.com/go-sql-driver/mysql) pela conectividade MySQL
- [AMQP 0-9-1 Go Client](https://github.com/rabbitmq/amqp091-go) pela integração RabbitMQ
- As comunidades Go e Node.js por seus excelentes ecossistemas

---

<div align="center">
  <p>Feito com ❤️ pela equipe burrowctl</p>
  <p>
    <a href="https://github.com/lordbasex/burrowctl/stargazers">⭐ Dar estrela a este projeto</a> | 
    <a href="https://github.com/lordbasex/burrowctl/issues">🐛 Reportar Bug</a> | 
    <a href="https://github.com/lordbasex/burrowctl/issues">💡 Solicitar Recurso</a>
  </p>
</div>