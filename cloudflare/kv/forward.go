// Deprecated: use github.com/syumai/workers-go/cloudflare/kv instead.
package kv

import (
	src "github.com/syumai/workers-go/cloudflare/kv"
)

type GetOptions = src.GetOptions
type ListKey = src.ListKey
type ListOptions = src.ListOptions
type ListResult = src.ListResult
type Namespace = src.Namespace

var NewNamespace = src.NewNamespace

type PutOptions = src.PutOptions
