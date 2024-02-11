package hono

import (
	"strings"
	"syscall/js"
)

type Header interface {
	Add(key, value string)
	Set(key, value string)
	Get(key string) string
	Values(key string) []string
	Entries() []HeaderEntry
	// Write(w io.Writer) // TODO: implement
	// Clone() httpHeader // Not planned to be implemented
}

type HeaderEntry struct {
	Key    string
	Values []string
}

type header struct {
	headerObj js.Value
}

var _ Header = (*header)(nil)

func (h *header) Add(key, value string) {
	h.headerObj.Call("append", key, value)
}

func (h *header) Set(key, value string) {
	h.headerObj.Call("set", key, value)
}

func (h *header) Get(key string) string {
	vs := h.Values(key)
	if len(vs) == 0 {
		return ""
	}
	return vs[0]
}

func (h *header) Values(key string) []string {
	values := h.headerObj.Call("get", key).String()
	return strings.Split(values, ",")
}

func (h *header) Entries() []HeaderEntry {
	var entries []HeaderEntry
	entriesObj := js.Global().Get("Object").Call("entries", h.headerObj)
	for i := 0; i < entriesObj.Length(); i++ {
		entryObj := entriesObj.Index(i)
		key := entryObj.Index(0).String()
		values := entryObj.Index(1).String()
		entries[i] = HeaderEntry{
			Key:    key,
			Values: strings.Split(values, ","),
		}
	}
	return entries
}
