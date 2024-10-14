package queues

import (
	"bytes"
	"syscall/js"
	"testing"

	"github.com/syumai/workers/internal/jsutil"
)

func TestContentType_mapValue(t *testing.T) {
	tests := []struct {
		name        string
		contentType QueueContentType
		val         any
		want        js.Value
		wantErr     bool
	}{
		{
			name:        "string as text",
			contentType: QueueContentTypeText,
			val:         "hello",
			want:        js.ValueOf("hello"),
		},
		{
			name:        "[]byte as text",
			contentType: QueueContentTypeText,
			val:         []byte("hello"),
			want:        js.ValueOf("hello"),
		},
		{
			name:        "io.Reader as text",
			contentType: QueueContentTypeText,
			val:         bytes.NewBufferString("hello"),
			want:        js.ValueOf("hello"),
		},
		{
			name:        "number as text",
			contentType: QueueContentTypeText,
			val:         42,
			want:        js.Undefined(),
			wantErr:     true,
		},
		{
			name:        "function as text",
			contentType: QueueContentTypeText,
			val:         func() {},
			want:        js.Undefined(),
			wantErr:     true,
		},

		{
			name:        "string as json",
			contentType: QueueContentTypeJSON,
			val:         "hello",
			want:        js.ValueOf("hello"),
		},
		{
			name:        "number as json",
			contentType: QueueContentTypeJSON,
			val:         42,
			want:        js.ValueOf(42),
		},
		{
			name:        "bool as json",
			contentType: QueueContentTypeJSON,
			val:         true,
			want:        js.ValueOf(true),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.contentType.mapValue(tt.val)
			if (err != nil) != tt.wantErr {
				t.Fatalf("%s.mapValue() error = %v, wantErr %v", tt.contentType, err, tt.wantErr)
			}
			if got.String() != tt.want.String() {
				t.Errorf("%s.mapValue() = %v, want %v", tt.contentType, got, tt.want)
			}
		})
	}
}

func TestContentType_mapValue_bytes(t *testing.T) {
	jsOf := func(b []byte) js.Value {
		ua := jsutil.NewUint8Array(len(b))
		js.CopyBytesToJS(ua, b)
		return ua
	}

	tests := []struct {
		name string
		val  any
		want js.Value
	}{
		{
			name: "[]byte as bytes",
			val:  []byte("hello"),
			want: jsOf([]byte("hello")),
		},
		{
			name: "string as bytes",
			val:  "hello",
			want: jsOf([]byte("hello"))},
		{
			name: "io.Reader as bytes",
			val:  bytes.NewBufferString("hello"),
			want: jsOf([]byte("hello")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := QueueContentTypeBytes.mapValue(tt.val)
			if err != nil {
				t.Fatalf("%s.mapValue() got error = %v", QueueContentTypeBytes, err)
			}
			if got.Type() != tt.want.Type() {
				t.Errorf("%s.mapValue() = type %v, want type %v", QueueContentTypeBytes, got, tt.want)
			}
			if got.String() != tt.want.String() {
				t.Errorf("%s.mapValue() = %v, want %v", QueueContentTypeBytes, got, tt.want)
			}
		})
	}
}

func TestContentType_mapValue_map(t *testing.T) {
	val := map[string]interface{}{
		"Name": "Alice",
		"Age":  42,
	}

	tests := []struct {
		name        string
		contentType QueueContentType
	}{
		{
			name:        "json",
			contentType: QueueContentTypeJSON,
		},
		{
			name:        "v8",
			contentType: QueueContentTypeV8,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := tt.contentType.mapValue(val)
			if err != nil {
				t.Fatalf("QueueContentTypeJSON.mapValue() got error = %v", err)
			}
			if got.Type() != js.TypeObject {
				t.Errorf("QueueContentTypeJSON.mapValue() = type %v, want type %v", got, js.TypeObject)
			}
			if got.Get("Name").String() != "Alice" {
				t.Errorf("QueueContentTypeJSON.mapValue() = %v, want %v", got.Get("Name").String(), "Alice")
			}
			if got.Get("Age").Int() != 42 {
				t.Errorf("QueueContentTypeJSON.mapValue() = %v, want %v", got.Get("Age").Int(), 42)
			}
		})
	}
}

type User struct {
	Name string
}

func TestContentType_mapValue_unsupported_types(t *testing.T) {
	t.Run("struct as json", func(t *testing.T) {
		defer func() {
			if p := recover(); p == nil {
				t.Fatalf("QueueContentTypeJSON.mapValue() did not panic")
			}
		}()

		val := User{Name: "Alice"}
		_, _ = QueueContentTypeJSON.mapValue(val)
	})

	t.Run("slice of structs as json", func(t *testing.T) {
		defer func() {
			if p := recover(); p == nil {
				t.Fatalf("QueueContentTypeJSON.mapValue() did not panic")
			}
		}()

		val := User{Name: "Alice"}
		_, _ = QueueContentTypeJSON.mapValue([]User{val})
	})

	t.Run("slice of bytes as json", func(t *testing.T) {
		defer func() {
			if p := recover(); p == nil {
				t.Fatalf("QueueContentTypeJSON.mapValue() did not panic")
			}
		}()

		_, _ = QueueContentTypeJSON.mapValue([]byte("hello"))
	})
}
