package d1

import "database/sql"

type result struct{}

var _ sql.Result = (*result)(nil)

func (r *result) LastInsertId() (int64, error) {
	panic("implement me")
}

func (r *result) RowsAffected() (int64, error) {
	panic("implement me")
}
