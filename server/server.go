package server

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os/exec"
	"reflect"
	"strconv"
	"strings"
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

// Estructuras para manejo de funciones
type FunctionParam struct {
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}

type FunctionRequest struct {
	Name   string          `json:"name"`
	Params []FunctionParam `json:"params"`
}

// Struct para ejemplo (copiado de demo-func.go)
type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
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
	ClientIP string        `json:"clientIP"`
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

	log.Printf("[server] received ip=%s type=%s query=%s", req.ClientIP, req.Type, req.Query)

	switch req.Type {
	case "sql":
		h.handleSQL(ch, msg, req)

	case "function":
		h.handleFunction(ch, msg, req)

	case "command":
		h.handleCommand(ch, msg, req)

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

	// Obtener tipos de columnas para mejor manejo de datos
	colTypes, err := rows.ColumnTypes()
	if err != nil {
		h.respond(ch, msg.ReplyTo, msg.CorrelationId, RPCResponse{Error: err.Error()})
		return
	}

	var data [][]interface{}
	for rows.Next() {
		// Crear scan destinos basados en tipos de columnas
		scanDest := make([]interface{}, len(cols))
		for i := range scanDest {
			scanDest[i] = new(interface{})
		}

		if err := rows.Scan(scanDest...); err != nil {
			h.respond(ch, msg.ReplyTo, msg.CorrelationId, RPCResponse{Error: err.Error()})
			return
		}

		// Convertir y limpiar tipos de datos
		row := make([]interface{}, len(cols))
		for i, val := range scanDest {
			v := *(val.(*interface{}))
			row[i] = h.convertDatabaseValue(v, colTypes[i])
		}
		data = append(data, row)
	}

	h.respond(ch, msg.ReplyTo, msg.CorrelationId, RPCResponse{
		Columns: cols,
		Rows:    data,
	})
}

// Función auxiliar para convertir valores de base de datos a tipos JSON apropiados
func (h *Handler) convertDatabaseValue(val interface{}, colType *sql.ColumnType) interface{} {
	if val == nil {
		return nil
	}

	switch v := val.(type) {
	case []byte:
		// Determinar si debería ser string o número basado en el tipo de columna
		dbType := colType.DatabaseTypeName()
		switch dbType {
		case "TINYINT", "SMALLINT", "MEDIUMINT", "INT", "INTEGER", "BIGINT":
			// Intentar convertir bytes a número
			str := string(v)
			if str == "" {
				return 0
			}
			// Para números, retornar como string para que el cliente lo parsee
			return str
		case "DECIMAL", "NUMERIC", "FLOAT", "DOUBLE", "REAL":
			// Para decimales, retornar como string para que el cliente lo parsee
			return string(v)
		case "CHAR", "VARCHAR", "TEXT", "TINYTEXT", "MEDIUMTEXT", "LONGTEXT":
			// Para texto, retornar como string
			return string(v)
		default:
			// Por defecto, convertir a string
			return string(v)
		}
	case string:
		return v
	case int, int8, int16, int32, int64:
		return v
	case uint, uint8, uint16, uint32, uint64:
		return v
	case float32, float64:
		return v
	case bool:
		return v
	default:
		// Para otros tipos, convertir a string
		return fmt.Sprintf("%v", v)
	}
}

// handleCommand ejecuta un comando del sistema y devuelve su salida
func (h *Handler) handleCommand(ch *amqp.Channel, msg amqp.Delivery, req RPCRequest) {
	// Crear contexto con timeout para evitar que comandos se ejecuten indefinidamente
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	log.Printf("[server] executing command: %s", req.Query)

	// Parsear el comando y sus argumentos
	// El comando viene en req.Query, necesitamos dividirlo en comando y argumentos
	parts := strings.Fields(req.Query)
	if len(parts) == 0 {
		h.respond(ch, msg.ReplyTo, msg.CorrelationId, RPCResponse{
			Error: "empty command",
		})
		return
	}

	command := parts[0]
	args := parts[1:]

	// Crear y ejecutar el comando
	cmd := exec.CommandContext(ctx, command, args...)

	// Capturar tanto stdout como stderr
	output, err := cmd.CombinedOutput()

	if err != nil {
		// Si hay error, incluir tanto el error como la salida (si la hay)
		errorMsg := fmt.Sprintf("command failed: %v", err)
		if len(output) > 0 {
			errorMsg += fmt.Sprintf("\nOutput: %s", string(output))
		}
		h.respond(ch, msg.ReplyTo, msg.CorrelationId, RPCResponse{
			Error: errorMsg,
		})
		return
	}

	// Convertir la salida a string y dividir en líneas
	outputStr := string(output)

	// Dividir la salida en líneas
	lines := strings.Split(outputStr, "\n")

	// Preparar las filas para la respuesta
	var rows [][]interface{}

	// Agregar cada línea como una fila
	for _, line := range lines {
		// Incluir líneas vacías también, ya que pueden ser parte de la salida
		rows = append(rows, []interface{}{line})
	}

	// Si no hay salida, agregar una fila indicando que el comando se ejecutó correctamente
	if len(rows) == 0 {
		rows = append(rows, []interface{}{"(command executed successfully - no output)"})
	}

	// Responder con la salida del comando
	h.respond(ch, msg.ReplyTo, msg.CorrelationId, RPCResponse{
		Columns: []string{"output"},
		Rows:    rows,
	})

	log.Printf("[server] command executed successfully, returned %d lines", len(rows))
}

