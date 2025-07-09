# Exemplo de Servidor Básico

Uma implementação simples do servidor burrowctl que demonstra a funcionalidade básica.

## Características

- Manuseio básico de mensagens AMQP
- Conexão com banco de dados MySQL
- Execução simples de comandos
- Padrão básico de requisição/resposta

## Uso

### Execução direta
```bash
go run main.go
```

### Usando Makefile
```bash
make run-server-example
```

### Docker
```bash
make docker-up
```

## Configuração

O servidor usa valores de configuração codificados para simplicidade:

- **ID do Dispositivo**: `my-device`
- **URL AMQP**: `amqp://burrowuser:burrowpass123@localhost:5672/`
- **DSN MySQL**: `burrowuser:burrowpass123@tcp(localhost:3306)/burrowdb`

## Dependências

- Servidor RabbitMQ (porta 5672)
- Servidor MariaDB/MySQL (porta 3306)
- Go 1.22 ou superior

## Primeiros Passos

1. Iniciar os serviços necessários:
   ```bash
   make docker-up
   ```

2. Executar o servidor:
   ```bash
   make run-server-example
   ```

3. Testar com um cliente:
   ```bash
   cd ../../client/command-example
   go run main.go "ls -la"
   ```

## Arquitetura

Este servidor básico fornece:

- **Integração de Fila de Mensagens**: Conecta ao RabbitMQ para receber comandos
- **Conexão com Banco de Dados**: Usa MySQL para armazenamento de dados
- **Processamento de Comandos**: Executa comandos recebidos e retorna resultados
- **Tratamento de Erros**: Tratamento básico de erros e logging

## Próximos Passos

Para recursos mais avançados, consulte:
- [Servidor Avançado](../advanced/README.pt.md)
- [Servidor de Cache](../advanced/cache-server/README.pt.md)
- [Servidor de Validação](../advanced/validation-server/README.pt.md)
- [Servidor Completo](../advanced/full-featured-server/README.pt.md)