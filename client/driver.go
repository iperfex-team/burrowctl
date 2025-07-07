package client

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"net/url"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func init() {
	sql.Register("rabbitsql", &Driver{})
}

type Driver struct{}

func (d *Driver) Open(dsn string) (driver.Conn, error) {
	conf, err := parseDSN(dsn)
	if err != nil {
		return nil, fmt.Errorf("DSN parsing failed: %v", err)
	}

	// Intentar conectar a RabbitMQ
	conn, err := amqp.Dial(conf.AMQPURL)
	if err != nil {
		return nil, fmt.Errorf("RabbitMQ connection failed to '%s': %v\nPlease check:\n- RabbitMQ server is running\n- Credentials are correct\n- Network connectivity", conf.AMQPURL, err)
	}

	return &Conn{
		deviceID: conf.DeviceID,
		conn:     conn,
		config:   conf,
	}, nil
}

type DSNConfig struct {
	DeviceID string
	AMQPURL  string
	Timeout  time.Duration
}

func parseDSN(dsn string) (*DSNConfig, error) {
	// Parsear el DSN como query parameters
	u, err := url.Parse("?" + dsn)
	if err != nil {
		return nil, fmt.Errorf("invalid DSN format: %v", err)
	}

	values := u.Query()

	// Verificar parámetros requeridos
	deviceID := values.Get("deviceID")
	if deviceID == "" {
		return nil, fmt.Errorf("missing required parameter 'deviceID' in DSN")
	}

	amqpURI := values.Get("amqp_uri")
	if amqpURI == "" {
		return nil, fmt.Errorf("missing required parameter 'amqp_uri' in DSN")
	}

	// Validar que amqp_uri tenga el formato correcto
	if len(amqpURI) < 7 || amqpURI[:7] != "amqp://" {
		return nil, fmt.Errorf("invalid amqp_uri format: must start with 'amqp://'")
	}

	// Parsear timeout (opcional)
	timeoutStr := values.Get("timeout")
	timeout := 5 * time.Second // valor por defecto
	if timeoutStr != "" {
		parsedTimeout, err := time.ParseDuration(timeoutStr)
		if err != nil {
			return nil, fmt.Errorf("invalid timeout format '%s': %v (example: '5s', '30s', '1m')", timeoutStr, err)
		}
		timeout = parsedTimeout
	}

	conf := &DSNConfig{
		DeviceID: deviceID,
		AMQPURL:  amqpURI,
		Timeout:  timeout,
	}

	return conf, nil
}
