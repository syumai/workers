package main

import (
	"fmt"
	"strconv"
	"syscall/js"
	"time"
)

var (
	global      = js.Global()
	objectClass = global.Get("Object")
)

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
