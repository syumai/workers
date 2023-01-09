package d1

import (
	"context"
	"database/sql/driver"

	"github.com/syumai/workers/internal/jsutil"
)

type Connector struct {
	name string
}

var (
	_ driver.Connector = (*Connector)(nil)
)

func (c *Connector) Connect(context.Context) (driver.Conn, error) {
	v := jsutil.Global.Get(c.name)
	if v.IsUndefined() {
		return nil, ErrDatabaseNotFound
	}
	return &Conn{dbObj: v}, nil
}

func (c *Connector) Driver() driver.Driver {
	return &Driver{}
}
