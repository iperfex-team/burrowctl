package client

import (
	"database/sql/driver"
	"fmt"
	"io"
	"strconv"
)

type Rows struct {
	columns []string
	rows    [][]interface{}
	pos     int
}

func (r *Rows) Columns() []string {
	return r.columns
}

func (r *Rows) Next(dest []driver.Value) error {
	if r.pos >= len(r.rows) {
		return io.EOF
	}

	for i, val := range r.rows[r.pos] {
		dest[i] = r.convertValue(val)
	}
	r.pos++
	return nil
}

// Función auxiliar para convertir valores del servidor a tipos apropiados
func (r *Rows) convertValue(val interface{}) driver.Value {
	if val == nil {
		return nil
	}

	switch v := val.(type) {
	case string:
		// Intentar convertir strings que representan números
		if intVal, err := strconv.ParseInt(v, 10, 64); err == nil {
			return intVal
		}
		if floatVal, err := strconv.ParseFloat(v, 64); err == nil {
			return floatVal
		}
		return v
	case float64:
		// JSON unmarshaling siempre devuelve float64 para números
		// Si es un entero, convertirlo
		if v == float64(int64(v)) {
			return int64(v)
		}
		return v
	case bool:
		return v
	default:
		// Para otros tipos, convertir a string
		return fmt.Sprintf("%v", v)
	}
}

func (r *Rows) Close() error {
	return nil
}
