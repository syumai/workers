//go:build js && wasm

package translate_test

import (
	"syscall/js"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
	"github.com/syumai/workers/internal/translate"
)

var eval = js.Global().Call("eval", `(s) => eval(s)`)
var valueCmp = eval.Invoke(`(a, b) => a === b`)
var arrayCmp = eval.Invoke(`(a, b) => {
	for (var i = 0; i < a.length; i++) {
		if (a[i] !== b[i]) {
			return false;
		}
	}
	return true;
}`)
var objectCmp = eval.Invoke(`(a, b) => JSON.stringify(a) === JSON.stringify(b)`)

func TestToJS(t *testing.T) {
	t.Run("Value", func(t *testing.T) {
		t.Run("Func", func(t *testing.T) {
			// WIP: to do
		})
		t.Run("Bool", func(t *testing.T) {
			expected := eval.Invoke(`true`)
			actual := translate.ToJS(true)

			assert.True(t, valueCmp.Invoke(expected, actual).Bool())
		})
		t.Run("Int", func(t *testing.T) {
			expected := eval.Invoke(`100`)
			actual := translate.ToJS(int(100))

			assert.True(t, valueCmp.Invoke(expected, actual).Bool())
		})
		t.Run("Int8", func(t *testing.T) {
			expected := eval.Invoke(`100`)
			actual := translate.ToJS(int8(100))

			assert.True(t, valueCmp.Invoke(expected, actual).Bool())
		})
		t.Run("Int16", func(t *testing.T) {
			expected := eval.Invoke(`100`)
			actual := translate.ToJS(int16(100))

			assert.True(t, valueCmp.Invoke(expected, actual).Bool())
		})
		t.Run("Int32", func(t *testing.T) {
			expected := eval.Invoke(`100`)
			actual := translate.ToJS(int32(100))

			assert.True(t, valueCmp.Invoke(expected, actual).Bool())
		})
		t.Run("Int64", func(t *testing.T) {
			expected := eval.Invoke(`100`)
			actual := translate.ToJS(int64(100))

			assert.True(t, valueCmp.Invoke(expected, actual).Bool())
		})
		t.Run("Uint", func(t *testing.T) {
			expected := eval.Invoke(`100`)
			actual := translate.ToJS(uint(100))

			assert.True(t, valueCmp.Invoke(expected, actual).Bool())
		})
		t.Run("Uint8", func(t *testing.T) {
			expected := eval.Invoke(`100`)
			actual := translate.ToJS(uint8(100))

			assert.True(t, valueCmp.Invoke(expected, actual).Bool())
		})
		t.Run("Uint16", func(t *testing.T) {
			expected := eval.Invoke(`100`)
			actual := translate.ToJS(uint16(100))

			assert.True(t, valueCmp.Invoke(expected, actual).Bool())
		})
		t.Run("Uint32", func(t *testing.T) {
			expected := eval.Invoke(`100`)
			actual := translate.ToJS(uint32(100))

			assert.True(t, valueCmp.Invoke(expected, actual).Bool())
		})
		t.Run("Uint64", func(t *testing.T) {
			expected := eval.Invoke(`100`)
			actual := translate.ToJS(uint64(100))

			assert.True(t, valueCmp.Invoke(expected, actual).Bool())
		})
		t.Run("Uintptr", func(t *testing.T) {
			expected := eval.Invoke(`100`)
			actual := translate.ToJS(uintptr(100))

			assert.True(t, valueCmp.Invoke(expected, actual).Bool())
		})
		t.Run("UnsafePointer", func(t *testing.T) {
			expected := eval.Invoke(`100`)
			actual := translate.ToJS(unsafe.Pointer(uintptr(100)))

			assert.True(t, valueCmp.Invoke(expected, actual).Bool())
		})
		t.Run("Float32", func(t *testing.T) {
			expected := eval.Invoke(`3.5`)
			actual := translate.ToJS(float32(3.5))

			assert.True(t, valueCmp.Invoke(expected, actual).Bool())
		})
		t.Run("Float64", func(t *testing.T) {
			expected := eval.Invoke(`3.5`)
			actual := translate.ToJS(float64(3.5))

			assert.True(t, valueCmp.Invoke(expected, actual).Bool())
		})
		t.Run("String", func(t *testing.T) {
			expected := eval.Invoke(`"hello"`)
			actual := translate.ToJS("hello")

			assert.True(t, valueCmp.Invoke(expected, actual).Bool())
		})
	})

	t.Run("Array", func(t *testing.T) {
		t.Run("Standard Array", func(t *testing.T) {
			expected := eval.Invoke(`['a',1,true]`)
			actual := translate.ToJS([3]any{"a", 1, true})

			assert.True(t, arrayCmp.Invoke(expected, actual).Bool())
		})
		t.Run("Standard Slice", func(t *testing.T) {
			expected := eval.Invoke(`['a',1,true]`)
			actual := translate.ToJS([]any{"a", 1, true})

			assert.True(t, arrayCmp.Invoke(expected, actual).Bool())
		})
		t.Run("Null Slice", func(t *testing.T) {
			// WIP: need to fix later
			expected := eval.Invoke(`null`)
			actual := translate.ToJS(make([]any, 3))

			assert.True(t, valueCmp.Invoke(expected, actual).Bool())
		})
	})

	t.Run("Map", func(t *testing.T) {
		t.Run("Standard", func(t *testing.T) {
			expected := eval.Invoke(`({a:"123"})`)
			actual := translate.ToJS(map[string]string{"a": "123"})

			assert.True(t, objectCmp.Invoke(expected, actual).Bool())
		})
		t.Run("Child", func(t *testing.T) {
			expected := eval.Invoke(`({a:{b:"c"}})`)
			actual := translate.ToJS(map[string]map[string]string{"a": {"b": "c"}})

			assert.True(t, objectCmp.Invoke(expected, actual).Bool())
		})
	})

	t.Run("Struct", func(t *testing.T) {
		t.Run("Standard", func(t *testing.T) {
			expected := eval.Invoke(`({a:"123"})`)
			actual := translate.ToJS(struct {
				A string `key:"a"`
			}{
				A: "123",
			})

			assert.True(t, objectCmp.Invoke(expected, actual).Bool())
		})
		t.Run("Child", func(t *testing.T) {
			expected := eval.Invoke(`({a:{b:"123"}})`)
			actual := translate.ToJS(struct {
				A struct {
					B string `key:"b"`
				} `key:"a"`
			}{
				A: struct {
					B string `key:"b"`
				}{
					B: "123",
				},
			})

			assert.True(t, objectCmp.Invoke(expected, actual).Bool())
		})
		t.Run("Ignore Tag", func(t *testing.T) {
			expected := eval.Invoke(`({a:"123"})`)
			actual := translate.ToJS(struct {
				A string `key:"a"`
				B string `key:"-"`
			}{
				A: "123",
				B: "456",
			})

			assert.True(t, objectCmp.Invoke(expected, actual).Bool())
		})
		t.Run("Without Tag", func(t *testing.T) {
			expected := eval.Invoke(`({A:"123"})`)
			actual := translate.ToJS(struct {
				A string
			}{
				A: "123",
			})

			assert.True(t, objectCmp.Invoke(expected, actual).Bool())
		})
		t.Run("Omitempty Tag", func(t *testing.T) {
			expected := eval.Invoke(`({a:""})`)
			actual := translate.ToJS(struct {
				A string `key:"a"`
				B string `key:"b,omitempty"`
			}{})

			assert.True(t, objectCmp.Invoke(expected, actual).Bool())
		})
	})
}

type S struct {
}

func (S) ToJS() js.Value {
	return js.ValueOf("hello")
}

func TestGo2JS(t *testing.T) {
	t.Run("Standard", func(t *testing.T) {
		expected := eval.Invoke(`"hello"`)
		actual := translate.ToJS(S{})

		assert.True(t, valueCmp.Invoke(expected, actual).Bool())
	})
	t.Run("Child", func(t *testing.T) {
		expected := eval.Invoke(`({a:"hello"})`)
		actual := translate.ToJS(struct {
			A S `key:"a"`
		}{
			A: S{},
		})

		assert.True(t, objectCmp.Invoke(expected, actual).Bool())
	})
}
