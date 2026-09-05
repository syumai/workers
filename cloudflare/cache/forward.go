// Deprecated: use github.com/syumai/workers-go/cloudflare/cache instead.
package cache

import (
	src "github.com/syumai/workers-go/cloudflare/cache"
)

type Cache = src.Cache
type CacheOption = src.CacheOption
type DeleteOptions = src.DeleteOptions

var ErrCacheNotFound = src.ErrCacheNotFound

type MatchOptions = src.MatchOptions

var New = src.New
var WithNamespace = src.WithNamespace