// handleFunction ejecuta una función remota y devuelve su resultado
func (h *Handler) handleFunction(ch *amqp.Channel, msg amqp.Delivery, req RPCRequest) {
	// Crear contexto con timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	log.Printf("[server] executing function: %s", req.Query)

	// Parsear la solicitud de función desde req.Query
	var funcReq FunctionRequest
	if err := json.Unmarshal([]byte(req.Query), &funcReq); err != nil {
		h.respond(ch, msg.ReplyTo, msg.CorrelationId, RPCResponse{
			Error: fmt.Sprintf("invalid function request: %v", err),
		})
		return
	}

	// Ejecutar la función
	result, err := h.executeFunction(ctx, funcReq)
	if err != nil {
		h.respond(ch, msg.ReplyTo, msg.CorrelationId, RPCResponse{
			Error: fmt.Sprintf("function execution failed: %v", err),
		})
		return
	}

	// Convertir resultado a formato de respuesta
	columns, rows := h.convertFunctionResult(result)

	h.respond(ch, msg.ReplyTo, msg.CorrelationId, RPCResponse{
		Columns: columns,
		Rows:    rows,
	})

	log.Printf("[server] function executed successfully")
}

// executeFunction ejecuta una función por nombre usando reflection
func (h *Handler) executeFunction(ctx context.Context, funcReq FunctionRequest) ([]interface{}, error) {
	// Verificar si el contexto fue cancelado
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	// Obtener la función por nombre
	funcValue := h.getFunctionByName(funcReq.Name)
	if !funcValue.IsValid() {
		return nil, fmt.Errorf("function '%s' not found", funcReq.Name)
	}

	// Preparar parámetros
	params, err := h.prepareFunctionParams(funcReq.Params, funcValue.Type())
	if err != nil {
		return nil, fmt.Errorf("error preparing parameters: %v", err)
	}

	// Ejecutar función
	results := funcValue.Call(params)

	// Convertir resultados a []interface{}
	var output []interface{}
	for _, result := range results {
		output = append(output, result.Interface())
	}

	return output, nil
}

// getFunctionByName devuelve la función por nombre
func (h *Handler) getFunctionByName(name string) reflect.Value {
	// Mapa de funciones disponibles
	functions := map[string]interface{}{
		"returnError":       returnError,
		"returnBool":        returnBool,
		"returnInt":         returnInt,
		"returnString":      returnString,
		"returnStruct":      returnStruct,
		"returnIntArray":    returnIntArray,
		"returnStringArray": returnStringArray,
		"returnJSON":        returnJSON,
		"lengthOfString":    lengthOfString,
		"isEven":            isEven,
		"greetPerson":       greetPerson,
		"sumArray":          sumArray,
		"validateString":    validateString,
		"complexFunction":   complexFunction,
		"flagToPerson":      flagToPerson,
		"modifyJSON":        modifyJSON,
	}

	if fn, exists := functions[name]; exists {
		return reflect.ValueOf(fn)
	}

	return reflect.Value{}
}

// prepareFunctionParams convierte los parámetros al tipo correcto
func (h *Handler) prepareFunctionParams(params []FunctionParam, funcType reflect.Type) ([]reflect.Value, error) {
	if len(params) != funcType.NumIn() {
		return nil, fmt.Errorf("expected %d parameters, got %d", funcType.NumIn(), len(params))
	}

	var values []reflect.Value
	for i, param := range params {
		expectedType := funcType.In(i)
		value, err := h.convertToType(param.Value, expectedType)
		if err != nil {
			return nil, fmt.Errorf("parameter %d: %v", i, err)
		}
		values = append(values, value)
	}

	return values, nil
}

