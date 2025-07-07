package client

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Conn struct {
	deviceID string
	conn     *amqp.Connection
	config   *DSNConfig
}

func (c *Conn) logf(format string, args ...interface{}) {
	if c.config != nil && c.config.Debug {
		log.Printf("[client debug] "+format, args...)
	}
}

func (c *Conn) Prepare(query string) (driver.Stmt, error) {
	return nil, fmt.Errorf("Prepare not implemented")
}

func (c *Conn) Close() error {
	c.logf("Closing connection to RabbitMQ")
	return c.conn.Close()
}

func (c *Conn) Begin() (driver.Tx, error) {
	return nil, fmt.Errorf("transactions not supported")
}

func (c *Conn) Query(query string, args []driver.Value) (driver.Rows, error) {
	startTotal := time.Now()
	c.logf("Executing query: %s", query)

	ctx, cancel := context.WithTimeout(context.Background(), c.config.Timeout)
	defer cancel()

	named := make([]driver.NamedValue, len(args))
	for i, v := range args {
		named[i] = driver.NamedValue{Ordinal: i + 1, Value: v}
	}

	rows, err := c.queryRPC(ctx, query, named)

	total := time.Since(startTotal)
	c.logf("total time: %v", total)

	return rows, err
}

func (c *Conn) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	startTotal := time.Now()
	c.logf("Executing query (Context): %s", query)

	rows, err := c.queryRPC(ctx, query, args)

	total := time.Since(startTotal)
	c.logf("total time (QueryContext): %v", total)

	return rows, err
}

func getOutboundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "unknown"
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}

// parseCommand detecta el tipo de comando y extrae el comando real
func parseCommand(query string) (cmdType string, actualQuery string) {
	// Detectar prefijos especiales para funciones y comandos
	if len(query) > 9 && query[:9] == "FUNCTION:" {
		return "function", query[9:]
	}
	if len(query) > 8 && query[:8] == "COMMAND:" {
		return "command", query[8:]
	}
	// Por defecto, tratar como SQL
	return "sql", query
}

func (c *Conn) queryRPC(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	ch, err := c.conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to create RabbitMQ channel: %v", err)
	}
	defer ch.Close()
	c.logf("RabbitMQ channel opened")

	replyQueue, err := ch.QueueDeclare("", false, true, true, false, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to declare reply queue: %v", err)
	}
	c.logf("Reply queue declared: %s", replyQueue.Name)

	corrID := fmt.Sprintf("%d", time.Now().UnixNano())

	// Detectar el tipo de comando y extraer el comando real
	cmdType, actualQuery := parseCommand(query)
	c.logf("Detected command type: %s, actual query: %s", cmdType, actualQuery)

	req := map[string]interface{}{
		"type":     cmdType,
		"deviceID": c.deviceID,
		"query":    actualQuery,
		"params":   argsToSlice(args),
		"clientIP": getOutboundIP(),
	}

	body, _ := json.Marshal(req)

	startRT := time.Now()
	c.logf("Publishing query to device queue '%s'", c.deviceID)

	err = ch.PublishWithContext(ctx, "", c.deviceID, false, false, amqp.Publishing{
		ContentType:   "application/json",
		CorrelationId: corrID,
		ReplyTo:       replyQueue.Name,
		Body:          body,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to publish query to device queue '%s': %v\nPlease check:\n- Server is running\n- Device ID '%s' is correct\n- Queue exists", c.deviceID, err, c.deviceID)
	}
	c.logf("Query published, waiting for response...")

	msgs, err := ch.Consume(replyQueue.Name, "", true, true, false, false, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to consume from reply queue: %v", err)
	}

	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("timeout (%v) waiting for device response from '%s'\nPlease check:\n- Server is running and responding\n- Device ID '%s' is correct\n- Database is accessible", c.config.Timeout, c.deviceID, c.deviceID)
	case msg := <-msgs:
		rt := time.Since(startRT)
		c.logf("RabbitMQ roundtrip time: %v", rt)

		if msg.CorrelationId != corrID {
			return nil, fmt.Errorf("correlation id mismatch: expected %s, got %s", corrID, msg.CorrelationId)
		}
		var resp RPCResponse
		if err := json.Unmarshal(msg.Body, &resp); err != nil {
			return nil, fmt.Errorf("failed to parse server response: %v", err)
		}
		if resp.Error != "" {
			return nil, fmt.Errorf("server error: %s", resp.Error)
		}
		c.logf("Response received with %d rows", len(resp.Rows))
		return &Rows{columns: resp.Columns, rows: resp.Rows}, nil
	}
}

func argsToSlice(args []driver.NamedValue) []interface{} {
	var out []interface{}
	for _, a := range args {
		out = append(out, a.Value)
	}
	return out
}
