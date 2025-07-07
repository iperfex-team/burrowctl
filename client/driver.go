package client

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
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
		return nil, err
	}

	conn, err := amqp.Dial(conf.AMQPURL)
	if err != nil {
		return nil, err
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
	conf := &DSNConfig{}
	var timeoutStr string
	_, err := fmt.Sscanf(dsn, "deviceID=%s&amqp_uri=%s&timeout=%s", &conf.DeviceID, &conf.AMQPURL, &timeoutStr)
	if err != nil {
		return nil, fmt.Errorf("invalid DSN: %v", err)
	}
	conf.Timeout, _ = time.ParseDuration(timeoutStr)
	if conf.Timeout == 0 {
		conf.Timeout = 5 * time.Second
	}
	return conf, nil
}
