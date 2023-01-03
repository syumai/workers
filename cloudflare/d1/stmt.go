package d1

import (
	"context"
	"database/sql/driver"
	"errors"
)

type stmt struct{}

var (
	_ driver.Stmt             = (*stmt)(nil)
	_ driver.StmtExecContext  = (*stmt)(nil)
	_ driver.StmtQueryContext = (*stmt)(nil)
)

func (s *stmt) Close() error {
	panic("implement me")
}

func (s *stmt) NumInput() int {
	panic("implement me")
}

func (s *stmt) Exec(args []driver.Value) (driver.Result, error) {
	return nil, errors.New("d1: Exec is deprecated and not implemented")
}

func (s *stmt) ExecContext(ctx context.Context, args []driver.NamedValue) (driver.Result, error) {
	panic("implement me")
}

func (s *stmt) Query(args []driver.Value) (driver.Rows, error) {
	return nil, errors.New("d1: Query is deprecated and not implemented")
}

func (s *stmt) QueryContext(ctx context.Context, args []driver.NamedValue) (driver.Rows, error) {
	panic("implement me")
}
