package fresh

import (
	"log"
	"testing"
)

func TestFresh(t *testing.T) {
	// when a non-conditional GET is performed
	reqHeader := &RequestHeader{}
	resHeader := &ResponseHeader{}
	if Fresh(reqHeader, resHeader) != false {
		t.Fatalf("should be stale")
	}

	// when ETags match
	reqHeader = &RequestHeader{
		IfNoneMatch: []byte("\"foo\""),
	}
	resHeader = &ResponseHeader{
		ETag: []byte("\"foo\""),
	}
	if Fresh(reqHeader, resHeader) != true {
		t.Fatalf("should be fresh'")
	}

	// when ETags mismatch
	reqHeader = &RequestHeader{
		IfNoneMatch: []byte("\"foo\""),
	}
	resHeader = &ResponseHeader{
		ETag: []byte("\"bar\""),
	}
	if Fresh(reqHeader, resHeader) != false {
		t.Fatalf("should be stale'")
	}

	// when at least one matches
	reqHeader = &RequestHeader{
		IfNoneMatch: []byte(" \"bar\" , \"foo\""),
	}
	resHeader = &ResponseHeader{
		ETag: []byte("\"foo\""),
	}
	if Fresh(reqHeader, resHeader) != true {
		t.Fatalf("should be fresh")
	}

	// when etag is missing
	reqHeader = &RequestHeader{
		IfNoneMatch: []byte("\"foo\""),
	}
	resHeader = &ResponseHeader{}
	if Fresh(reqHeader, resHeader) != false {
		t.Fatalf("should be stale'")
	}

	// when ETag is weak
	reqHeader = &RequestHeader{
		IfNoneMatch: []byte("W/\"foo\""),
	}
	resHeader = &ResponseHeader{
		ETag: []byte("W/\"foo\""),
	}
	if Fresh(reqHeader, resHeader) != true {
		t.Fatalf("should be fresh on exact match")
	}
	resHeader = &ResponseHeader{
		ETag: []byte("\"foo\""),
	}
	if Fresh(reqHeader, resHeader) != true {
		t.Fatalf("should be fresh on strong match")
	}

	// when ETag is strong
	reqHeader = &RequestHeader{
		IfNoneMatch: []byte("\"foo\""),
	}
	resHeader = &ResponseHeader{
		ETag: []byte("\"foo\""),
	}
	if Fresh(reqHeader, resHeader) != true {
		t.Fatalf("should be fresh on exact match")
	}
	resHeader = &ResponseHeader{
		ETag: []byte("W/\"foo\""),
	}
	if Fresh(reqHeader, resHeader) != true {
		t.Fatalf("sshould be fresh on weak match")
	}

	// when * is given
	reqHeader = &RequestHeader{
		IfNoneMatch: []byte("*"),
	}
	resHeader = &ResponseHeader{
		ETag: []byte("\"foo\""),
	}
	if Fresh(reqHeader, resHeader) != true {
		t.Fatalf("should be fresh")
	}
	reqHeader = &RequestHeader{
		IfNoneMatch: []byte("*, \"bar\""),
	}
	if Fresh(reqHeader, resHeader) != false {
		t.Fatalf("should get ignored if not only value")
	}

	// when modified since the date
	reqHeader = &RequestHeader{
		IfModifiedSince: []byte("Sat, 01 Jan 2000 00:00:00 GMT"),
	}
	resHeader = &ResponseHeader{
		LastModified: []byte("Sat, 01 Jan 2000 01:00:00 GMT"),
	}
	if Fresh(reqHeader, resHeader) != false {
		t.Fatalf("should be stale")
	}

	// when unmodified since the date
	reqHeader = &RequestHeader{
		IfModifiedSince: []byte("Sat, 01 Jan 2000 01:00:00 GMT"),
	}
	resHeader = &ResponseHeader{
		LastModified: []byte("Sat, 01 Jan 2000 00:00:00 GMT"),
	}
	if Fresh(reqHeader, resHeader) != true {
		t.Fatalf("should be fresh")
	}

	// when Last-Modified is missing
	reqHeader = &RequestHeader{
		IfModifiedSince: []byte("Sat, 01 Jan 2000 01:00:00 GMT"),
	}
	resHeader = &ResponseHeader{}
	if Fresh(reqHeader, resHeader) != false {
		t.Fatalf("should be stale")
	}

	// with invalid If-Modified-Since date
	reqHeader = &RequestHeader{
		IfModifiedSince: []byte("foo"),
	}
	resHeader = &ResponseHeader{
		LastModified: []byte("Sat, 01 Jan 2000 00:00:00 GMT"),
	}
	if Fresh(reqHeader, resHeader) != false {
		t.Fatalf("should be stale")
	}

	// with invalid Last-Modified date
	reqHeader = &RequestHeader{
		IfModifiedSince: []byte("Sat, 01 Jan 2000 00:00:00 GMT"),
	}
	resHeader = &ResponseHeader{
		LastModified: []byte("foo"),
	}
	if Fresh(reqHeader, resHeader) != false {
		t.Fatalf("should be stale")
	}

	// when requested with If-Modified-Since and If-None-Match

	// both match
	log.Println("both match")
	reqHeader = &RequestHeader{
		IfNoneMatch:     []byte("\"foo\""),
		IfModifiedSince: []byte("Sat, 01 Jan 2000 00:00:00 GMT"),
	}
	resHeader = &ResponseHeader{
		ETag:         []byte("\"foo\""),
		LastModified: []byte("Sat, 01 Jan 2000 00:00:00 GMT"),
	}
	if Fresh(reqHeader, resHeader) != true {
		t.Fatalf("should be fresh")
	}
	// when only ETag matches
	reqHeader = &RequestHeader{
		IfNoneMatch:     []byte("\"foo\""),
		IfModifiedSince: []byte("Sat, 01 Jan 2000 00:00:00 GMT"),
	}
	resHeader = &ResponseHeader{
		ETag:         []byte("\"foo\""),
		LastModified: []byte("Sat, 01 Jan 2000 01:00:00 GMT'"),
	}
	if Fresh(reqHeader, resHeader) != false {
		t.Fatalf("should be stale")
	}
	// when only Last-Modified matches
	reqHeader = &RequestHeader{
		IfNoneMatch:     []byte("\"foo\""),
		IfModifiedSince: []byte("Sat, 01 Jan 2000 00:00:00 GMT"),
	}
	resHeader = &ResponseHeader{
		ETag:         []byte("\"bar\""),
		LastModified: []byte("Sat, 01 Jan 2000 00:00:00 GMT'"),
	}
	if Fresh(reqHeader, resHeader) != false {
		t.Fatalf("should be stale")
	}

	// when none match
	reqHeader = &RequestHeader{
		IfNoneMatch:     []byte("\"foo\""),
		IfModifiedSince: []byte("Sat, 01 Jan 2000 00:00:00 GMT"),
	}
	resHeader = &ResponseHeader{
		ETag:         []byte("\"bar\""),
		LastModified: []byte("Sat, 01 Jan 2000 01:00:00 GMT"),
	}
	if Fresh(reqHeader, resHeader) != false {
		t.Fatalf("should be stale")
	}

	// when requested with Cache-Control: no-cache
	reqHeader = &RequestHeader{
		CacheControl: []byte("no-cache"),
	}
	resHeader = &ResponseHeader{}
	if Fresh(reqHeader, resHeader) != false {
		t.Fatalf("should be stale")
	}
	// when ETags match
	reqHeader = &RequestHeader{
		CacheControl: []byte("no-cache"),
		IfNoneMatch:  []byte("\"foo\""),
	}
	resHeader = &ResponseHeader{
		ETag: []byte("\"foo\""),
	}
	if Fresh(reqHeader, resHeader) != false {
		t.Fatalf("should be stale")
	}
	// when unmodified since the date
	reqHeader = &RequestHeader{
		CacheControl:    []byte("no-cache"),
		IfModifiedSince: []byte("Sat, 01 Jan 2000 00:00:00 GMT"),
	}
	resHeader = &ResponseHeader{
		ETag:         []byte("\"foo\""),
		LastModified: []byte("Sat, 01 Jan 2000 00:00:00 GMT"),
	}
	if Fresh(reqHeader, resHeader) != false {
		t.Fatalf("should be stale")
	}
}
