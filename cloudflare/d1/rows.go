package d1

import (
	"database/sql/driver"
)

type rows struct{}

var _ driver.Rows = (*rows)(nil)

func (r *rows) Columns() []string {
	panic("implement me")
}

func (r *rows) Close() error {
	panic("implement me")
}

func (r *rows) Next(dest []driver.Value) error {
	panic("implement me")
}
