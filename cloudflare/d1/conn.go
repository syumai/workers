package d1

import (
	"context"
	"database/sql/driver"
	"errors"
	"syscall/js"
)

type Conn struct {
	dbObj js.Value
}

var (
	_ driver.Conn               = (*Conn)(nil)
	_ driver.ConnBeginTx        = (*Conn)(nil)
	_ driver.ConnPrepareContext = (*Conn)(nil)
	_ driver.QueryerContext     = (*Conn)(nil)
)

func (c *Conn) Prepare(query string) (driver.Stmt, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Conn) PrepareContext(ctx context.Context, query string) (driver.Stmt, error) {
	//TODO implement me
	panic("implement me")
}

func (c *Conn) Close() error {
	// do nothing
	return nil
}

func (c *Conn) Begin() (driver.Tx, error) {
	return nil, errors.New("d1: transaction is not currently supported")
}

func (c *Conn) BeginTx(context.Context, driver.TxOptions) (driver.Tx, error) {
	return nil, errors.New("d1: transaction is not currently supported")
}

func (c *Conn) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	//TODO implement me
	panic("implement me")
}
