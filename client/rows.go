package client

import (
	"database/sql/driver"
	"errors"
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
		return errors.New("EOF")
	}
	for i, val := range r.rows[r.pos] {
		dest[i] = val
	}
	r.pos++
	return nil
}

func (r *Rows) Close() error {
	return nil
}
