package cfcontext

import (
	"context"
	"errors"
	"syscall/js"
)

type runtimeCtxKey struct{}
type incomingPropertyKey struct{}

func New(ctx context.Context, runtimeCtxObj, incomingPropertyObj js.Value) context.Context {
	ctx = context.WithValue(ctx, runtimeCtxKey{}, runtimeCtxObj)
	ctx = context.WithValue(ctx, incomingPropertyKey{}, incomingPropertyObj)
	return ctx
}

var ErrRuntimeContextNotFound = errors.New("runtime context was not found")
var ErrIncomingPropertyNotFound = errors.New("incoming property was not found")

// MustExtractRuntimeContext extracts runtime context object from context.
// This function panics when runtime context object was not found.
func MustExtractRuntimeContext(ctx context.Context) js.Value {
	v, ok := ctx.Value(runtimeCtxKey{}).(js.Value)
	if !ok {
		panic(ErrRuntimeContextNotFound)
	}
	return v
}

// MustExtractIncomingProperty extracts incoming property object from context.
// This function panics when incoming property object was not found.
func MustExtractIncomingProperty(ctx context.Context) js.Value {
	v, ok := ctx.Value(incomingPropertyKey{}).(js.Value)
	if !ok {
		panic(ErrIncomingPropertyNotFound)
	}
	return v
}
