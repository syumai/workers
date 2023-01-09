package d1

import (
	"database/sql"
	"errors"
	"syscall/js"
)

type result struct {
	resultObj js.Value
}

// Result is the interface which represents Cloudflare's D1Result type.
// For `changes` field, RowsAffected method can be used.
// see: https://github.com/cloudflare/workers-types/blob/v3.18.0/src/workers.json#L1608
type Result interface {
	// LastRowId returns id of result's last row.
	// If LastRowId can't be retrieved, this method returns nil.
	LastRowId() *int
	// Duration returns duration of executed query.
	Duration() int
}

var (
	_ sql.Result = (*result)(nil)
	_ Result     = (*result)(nil)
)

func (r *result) LastInsertId() (int64, error) {
	return 0, errors.New("d1: LastInsertId is not implemented. instead of it, please use d1.Result.LastRowId")
}

func (r *result) RowsAffected() (int64, error) {
	return int64(r.resultObj.Get("changes").Int()), nil
}

func (r *result) LastRowId() *int {
	v := r.resultObj.Get("lastRowId")
	if v.IsNull() {
		return nil
	}
	id := v.Int()
	return &id
}

func (r *result) Duration() int {
	return r.resultObj.Get("duration").Int()
}
