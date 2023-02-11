package d1

import (
	"context"
	"database/sql/driver"
	"errors"

	"github.com/syumai/workers/cloudflare/internal/cfruntimecontext"
)

type Conn struct {
	dbName string
}

var (
	_ driver.Conn               = (*Conn)(nil)
	_ driver.ConnBeginTx        = (*Conn)(nil)
	_ driver.ConnPrepareContext = (*Conn)(nil)
)

func (c *Conn) Prepare(query string) (driver.Stmt, error) {
	return nil, errors.New("d1: Prepare is not implemented. please use PrepareContext instead")
}

func (c *Conn) PrepareContext(ctx context.Context, query string) (driver.Stmt, error) {
	dbObj := cfruntimecontext.GetRuntimeContextEnv(ctx).Get(c.dbName)
	if dbObj.IsUndefined() {
		return nil, ErrDatabaseNotFound
	}
	stmtObj := dbObj.Call("prepare", query)
	return &stmt{
		stmtObj: stmtObj,
	}, nil
}

func (c *Conn) Close() error {
	// do nothing
	return nil
}

func (c *Conn) Begin() (driver.Tx, error) {
	return nil, errors.New("d1: Begin is deprecated and not implemented")
}

func (c *Conn) BeginTx(context.Context, driver.TxOptions) (driver.Tx, error) {
	return nil, errors.New("d1: transaction is not currently supported")
}
