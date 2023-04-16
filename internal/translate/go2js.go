//go:build js && wasm

package translate

import (
	"reflect"
	"strings"
	"syscall/js"
)

// Go2JS is the interface implemented by types that can translate themselves to js.Value.
type Go2JS interface {
	ToJS() js.Value
}

var go2jsType = reflect.TypeOf((*Go2JS)(nil)).Elem()

var (
	null   = js.ValueOf(nil)
	array  = js.Global().Get("Array")
	object = js.Global().Get("Object")
)

// ToJS translates Go values to JS values.
// If the encountered value implements the Go2JS interface, call its ToJS function.
func ToJS(x any) js.Value {
	if jv, ok := x.(js.Value); ok {
		return jv
	}
	v := reflect.ValueOf(x)
	t := reflect.TypeOf(x)

	if t.Implements(go2jsType) {
		g2s := v.Interface().(Go2JS)
		return g2s.ToJS()
	}

	switch v.Kind() {
	case reflect.Func, reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32,
		reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Uintptr, reflect.UnsafePointer, reflect.Float32, reflect.Float64, reflect.String:
		return valueToJSValue(v)
	case reflect.Array, reflect.Slice:
		return sliceToJSArray(v)
	case reflect.Map:
		return mapToJSObject(v)
	case reflect.Struct:
		return structToJSObject(v)
	case reflect.Ptr, reflect.Interface:
		return interfaceToJSValue(v)
	default:
		panic("ValueOf: invalid or unsupported value")
	}
}

// valueToJSValue returns js.ValueOf
func valueToJSValue(v reflect.Value) js.Value {
	return js.ValueOf(v.Interface())
}

// interfaceToJSValue converts reflect.Ptr, reflect.Interface to js.Value
func interfaceToJSValue(v reflect.Value) js.Value {
	if v.IsNil() {
		return null
	}
	return ToJS(v.Elem())
}

// sliceToJSArray converts reflect.Array, reflect.Slice to JS Array
func sliceToJSArray(v reflect.Value) js.Value {
	if v.Kind() == reflect.Slice && v.IsNil() {
		return null
	}
	a := array.New()
	for i := 0; i < v.Len(); i++ {
		if v.Index(i).Interface() == nil {
			a.SetIndex(i, null)
			continue
		}

		a.SetIndex(i, ToJS(v.Index(i).Interface()))
	}
	return a
}

// mapToJSObject converts reflect.Map to JS Object
// This is an experimental function because MapKeys() for TinyGo will be implemented in v0.28.0
func mapToJSObject(v reflect.Value) js.Value {
	if v.IsNil() {
		return null
	}
	o := object.New()
	for _, k := range v.MapKeys() {
		o.Set(k.String(), ToJS(v.MapIndex(k).Interface()))
	}
	return o
}

// structToJSObject converts reflect.Struct to JS Object.
// Use the key tag contained in struct as the object key.
func structToJSObject(v reflect.Value) js.Value {
	o := object.New()
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		tag := t.Field(i).Tag.Get("key")
		if tag == "-" {
			continue
		}
		key, opts, _ := strings.Cut(tag, ",")
		if key == "" {
			key = t.Field(i).Name
		}
		if strings.Contains(opts, "omitempty") && v.Field(i).IsZero() {
			continue
		}
		o.Set(key, ToJS(v.Field(i).Interface()))
	}
	return o
}
