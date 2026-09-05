// Deprecated: use github.com/syumai/workers-go/cloudflare/d1 instead.
package d1

import (
	src "github.com/syumai/workers-go/cloudflare/d1"
)

type Conn = src.Conn
type Connector = src.Connector
type Driver = src.Driver

var ErrDatabaseNotFound = src.ErrDatabaseNotFound
var OpenConnector = src.OpenConnector
