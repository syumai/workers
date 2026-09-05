// Deprecated: use github.com/syumai/workers-go/cloudflare instead.
package cloudflare

import (
	src "github.com/syumai/workers-go/cloudflare"
)

type DurableObjectId = src.DurableObjectId
type DurableObjectNamespace = src.DurableObjectNamespace
type DurableObjectStub = src.DurableObjectStub

var GetBinding = src.GetBinding
var Getenv = src.Getenv

type KVNamespace = src.KVNamespace
type KVNamespaceGetOptions = src.KVNamespaceGetOptions
type KVNamespaceListKey = src.KVNamespaceListKey
type KVNamespaceListOptions = src.KVNamespaceListOptions
type KVNamespaceListResult = src.KVNamespaceListResult
type KVNamespacePutOptions = src.KVNamespacePutOptions

var NewDurableObjectNamespace = src.NewDurableObjectNamespace
var NewKVNamespace = src.NewKVNamespace
var NewR2Bucket = src.NewR2Bucket
var PassThroughOnException = src.PassThroughOnException

type R2Bucket = src.R2Bucket
type R2HTTPMetadata = src.R2HTTPMetadata
type R2Object = src.R2Object
type R2Objects = src.R2Objects
type R2PutOptions = src.R2PutOptions

var WaitUntil = src.WaitUntil
