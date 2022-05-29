package main

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/syumai/workers"
)

// bucketName is R2 bucket name defined in wrangler.toml.
const bucketName = "BUCKET"

func handleErr(w http.ResponseWriter, msg string, err error) {
	log.Println(err)
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(msg))
}

// This example is based on implementation in syumai/workers-playground
// * https://github.com/syumai/workers-playground/blob/e32881648ccc055e3690a0d9c750a834261c333e/r2-image-viewer/src/index.ts#L30
func handler(w http.ResponseWriter, req *http.Request) {
	fmt.Println("new R2Bucket")
	bucket, err := NewR2Bucket(bucketName)
	if err != nil {
		handleErr(w, "failed to get R2Bucket\n", err)
		return
	}
	imgPath := req.URL.Path
	fmt.Println("bucket.get")
	imgObj, err := bucket.Get(imgPath)
	if err != nil {
		handleErr(w, "failed to get R2Object\n", err)
		return
	}
	w.Header().Set("Cache-Control", "public, max-age=14400")
	w.Header().Set("ETag", fmt.Sprintf("W/%s", imgObj.HTTPETag))
	contentType := "application/octet-stream"
	if imgObj.HTTPMetadata.ContentType != nil {
		contentType = *imgObj.HTTPMetadata.ContentType
	}
	w.Header().Set("Content-Type", contentType)
	fmt.Println("return result")
	io.Copy(w, imgObj.Body)
}

func main() {
	workers.Serve(http.HandlerFunc(handler))
}
