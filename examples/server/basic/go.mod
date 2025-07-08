module server

go 1.23.2

require github.com/lordbasex/burrowctl v0.0.0-00010101000000-000000000000

require (
	github.com/go-sql-driver/mysql v1.7.1 // indirect
	github.com/rabbitmq/amqp091-go v1.10.0 // indirect
)

replace github.com/lordbasex/burrowctl => ../../..
