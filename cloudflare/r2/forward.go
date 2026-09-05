// Deprecated: use github.com/syumai/workers-go/cloudflare/r2 instead.
package r2

import (
	src "github.com/syumai/workers-go/cloudflare/r2"
)

type Bucket = src.Bucket
type HTTPMetadata = src.HTTPMetadata

var NewBucket = src.NewBucket

type Object = src.Object
type Objects = src.Objects
type PutOptions = src.PutOptions
