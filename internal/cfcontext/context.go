package cfcontext

import (
	"context"
	"errors"
	"syscall/js"
)

type runtimeCtxKey struct{}

func New(ctx context.Context, runtimeCtxObj js.Value) context.Context {
	return context.WithValue(ctx, runtimeCtxKey{}, runtimeCtxObj)
}

var ErrRuntimeContextNotFound = errors.New("runtime context was not found")

// MustExtract extracts runtime context object from context.
// This function panics when runtime context object was not found.
func MustExtract(ctx context.Context) js.Value {
	v, ok := ctx.Value(runtimeCtxKey{}).(js.Value)
	if !ok {
		panic(ErrRuntimeContextNotFound)
	}
	return v
}
