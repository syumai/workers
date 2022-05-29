package workers

import (
	"bytes"
	"fmt"
	"io"
	"syscall/js"
)

// readableStreamToReadCloser implements io.ReadCloser sourced from ReadableStream.
// * ReadableStream: https://developer.mozilla.org/en-US/docs/Web/API/ReadableStream
// * ReadableStreamDefaultReader: https://developer.mozilla.org/en-US/docs/Web/API/ReadableStreamDefaultReader
// * This implementation is based on: https://deno.land/std@0.139.0/streams/conversion.ts#L76
type readableStreamToReadCloser struct {
	buf            bytes.Buffer
	readableStream js.Value
	streamReader   js.Value
}

// Read reads bytes from ReadableStreamDefaultReader.
func (rs *readableStreamToReadCloser) Read(p []byte) (n int, err error) {
	if rs.buf.Len() == 0 {
		promise := rs.streamReader.Call("read")
		resultCh := make(chan js.Value)
		errCh := make(chan error)
		var then, catch js.Func
		then = js.FuncOf(func(_ js.Value, args []js.Value) any {
			defer then.Release()
			result := args[0]
			if result.Get("done").Bool() {
				errCh <- io.EOF
				return js.Undefined()
			}
			resultCh <- result.Get("value")
			return js.Undefined()
		})
		catch = js.FuncOf(func(_ js.Value, args []js.Value) any {
			defer catch.Release()
			result := args[0]
			errCh <- fmt.Errorf("JavaScript error on read: %s", result.Call("toString").String())
			return js.Undefined()
		})
		promise.Call("then", then).Call("catch", catch)
		select {
		case result := <-resultCh:
			chunk := make([]byte, result.Get("byteLength").Int())
			_ = js.CopyBytesToGo(chunk, result)
			// The length written is always the same as the length of chunk, so it can be discarded.
			// - https://pkg.go.dev/bytes#Buffer.Write
			_, err := rs.buf.Write(chunk)
			if err != nil {
				return 0, err
			}
		case err := <-errCh:
			return 0, err
		}
	}
	return rs.buf.Read(p)
}

// Close cancels ReadableStream.
// * ReadableStream.cancel: https://developer.mozilla.org/en-US/docs/Web/API/ReadableStream/cancel
func (rs *readableStreamToReadCloser) Close() error {
	fmt.Println("called cancel")

	srP := rs.streamReader.Call("cancel")
	var (
		v   js.Value
		err error
	)
	if v, err = awaitPromise(srP); err != nil {
		return err
	}
	fmt.Println("canceled", v.Call("toString").String())
	return nil
}

// convertReadableStreamToReadCloser converts ReadableStream to io.ReadCloser.
func convertReadableStreamToReadCloser(rs js.Value) io.ReadCloser {
	return &readableStreamToReadCloser{
		readableStream: rs,
		streamReader:   rs.Call("getReader"),
	}
}

// readerToReadableStream implements ReadableStream sourced from io.ReadCloser.
// * ReadableStream: https://developer.mozilla.org/docs/Web/API/ReadableStream
// * This implementation is based on: https://deno.land/std@0.139.0/streams/conversion.ts#L230
type readerToReadableStream struct {
	reader   io.ReadCloser
	chunkBuf []byte
}

// Pull implements ReadableStream's pull method.
// * https://developer.mozilla.org/en-US/docs/Web/API/ReadableStream/ReadableStream#pull
func (rs *readerToReadableStream) Pull(controller js.Value) error {
	n, err := rs.reader.Read(rs.chunkBuf)
	if err == io.EOF {
		if err := rs.reader.Close(); err != nil {
			return err
		}
		controller.Call("close")
		return nil
	}
	if err != nil {
		jsErr := errorClass.New(err.Error())
		controller.Call("error", jsErr)
		if err := rs.reader.Close(); err != nil {
			return err
		}
		return err
	}
	ua := newUint8Array(n)
	_ = js.CopyBytesToJS(ua, rs.chunkBuf[:n])
	controller.Call("enqueue", ua)
	return nil
}

// Cancel implements ReadableStream's cancel method.
// * https://developer.mozilla.org/en-US/docs/Web/API/ReadableStream/ReadableStream#cancel
func (rs *readerToReadableStream) Cancel() error {
	return rs.reader.Close()
}

// https://deno.land/std@0.139.0/streams/conversion.ts#L5
const defaultChunkSize = 16_640

// convertReaderToReadableStream converts io.ReadCloser to ReadableStream.
func convertReaderToReadableStream(reader io.ReadCloser) js.Value {
	stream := &readerToReadableStream{
		reader:   reader,
		chunkBuf: make([]byte, defaultChunkSize),
	}
	rsInit := newObject()
	rsInit.Set("pull", js.FuncOf(func(_ js.Value, args []js.Value) any {
		var cb js.Func
		cb = js.FuncOf(func(this js.Value, pArgs []js.Value) any {
			defer cb.Release()
			resolve := pArgs[0]
			reject := pArgs[1]
			controller := args[0]
			err := stream.Pull(controller)
			if err != nil {
				reject.Invoke(errorClass.New(err.Error()))
				return js.Undefined()
			}
			resolve.Invoke()
			return js.Undefined()
		})
		return newPromise(cb)
	}))
	rsInit.Set("cancel", js.FuncOf(func(js.Value, []js.Value) any {
		err := stream.Cancel()
		if err != nil {
			panic(err)
		}
		return js.Undefined()
	}))
	return readableStreamClass.New(rsInit)
}
