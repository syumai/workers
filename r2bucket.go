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
	Put(key string, value io.Reader) error
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
	fmt.Println(v)
	if v.IsNull() {
		return nil, nil
	}
	return toR2Object(v)
}

func (r *r2Bucket) Put(key string, value io.Reader) error {
	panic("implement me")
}

func (r *r2Bucket) Delete(key string) error {
	panic("implement me")
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
