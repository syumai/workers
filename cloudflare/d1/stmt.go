package d1

import (
	"context"
	"database/sql/driver"
	"errors"
	"syscall/js"

	"github.com/syumai/workers/internal/jsutil"
)

type stmt struct {
	stmtObj js.Value
}

var (
	_ driver.Stmt             = (*stmt)(nil)
	_ driver.StmtExecContext  = (*stmt)(nil)
	_ driver.StmtQueryContext = (*stmt)(nil)
)

func (s *stmt) Close() error {
	// do nothing
	return nil
}

// NumInput is not supported and always returns -1.
func (s *stmt) NumInput() int {
	return -1
}

func (s *stmt) Exec([]driver.Value) (driver.Result, error) {
	return nil, errors.New("d1: Exec is deprecated and not implemented")
}

// ExecContext executes prepared statement.
// Given []drier.NamedValue's `Name` field will be ignored because Cloudflare D1 client doesn't support it.
func (s *stmt) ExecContext(_ context.Context, args []driver.NamedValue) (driver.Result, error) {
	argValues := make([]any, len(args))
	for i, arg := range args {
		argValues[i] = arg.Value
	}
	resultPromise := s.stmtObj.Call("bind", argValues...).Call("run")
	resultObj, err := jsutil.AwaitPromise(resultPromise)
	if err != nil {
		return nil, err
	}
	return &result{
		resultObj: resultObj,
	}, nil
}

func (s *stmt) Query([]driver.Value) (driver.Rows, error) {
	return nil, errors.New("d1: Query is deprecated and not implemented")
}

func (s *stmt) QueryContext(_ context.Context, args []driver.NamedValue) (driver.Rows, error) {
	argValues := make([]any, len(args))
	for i, arg := range args {
		argValues[i] = arg.Value
	}
	resultPromise := s.stmtObj.Call("bind", argValues...).Call("raw", map[string]any{"columnNames": true})
	rowsObj, err := jsutil.AwaitPromise(resultPromise)
	if err != nil {
		return nil, err
	}
	// If there are no rows to retrieve, length is 0.
	if rowsObj.Length() == 0 {
		return &rows{
			columns: nil,
			rowsObj: rowsObj,
		}, nil
	}

	// The first result array includes the column names.
	colsArray := rowsObj.Index(0)
	colsLen := colsArray.Length()
	cols := make([]string, colsLen)
	for i := 0; i < colsLen; i++ {
		cols[i] = colsArray.Index(i).String()
	}
	// Remove the first result array from the rowsObj.
	rowsObj.Call("shift")
	return &rows{
		columns: cols,
		rowsObj: rowsObj,
	}, nil
}
