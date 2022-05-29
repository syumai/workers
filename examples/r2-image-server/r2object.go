package main

import (
	"errors"
	"fmt"
	"io"
	"syscall/js"
	"time"
)

// R2Object represents JavaScript side's R2Object.
// * https://github.com/cloudflare/workers-types/blob/3012f263fb1239825e5f0061b267c8650d01b717/index.d.ts#L1094
type R2Object struct {
	instance       js.Value
	Key            string
	Version        string
	Size           int
	ETag           string
	HTTPETag       string
	Uploaded       time.Time
	HTTPMetadata   R2HTTPMetadata
	CustomMetadata map[string]string
	// Body is a body of R2Object.
	// This value becomes nil when `Head` method of R2Bucket is called.
	Body io.Reader
}

// TODO: implement
// - https://github.com/cloudflare/workers-types/blob/3012f263fb1239825e5f0061b267c8650d01b717/index.d.ts#L1106
// func (o *R2Object) WriteHTTPMetadata(headers http.Header) {
// }

func (o *R2Object) BodyUsed() (bool, error) {
	v := o.instance.Get("bodyUsed")
	if v.IsUndefined() {
		return false, errors.New("bodyUsed doesn't exist for this R2Object")
	}
	return v.Bool(), nil
}

// toR2Object converts JavaScript side's R2Object to *R2Object.
// * https://github.com/cloudflare/workers-types/blob/3012f263fb1239825e5f0061b267c8650d01b717/index.d.ts#L1094
func toR2Object(v js.Value) (*R2Object, error) {
	uploaded, err := dateToTime(v.Get("uploaded"))
	if err != nil {
		return nil, fmt.Errorf("error converting uploaded: %w", err)
	}
	r2Meta, err := toR2HTTPMetadata(v.Get("httpMetadata"))
	if err != nil {
		return nil, fmt.Errorf("error converting httpMetadata: %w", err)
	}
	bodyVal := v.Get("body")
	var body io.Reader
	if !bodyVal.IsUndefined() {
		body = convertStreamReaderToReader(v.Get("body").Call("getReader"))
	}
	return &R2Object{
		instance:       v,
		Key:            v.Get("key").String(),
		Version:        v.Get("version").String(),
		Size:           v.Get("size").Int(),
		ETag:           v.Get("etag").String(),
		HTTPETag:       v.Get("httpEtag").String(),
		Uploaded:       uploaded,
		HTTPMetadata:   r2Meta,
		CustomMetadata: strRecordToMap(v.Get("customMetadata")),
		Body:           body,
	}, nil
}

// R2HTTPMetadata represents metadata of R2 Object.
// * https://github.com/cloudflare/workers-types/blob/3012f263fb1239825e5f0061b267c8650d01b717/index.d.ts#L1053
type R2HTTPMetadata struct {
	ContentType        *string
	ContentLanguage    *string
	ContentDisposition *string
	ContentEncoding    *string
	CacheControl       *string
	CacheExpiry        *time.Time
}

func toR2HTTPMetadata(v js.Value) (R2HTTPMetadata, error) {
	cacheExpiry, err := maybeDate(v.Get("cacheExpiry"))
	if err != nil {
		return R2HTTPMetadata{}, fmt.Errorf("error converting cacheExpiry: %w", err)
	}
	return R2HTTPMetadata{
		ContentType:        maybeString(v.Get("contentType")),
		ContentLanguage:    maybeString(v.Get("contentLanguage")),
		ContentDisposition: maybeString(v.Get("contentDisposition")),
		ContentEncoding:    maybeString(v.Get("contentEncoding")),
		CacheControl:       maybeString(v.Get("cacheControl")),
		CacheExpiry:        cacheExpiry,
	}, nil
}
