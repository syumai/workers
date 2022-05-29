package main

import (
	"fmt"
	"github.com/syumai/workers"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
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
	bucket, err := workers.NewR2Bucket(bucketName)
	if err != nil {
		handleErr(w, "failed to get R2Bucket\n", err)
		return
	}
	imgPath := strings.TrimPrefix(req.URL.Path, "/")
	imgObj, err := bucket.Get(imgPath)
	if err != nil {
		handleErr(w, "failed to get R2Object\n", err)
		return
	}
	if imgObj == nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(fmt.Sprintf("image not found: %s", imgPath)))
		return
	}
	defer imgObj.Body.Close()
	w.Header().Set("Cache-Control", "public, max-age=14400")
	w.Header().Set("ETag", fmt.Sprintf("W/%s", imgObj.HTTPETag))
	contentType := "application/octet-stream"
	if imgObj.HTTPMetadata.ContentType != nil {
		contentType = *imgObj.HTTPMetadata.ContentType
	}
	w.Header().Set("Content-Type", contentType)
	go func() {
		time.Sleep(100 * time.Millisecond)
		if err := imgObj.Body.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	io.Copy(w, imgObj.Body)
}

func main() {
	workers.Serve(http.HandlerFunc(handler))
}
