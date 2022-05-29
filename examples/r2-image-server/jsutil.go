package main

import (
	"fmt"
	"strconv"
	"syscall/js"
	"time"
)

var (
	global              = js.Global()
	objectClass         = global.Get("Object")
	promiseClass        = global.Get("Promise")
	uint8ArrayClass     = global.Get("Uint8Array")
	errorClass          = global.Get("Error")
	readableStreamClass = global.Get("ReadableStream")
)

func newObject() js.Value {
	return objectClass.New()
}

func newUint8Array(size int) js.Value {
	return uint8ArrayClass.New(size)
}

func newPromise(fn js.Func) js.Value {
	return promiseClass.New(fn)
}

func awaitPromise(promiseVal js.Value) (js.Value, error) {
	fmt.Println("await promise")
	fmt.Println(promiseVal.Call("toString").String())
	resultCh := make(chan js.Value)
	errCh := make(chan error)
	var then, catch js.Func
	then = js.FuncOf(func(_ js.Value, args []js.Value) any {
		defer then.Release()
		result := args[0]
		resultCh <- result
		return js.Undefined()
	})
	catch = js.FuncOf(func(_ js.Value, args []js.Value) any {
		defer catch.Release()
		result := args[0]
		errCh <- fmt.Errorf("failed on promise: %s", result.Call("toString").String())
		return js.Undefined()
	})
	promiseVal.Call("then", then).Call("catch", catch)
	select {
	case result := <-resultCh:
		fmt.Println("got result of promise")
		return result, nil
	case err := <-errCh:
		fmt.Println("got error of promise")
		return js.Value{}, err
	}
}

// dateToTime converts JavaScript side's Data object into time.Time.
func dateToTime(v js.Value) (time.Time, error) {
	milliStr := v.Call("getTime").Call("toString").String()
	milli, err := strconv.ParseInt(milliStr, 10, 64)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to convert Date to time.Time: %w", err)
	}
	return time.UnixMilli(milli), nil
}

// strRecordToMap converts JavaScript side's Record<string, string> into map[string]string.
func strRecordToMap(v js.Value) map[string]string {
	entries := objectClass.Call("entries", v)
	entriesLen := entries.Get("length").Int()
	result := make(map[string]string, entriesLen)
	for i := 0; i < entriesLen; i++ {
		entry := entries.Index(i)
		key := entry.Index(0).String()
		value := entry.Index(1).String()
		result[key] = value
	}
	return result
}

// maybeString returns string value of given JavaScript value or returns nil if the value is undefined.
func maybeString(v js.Value) *string {
	if v.IsUndefined() {
		return nil
	}
	s := v.String()
	return &s
}

// maybeDate returns time.Time value of given JavaScript Date value or returns nil if the value is undefined.
func maybeDate(v js.Value) (*time.Time, error) {
	if v.IsUndefined() {
		return nil, nil
	}
	d, err := dateToTime(v)
	if err != nil {
		return nil, err
	}
	return &d, nil
}
