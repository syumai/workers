package d1

import (
	"context"
	"database/sql/driver"
)

type Connector struct {
	name string
}

var (
	_ driver.Connector = (*Connector)(nil)
)

// Connect returns Conn object of D1.
// This method doesn't check DB existence, so this function never return errors.
func (c *Connector) Connect(context.Context) (driver.Conn, error) {
	return &Conn{dbName: c.name}, nil
}

func (c *Connector) Driver() driver.Driver {
	return &Driver{}
}
