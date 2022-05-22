package main

import (
	"fmt"
	"io"
	"syscall/js"
)

// R2Bucket represents interface of Cloudflare Worker's R2 Bucket instance.
// - https://developers.cloudflare.com/r2/runtime-apis/#bucket-method-definitions
// - https://github.com/cloudflare/workers-types/blob/3012f263fb1239825e5f0061b267c8650d01b717/index.d.ts#L1006
type R2Bucket interface {
	Head(key string)
	Get(key string) *R2Object
	Put(key string, value io.Reader)
	Delete(key string)
	List() []*R2Object
}

type r2Bucket struct {
	instance js.Value
}

var _ R2Bucket = &r2Bucket{}

func (r *r2Bucket) Head(key string) {
	panic("implement me")
}

func (r *r2Bucket) Get(key string) *R2Object {
	return nil
}

func (r *r2Bucket) Put(key string, value io.Reader) {
	panic("implement me")
}

func (r *r2Bucket) Delete(key string) {
	panic("implement me")
}

func (r *r2Bucket) List() []*R2Object {
	panic("implement me")
}

func NewR2Bucket(varName string) (R2Bucket, error) {
	inst := js.Global().Get(varName)
	if inst.IsUndefined() {
		return nil, fmt.Errorf("%s is undefined", varName)
	}
	return &r2Bucket{instance: inst}, nil
}
