package main

import (
	"fmt"
	"io"
	"net/http"

	"github.com/syumai/workers"
)

// bucketName is R2 bucket name defined in wrangler.toml.
const bucketName = "BUCKET"

func handleErr(w http.ResponseWriter, msg string) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(msg))
}

// This example is based on implementation in syumai/workers-playground
// * https://github.com/syumai/workers-playground/blob/e32881648ccc055e3690a0d9c750a834261c333e/r2-image-viewer/src/index.ts#L30
func handler(w http.ResponseWriter, req *http.Request) {
	bucket, err := NewR2Bucket(bucketName)
	if err != nil {
		handleErr(w, "failed to get R2Bucket\n")
		return
	}
	imgPath := req.URL.Path
	imgObj, err := bucket.Get(imgPath)
	if err != nil {
		handleErr(w, "failed to get R2Object\n")
		return
	}
	w.Header().Set("Cache-Control", "public, max-age=14400")
	w.Header().Set("ETag", fmt.Sprintf("W/%s", imgObj.HTTPETag))
	contentType := "application/octet-stream"
	if imgObj.HTTPMetadata.ContentType != nil {
		contentType = *imgObj.HTTPMetadata.ContentType
	}
	w.Header().Set("Content-Type", contentType)
	io.Copy(w, imgObj.Body)
}

func main() {
	workers.Serve(http.HandlerFunc(handler))
}