// convertToType convierte un valor al tipo especificado
func (h *Handler) convertToType(value interface{}, targetType reflect.Type) (reflect.Value, error) {
	if value == nil {
		return reflect.Zero(targetType), nil
	}

	valueType := reflect.TypeOf(value)
	if valueType == targetType {
		return reflect.ValueOf(value), nil
	}

	// Conversiones específicas
	switch targetType.Kind() {
	case reflect.String:
		return reflect.ValueOf(fmt.Sprintf("%v", value)), nil

	case reflect.Int:
		switch v := value.(type) {
		case float64:
			return reflect.ValueOf(int(v)), nil
		case string:
			if i, err := strconv.Atoi(v); err == nil {
				return reflect.ValueOf(i), nil
			}
		}

	case reflect.Bool:
		switch v := value.(type) {
		case string:
			if b, err := strconv.ParseBool(v); err == nil {
				return reflect.ValueOf(b), nil
			}
		}

	case reflect.Slice:
		if valueType.Kind() == reflect.Slice {
			// Convertir slice a slice del tipo correcto
			sourceSlice := reflect.ValueOf(value)
			targetSlice := reflect.MakeSlice(targetType, sourceSlice.Len(), sourceSlice.Len())
			for i := 0; i < sourceSlice.Len(); i++ {
				convertedValue, err := h.convertToType(sourceSlice.Index(i).Interface(), targetType.Elem())
				if err != nil {
					return reflect.Value{}, err
				}
				targetSlice.Index(i).Set(convertedValue)
			}
			return targetSlice, nil
		}

	case reflect.Struct:
		if targetType == reflect.TypeOf(Person{}) {
			// Convertir a Person struct
			if jsonData, err := json.Marshal(value); err == nil {
				var person Person
				if json.Unmarshal(jsonData, &person) == nil {
					return reflect.ValueOf(person), nil
				}
			}
		}
	}

	return reflect.Value{}, fmt.Errorf("cannot convert %v to %v", valueType, targetType)
}

// convertFunctionResult convierte el resultado de la función a formato de respuesta
func (h *Handler) convertFunctionResult(results []interface{}) ([]string, [][]interface{}) {
	if len(results) == 0 {
		return []string{"result"}, [][]interface{}{{"no output"}}
	}

	var columns []string
	var rows [][]interface{}

	if len(results) == 1 {
		// Un solo resultado
		result := results[0]
		if err, ok := result.(error); ok {
			if err != nil {
				columns = []string{"error"}
				rows = [][]interface{}{{err.Error()}}
			} else {
				columns = []string{"result"}
				rows = [][]interface{}{{"success"}}
			}
		} else {
			columns = []string{"result"}
			rows = [][]interface{}{{h.formatResult(result)}}
		}
	} else {
		// Múltiples resultados
		for i := range results {
			columns = append(columns, fmt.Sprintf("result_%d", i+1))
		}

		var row []interface{}
		for _, res := range results {
			if err, ok := res.(error); ok {
				if err != nil {
					row = append(row, err.Error())
				} else {
					row = append(row, "success")
				}
			} else {
				row = append(row, h.formatResult(res))
			}
		}
		rows = [][]interface{}{row}
	}

	return columns, rows
}

// formatResult formatea un resultado para mostrar
func (h *Handler) formatResult(result interface{}) interface{} {
	if result == nil {
		return "null"
	}

	switch v := result.(type) {
	case []int:
		return fmt.Sprintf("%v", v)
	case []string:
		return fmt.Sprintf("%v", v)
	case Person:
		if jsonData, err := json.Marshal(v); err == nil {
			return string(jsonData)
		}
		return fmt.Sprintf("%+v", v)
	default:
		return result
	}
}

// Funciones de ejemplo (copiadas de demo-func.go)
func returnError() error {
	return errors.New("algo salió mal")
}

func returnBool() bool {
	return true
}

func returnInt() int {
	return 42
}

func returnString() string {
	return "Hola mundo"
}

func returnStruct() Person {
	return Person{Name: "Juan", Age: 30}
}

func returnIntArray() []int {
	return []int{1, 2, 3, 4, 5}
}

func returnStringArray() []string {
	return []string{"uno", "dos", "tres"}
}

func returnJSON() string {
	p := Person{Name: "Ana", Age: 25}
	data, _ := json.Marshal(p)
	return string(data)
}

func lengthOfString(s string) int {
	return len(s)
}

func isEven(n int) bool {
	return n%2 == 0
}

func greetPerson(p Person) string {
	return fmt.Sprintf("Hola, %s. Tienes %d años.", p.Name, p.Age)
}

func sumArray(arr []int) int {
	sum := 0
	for _, n := range arr {
		sum += n
	}
	return sum
}

func validateString(s string) error {
	if s == "" {
		return errors.New("cadena vacía")
	}
	return nil
}

func complexFunction(s string, n int) (string, int, error) {
	if s == "" {
		return "", 0, errors.New("string vacío")
	}
	return s, n * 2, nil
}

func flagToPerson(flag bool) Person {
	if flag {
		return Person{Name: "Verdadero", Age: 1}
	}
	return Person{Name: "Falso", Age: 0}
}

func modifyJSON(jsonStr string) (string, error) {
	var p Person
	err := json.Unmarshal([]byte(jsonStr), &p)
	if err != nil {
		return "", err
	}
	p.Age += 1
	data, err := json.Marshal(p)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (h *Handler) respond(ch *amqp.Channel, replyTo, corrID string, resp RPCResponse) {
	body, _ := json.Marshal(resp)
	ch.PublishWithContext(context.Background(), "", replyTo, false, false, amqp.Publishing{
		ContentType:   "application/json",
		CorrelationId: corrID,
		Body:          body,
	})
}
