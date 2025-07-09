module validation-server

go 1.22.0

toolchain go1.23.2

replace github.com/lordbasex/burrowctl => ../../../../

require github.com/lordbasex/burrowctl v0.0.0-00010101000000-000000000000

require (
	github.com/go-sql-driver/mysql v1.7.1 // indirect
	github.com/rabbitmq/amqp091-go v1.10.0 // indirect
)
