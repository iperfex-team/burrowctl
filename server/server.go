package server

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	amqp "github.com/rabbitmq/amqp091-go"
)

type PoolConfig struct {
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime time.Duration
}

type Handler struct {
	deviceID string
	amqpURL  string
	mysqlDSN string
	conn     *amqp.Connection
	db       *sql.DB
	mode     string
	poolConf PoolConfig
}

func NewHandler(deviceID, amqpURL, mysqlDSN, mode string, poolConf *PoolConfig) *Handler {
	if mode == "" {
		mode = "open"
	}
	defaultPool := PoolConfig{
		MaxIdleConns:    10,
		MaxOpenConns:    20,
		ConnMaxLifetime: 3 * time.Minute,
	}

	if poolConf == nil {
		poolConf = &defaultPool
	} else {
		if poolConf.MaxIdleConns == 0 {
			poolConf.MaxIdleConns = defaultPool.MaxIdleConns
		}
		if poolConf.MaxOpenConns == 0 {
			poolConf.MaxOpenConns = defaultPool.MaxOpenConns
		}
		if poolConf.ConnMaxLifetime == 0 {
			poolConf.ConnMaxLifetime = defaultPool.ConnMaxLifetime
		}
	}

	return &Handler{
		deviceID: deviceID,
		amqpURL:  amqpURL,
		mysqlDSN: mysqlDSN,
		mode:     mode,
		poolConf: *poolConf,
	}
}

func (h *Handler) Start(ctx context.Context) error {
	var err error

	h.conn, err = amqp.Dial(h.amqpURL)
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}
	defer h.conn.Close()

	if h.mode == "open" {
		h.db, err = sql.Open("mysql", h.mysqlDSN)
		if err != nil {
			return fmt.Errorf("failed to connect to MySQL: %w", err)
		}

		h.db.SetMaxIdleConns(h.poolConf.MaxIdleConns)
		h.db.SetMaxOpenConns(h.poolConf.MaxOpenConns)
		h.db.SetConnMaxLifetime(h.poolConf.ConnMaxLifetime)
		defer h.db.Close()

		log.Printf("[server] Database pool initialized: idle=%d open=%d lifetime=%s",
			h.poolConf.MaxIdleConns, h.poolConf.MaxOpenConns, h.poolConf.ConnMaxLifetime)
	} else {
		log.Println("[server] Using 'close' mode: opening/closing DB connection per query")
	}

	ch, err := h.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	// Declarar la cola antes de consumir de ella
	// Esto asegura que la cola existe antes de intentar consumir
	_, err = ch.QueueDeclare(
		h.deviceID, // name - nombre de la cola (usamos el deviceID)
		false,      // durable - no persistente (se pierde si se reinicia RabbitMQ)
		false,      // delete when unused - no borrar automáticamente
		false,      // exclusive - no exclusiva
		false,      // no-wait - esperar confirmación
		nil,        // arguments - sin argumentos adicionales
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	log.Printf("[server] Queue '%s' declared successfully", h.deviceID)

	msgs, err := ch.Consume(h.deviceID, "", true, true, false, false, nil)
	if err != nil {
		return err
	}

	log.Printf("[server] Listening on queue %s", h.deviceID)

	for {
		select {
		case <-ctx.Done():
			return nil
		case msg := <-msgs:
			go h.handleMessage(ch, msg)
		}
	}
}

type RPCRequest struct {
	Type     string        `json:"type"`
	DeviceID string        `json:"deviceID"`
	Query    string        `json:"query"`
	Params   []interface{} `json:"params"`
}

type RPCResponse struct {
	Columns []string        `json:"columns"`
	Rows    [][]interface{} `json:"rows"`
	Error   string          `json:"error"`
}

func (h *Handler) handleMessage(ch *amqp.Channel, msg amqp.Delivery) {
	var req RPCRequest
	if err := json.Unmarshal(msg.Body, &req); err != nil {
		h.respond(ch, msg.ReplyTo, msg.CorrelationId, RPCResponse{Error: err.Error()})
		return
	}

	log.Printf("[server] received type=%s query=%s", req.Type, req.Query)

	switch req.Type {
	case "sql":
		h.handleSQL(ch, msg, req)

	case "function":
		h.respond(ch, msg.ReplyTo, msg.CorrelationId, RPCResponse{
			Columns: []string{"message"},
			Rows:    [][]interface{}{{"function executed (mock)"}},
		})

	case "command":
		h.respond(ch, msg.ReplyTo, msg.CorrelationId, RPCResponse{
			Columns: []string{"message"},
			Rows:    [][]interface{}{{"command executed (mock)"}},
		})

	default:
		h.respond(ch, msg.ReplyTo, msg.CorrelationId, RPCResponse{
			Error: fmt.Sprintf("unsupported type: %s", req.Type),
		})
	}
}

func (h *Handler) handleSQL(ch *amqp.Channel, msg amqp.Delivery, req RPCRequest) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var db *sql.DB
	var err error

	if h.mode == "open" {
		db = h.db
	} else {
		db, err = sql.Open("mysql", h.mysqlDSN)
		if err != nil {
			h.respond(ch, msg.ReplyTo, msg.CorrelationId, RPCResponse{Error: err.Error()})
			return
		}
		defer db.Close()
	}

	rows, err := db.QueryContext(ctx, req.Query, req.Params...)
	if err != nil {
		h.respond(ch, msg.ReplyTo, msg.CorrelationId, RPCResponse{Error: err.Error()})
		return
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		h.respond(ch, msg.ReplyTo, msg.CorrelationId, RPCResponse{Error: err.Error()})
		return
	}

	var data [][]interface{}
	for rows.Next() {
		dest := make([]interface{}, len(cols))
		ptrs := make([]interface{}, len(cols))
		for i := range dest {
			ptrs[i] = &dest[i]
		}
		if err := rows.Scan(ptrs...); err != nil {
			h.respond(ch, msg.ReplyTo, msg.CorrelationId, RPCResponse{Error: err.Error()})
			return
		}
		data = append(data, dest)
	}

	h.respond(ch, msg.ReplyTo, msg.CorrelationId, RPCResponse{
		Columns: cols,
		Rows:    data,
	})
}

func (h *Handler) respond(ch *amqp.Channel, replyTo, corrID string, resp RPCResponse) {
	body, _ := json.Marshal(resp)
	ch.PublishWithContext(context.Background(), "", replyTo, false, false, amqp.Publishing{
		ContentType:   "application/json",
		CorrelationId: corrID,
		Body:          body,
	})
}
