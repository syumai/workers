package kv

import (
	"fmt"
	"io"
	"syscall/js"

	"github.com/syumai/workers/cloudflare/internal/cfruntimecontext"
	"github.com/syumai/workers/internal/jsutil"
)

// Namespace represents interface of Cloudflare Worker's KV namespace instance.
//   - https://developers.cloudflare.com/workers/runtime-apis/kv/
//   - https://github.com/cloudflare/workers-types/blob/3012f263fb1239825e5f0061b267c8650d01b717/index.d.ts#L850
type Namespace struct {
	instance js.Value
}

// NewNamespace returns Namespace for given variable name.
//   - variable name must be defined in wrangler.toml as kv_namespace's binding.
//   - if the given variable name doesn't exist on runtime context, returns error.
//   - This function panics when a runtime context is not found.
func NewNamespace(varName string) (*Namespace, error) {
	inst := cfruntimecontext.MustGetRuntimeContextEnv().Get(varName)
	if inst.IsUndefined() {
		return nil, fmt.Errorf("%s is undefined", varName)
	}
	return &Namespace{instance: inst}, nil
}

// GetOptions represents Cloudflare KV namespace get options.
//   - https://github.com/cloudflare/workers-types/blob/3012f263fb1239825e5f0061b267c8650d01b717/index.d.ts#L930
type GetOptions struct {
	CacheTTL int
}

func (opts *GetOptions) toJS(type_ string) js.Value {
	obj := jsutil.NewObject()
	obj.Set("type", type_)
	if opts == nil {
		return obj
	}
	if opts.CacheTTL != 0 {
		obj.Set("cacheTtl", opts.CacheTTL)
	}
	return obj
}

// GetString gets string value by the specified key.
//   - if a network error happens, returns error.
func (ns *Namespace) GetString(key string, opts *GetOptions) (string, error) {
	p := ns.instance.Call("get", key, opts.toJS("text"))
	v, err := jsutil.AwaitPromise(p)
	if err != nil {
		return "", err
	}
	return v.String(), nil
}

// GetReader gets stream value by the specified key.
//   - if a network error happens, returns error.
func (ns *Namespace) GetReader(key string, opts *GetOptions) (io.Reader, error) {
	p := ns.instance.Call("get", key, opts.toJS("stream"))
	v, err := jsutil.AwaitPromise(p)
	if err != nil {
		return nil, err
	}
	return jsutil.ConvertReadableStreamToReadCloser(v), nil
}

// ListOptions represents Cloudflare KV namespace list options.
//   - https://github.com/cloudflare/workers-types/blob/3012f263fb1239825e5f0061b267c8650d01b717/index.d.ts#L946
type ListOptions struct {
	Limit  int
	Prefix string
	Cursor string
}

func (opts *ListOptions) toJS() js.Value {
	if opts == nil {
		return js.Undefined()
	}
	obj := jsutil.NewObject()
	if opts.Limit != 0 {
		obj.Set("limit", opts.Limit)
	}
	if opts.Prefix != "" {
		obj.Set("prefix", opts.Prefix)
	}
	if opts.Cursor != "" {
		obj.Set("cursor", opts.Cursor)
	}
	return obj
}

// ListKey represents Cloudflare KV namespace list key.
//   - https://github.com/cloudflare/workers-types/blob/3012f263fb1239825e5f0061b267c8650d01b717/index.d.ts#L940
type ListKey struct {
	Name string
	// Expiration is an expiration of KV value cache. The value `0` means no expiration.
	Expiration int
	// Metadata   map[string]any // TODO: implement
}

// toListKey converts JavaScript side's KVNamespaceListKey to *ListKey.
//   - https://github.com/cloudflare/workers-types/blob/3012f263fb1239825e5f0061b267c8650d01b717/index.d.ts#L940
func toListKey(v js.Value) (*ListKey, error) {
	expVal := v.Get("expiration")
	var exp int
	if !expVal.IsUndefined() {
		exp = expVal.Int()
	}
	return &ListKey{
		Name:       v.Get("name").String(),
		Expiration: exp,
		// Metadata // TODO: implement. This may return an error, so this func signature has an error in return parameters.
	}, nil
}

// ListResult represents Cloudflare KV namespace list result.
//   - https://github.com/cloudflare/workers-types/blob/3012f263fb1239825e5f0061b267c8650d01b717/index.d.ts#L952
type ListResult struct {
	Keys         []*ListKey
	ListComplete bool
	Cursor       string
}

// toListResult converts JavaScript side's KVNamespaceListResult to *ListResult.
//   - https://github.com/cloudflare/workers-types/blob/3012f263fb1239825e5f0061b267c8650d01b717/index.d.ts#L952
func toListResult(v js.Value) (*ListResult, error) {
	keysVal := v.Get("keys")
	keys := make([]*ListKey, keysVal.Length())
	for i := 0; i < len(keys); i++ {
		key, err := toListKey(keysVal.Index(i))
		if err != nil {
			return nil, fmt.Errorf("error converting to ListKey: %w", err)
		}
		keys[i] = key
	}

	cursorVal := v.Get("cursor")
	var cursor string
	if !cursorVal.IsUndefined() {
		cursor = cursorVal.String()
	}

	return &ListResult{
		Keys:         keys,
		ListComplete: v.Get("list_complete").Bool(),
		Cursor:       cursor,
	}, nil
}

// List lists keys stored into the KV namespace.
func (ns *Namespace) List(opts *ListOptions) (*ListResult, error) {
	p := ns.instance.Call("list", opts.toJS())
	v, err := jsutil.AwaitPromise(p)
	if err != nil {
		return nil, err
	}
	return toListResult(v)
}

// PutOptions represents Cloudflare KV namespace put options.
//   - https://github.com/cloudflare/workers-types/blob/3012f263fb1239825e5f0061b267c8650d01b717/index.d.ts#L958
type PutOptions struct {
	Expiration    int
	ExpirationTTL int
	// Metadata // TODO: implement
}

func (opts *PutOptions) toJS() js.Value {
	if opts == nil {
		return js.Undefined()
	}
	obj := jsutil.NewObject()
	if opts.Expiration != 0 {
		obj.Set("expiration", opts.Expiration)
	}
	if opts.ExpirationTTL != 0 {
		obj.Set("expirationTtl", opts.ExpirationTTL)
	}
	return obj
}

// PutString puts string value into KV with key.
//   - if a network error happens, returns error.
func (ns *Namespace) PutString(key string, value string, opts *PutOptions) error {
	p := ns.instance.Call("put", key, value, opts.toJS())
	_, err := jsutil.AwaitPromise(p)
	if err != nil {
		return err
	}
	return nil
}

// PutReader puts stream value into KV with key.
//   - This method copies all bytes into memory for implementation restriction.
//   - if a network error happens, returns error.
func (ns *Namespace) PutReader(key string, value io.Reader, opts *PutOptions) error {
	// fetch body cannot be ReadableStream. see: https://github.com/whatwg/fetch/issues/1438
	b, err := io.ReadAll(value)
	if err != nil {
		return err
	}
	ua := jsutil.NewUint8Array(len(b))
	js.CopyBytesToJS(ua, b)
	p := ns.instance.Call("put", key, ua.Get("buffer"), opts.toJS())
	_, err = jsutil.AwaitPromise(p)
	if err != nil {
		return err
	}
	return nil
}

// Delete deletes key-value pair specified by the key.
//   - if a network error happens, returns error.
func (ns *Namespace) Delete(key string) error {
	p := ns.instance.Call("delete", key)
	_, err := jsutil.AwaitPromise(p)
	if err != nil {
		return err
	}
	return nil
}
