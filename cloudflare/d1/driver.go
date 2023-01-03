package d1

import (
	"context"
	"database/sql"
	"database/sql/driver"
)

func init() {
	sql.Register("d1", &Driver{})
}

type Driver struct{}

var (
	_ driver.Driver        = (*Driver)(nil)
	_ driver.DriverContext = (*Driver)(nil)
)

func (d *Driver) Open(name string) (driver.Conn, error) {
	c, _ := d.OpenConnector(name)
	return c.Connect(context.Background())
}

func (d *Driver) OpenConnector(name string) (driver.Connector, error) {
	return &Connector{name: name}, nil
}
