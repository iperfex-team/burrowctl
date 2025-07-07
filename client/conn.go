package client

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Conn struct {
	deviceID string
	conn     *amqp.Connection
	config   *DSNConfig
}

func (c *Conn) Prepare(query string) (driver.Stmt, error) {
	return nil, errors.New("Prepare not implemented")
}

func (c *Conn) Close() error {
	return c.conn.Close()
}

func (c *Conn) Begin() (driver.Tx, error) {
	return nil, errors.New("transactions not supported")
}

func (c *Conn) Query(query string, args []driver.Value) (driver.Rows, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.config.Timeout)
	defer cancel()
	named := make([]driver.NamedValue, len(args))
	for i, v := range args {
		named[i] = driver.NamedValue{Ordinal: i + 1, Value: v}
	}
	return c.queryRPC(ctx, query, named)
}

func (c *Conn) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	return c.queryRPC(ctx, query, args)
}

func (c *Conn) queryRPC(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	ch, err := c.conn.Channel()
	if err != nil {
		return nil, err
	}
	defer ch.Close()

	replyQueue, err := ch.QueueDeclare("", false, true, true, false, nil)
	if err != nil {
		return nil, err
	}

	corrID := fmt.Sprintf("%d", time.Now().UnixNano())

	req := map[string]interface{}{
		"type":     "sql",
		"deviceID": c.deviceID,
		"query":    query,
		"params":   argsToSlice(args),
	}

	body, _ := json.Marshal(req)

	err = ch.PublishWithContext(ctx, "", c.deviceID, false, false, amqp.Publishing{
		ContentType:   "application/json",
		CorrelationId: corrID,
		ReplyTo:       replyQueue.Name,
		Body:          body,
	})
	if err != nil {
		return nil, err
	}

	msgs, err := ch.Consume(replyQueue.Name, "", true, true, false, false, nil)
	if err != nil {
		return nil, err
	}

	select {
	case <-ctx.Done():
		return nil, errors.New("timeout waiting for device response")
	case msg := <-msgs:
		if msg.CorrelationId != corrID {
			return nil, errors.New("correlation id mismatch")
		}
		var resp RPCResponse
		if err := json.Unmarshal(msg.Body, &resp); err != nil {
			return nil, err
		}
		if resp.Error != "" {
			return nil, errors.New(resp.Error)
		}
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
