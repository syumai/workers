// Deprecated: use github.com/syumai/workers-go/exp/hono instead.
package hono

import (
	src "github.com/syumai/workers-go/exp/hono"
)

var ChainMiddlewares = src.ChainMiddlewares

type Context = src.Context
type Middleware = src.Middleware

var NewJSResponse = src.NewJSResponse
var NewJSResponseWithBase = src.NewJSResponseWithBase
var ServeMiddleware = src.ServeMiddleware
