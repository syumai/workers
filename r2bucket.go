package workers

import (
	"fmt"
	"io"
	"syscall/js"
)

// R2Bucket represents interface of Cloudflare Worker's R2 Bucket instance.
// - https://developers.cloudflare.com/r2/runtime-apis/#bucket-method-definitions
// - https://github.com/cloudflare/workers-types/blob/3012f263fb1239825e5f0061b267c8650d01b717/index.d.ts#L1006
type R2Bucket interface {
	Head(key string) (*R2Object, error)
	Get(key string) (*R2Object, error)
	Put(key string, value io.ReadCloser, opts *R2PutOptions) (*R2Object, error)
	Delete(key string) error
	List() (*R2Objects, error)
}

type r2Bucket struct {
	instance js.Value
}

var _ R2Bucket = &r2Bucket{}

// NewR2Bucket returns R2Bucket for given variable name.
// * variable name must be defined in wrangler.toml.
//   - see example: https://github.com/syumai/workers/tree/main/examples/r2-image-server
// * if the given variable name doesn't exist on global object, returns error.
func NewR2Bucket(varName string) (R2Bucket, error) {
	inst := js.Global().Get(varName)
	if inst.IsUndefined() {
		return nil, fmt.Errorf("%s is undefined", varName)
	}
	return &r2Bucket{instance: inst}, nil
}

// Head returns the result of `head` call to R2Bucket.
// * Body field of *R2Object is always nil for Head call.
// * if the object for given key doesn't exist, returns nil.
// * if a network error happens, returns error.
func (r *r2Bucket) Head(key string) (*R2Object, error) {
	p := r.instance.Call("head", key)
	v, err := awaitPromise(p)
	if err != nil {
		return nil, err
	}
	if v.IsNull() {
		return nil, nil
	}
	return toR2Object(v)
}

// Get returns the result of `get` call to R2Bucket.
// * if the object for given key doesn't exist, returns nil.
// * if a network error happens, returns error.
func (r *r2Bucket) Get(key string) (*R2Object, error) {
	p := r.instance.Call("get", key)
	v, err := awaitPromise(p)
	if err != nil {
		return nil, err
	}
	if v.IsNull() {
		return nil, nil
	}
	return toR2Object(v)
}

// R2PutOptions represents Cloudflare R2 put options.
// * https://github.com/cloudflare/workers-types/blob/3012f263fb1239825e5f0061b267c8650d01b717/index.d.ts#L1128
type R2PutOptions struct {
	HTTPMetadata   R2HTTPMetadata
	CustomMetadata map[string]string
	MD5            string
}

func (opts *R2PutOptions) toJS() js.Value {
	if opts == nil {
		return js.Undefined()
	}
	obj := newObject()
	if opts.HTTPMetadata != (R2HTTPMetadata{}) {
		obj.Set("httpMetadata", opts.HTTPMetadata.toJS())
	}
	if opts.CustomMetadata != nil {
		// convert map[string]string to map[string]any.
		// This makes the map convertible to JS.
		// see: https://pkg.go.dev/syscall/js#ValueOf
		customMeta := make(map[string]any, len(opts.CustomMetadata))
		for k, v := range opts.CustomMetadata {
			customMeta[k] = v
		}
		obj.Set("customMetadata", customMeta)
	}
	if opts.MD5 != "" {
		obj.Set("md5", opts.MD5)
	}
	return obj
}

// Put returns the result of `put` call to R2Bucket.
// * Body field of *R2Object is always nil for Put call.
// * if a network error happens, returns error.
func (r *r2Bucket) Put(key string, value io.ReadCloser, opts *R2PutOptions) (*R2Object, error) {
	body := convertReaderToReadableStream(value)
	p := r.instance.Call("put", key, body, opts.toJS())
	v, err := awaitPromise(p)
	if err != nil {
		return nil, err
	}
	return toR2Object(v)
}

// Delete returns the result of `delete` call to R2Bucket.
// * if a network error happens, returns error.
func (r *r2Bucket) Delete(key string) error {
	p := r.instance.Call("delete", key)
	if _, err := awaitPromise(p); err != nil {
		return err
	}
	return nil
}

// List returns the result of `list` call to R2Bucket.
// * if a network error happens, returns error.
func (r *r2Bucket) List() (*R2Objects, error) {
	p := r.instance.Call("list")
	v, err := awaitPromise(p)
	if err != nil {
		return nil, err
	}
	return toR2Objects(v)
}
