package workers

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
	responseClass       = global.Get("Response")
	headersClass        = global.Get("Headers")
	arrayClass          = global.Get("Array")
	uint8ArrayClass     = global.Get("Uint8Array")
	errorClass          = global.Get("Error")
	readableStreamClass = global.Get("ReadableStream")
	stringClass         = global.Get("String")
	dateClass           = global.Get("Date")
	numberClass         = global.Get("Number")
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

// arrayFrom calls Array.from to given argument and returns result Array.
func arrayFrom(v js.Value) js.Value {
	return arrayClass.Call("from", v)
}

func awaitPromise(promiseVal js.Value) (js.Value, error) {
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
		return result, nil
	case err := <-errCh:
		return js.Value{}, err
	}
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
func maybeString(v js.Value) string {
	if v.IsUndefined() {
		return ""
	}
	return v.String()
}

// maybeDate returns time.Time value of given JavaScript Date value or returns nil if the value is undefined.
func maybeDate(v js.Value) (time.Time, error) {
	if v.IsUndefined() {
		return time.Time{}, nil
	}
	return dateToTime(v)
}

// dateToTime converts JavaScript side's Data object into time.Time.
func dateToTime(v js.Value) (time.Time, error) {
	milliStr := stringClass.Invoke(v.Call("getTime")).String()
	milli, err := strconv.ParseInt(milliStr, 10, 64)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to convert Date to time.Time: %w", err)
	}
	return time.UnixMilli(milli), nil
}

// timeToDate converts Go side's time.Time into Date object.
func timeToDate(t time.Time) js.Value {
	milliStr := strconv.FormatInt(t.UnixMilli(), 10)
	return dateClass.New(numberClass.Call(milliStr))
}
