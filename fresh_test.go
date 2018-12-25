package fresh

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func createRequestHeader(modifiedSince, noneMatch, cacheControl string) http.Header {
	req := httptest.NewRequest("GET", "/users/me", nil)
	header := req.Header
	if modifiedSince != "" {
		header.Set(HeaderIfModifiedSince, modifiedSince)
	}

	if noneMatch != "" {
		header.Set(HeaderIfNoneMatch, noneMatch)
	}

	if cacheControl != "" {
		header.Set(HeaderCacheControl, cacheControl)
	}
	return header
}

func createResponseHeader(lastModified, etag string) http.Header {
	resp := httptest.NewRecorder()
	header := resp.Header()
	if lastModified != "" {
		header.Set(HeaderLastModified, lastModified)
	}

	if etag != "" {
		header.Set(HeaderETag, etag)
	}
	return header
}
func TestFresh(t *testing.T) {
	// when a non-conditional GET is performed
	reqHeader := createRequestHeader("", "", "")
	resHeader := createResponseHeader("", "")
	if Fresh(reqHeader, resHeader) != false {
		t.Fatalf("should be stale")
	}

	// when ETags match
	reqHeader = createRequestHeader("", "\"foo\"", "")

	resHeader = createResponseHeader("", "\"foo\"")
	if Fresh(reqHeader, resHeader) != true {
		t.Fatalf("should be fresh'")
	}

	reqHeader = createRequestHeader("", "W/\"foo\"", "")
	resHeader = createResponseHeader("", "\"foo\"")
	if Fresh(reqHeader, resHeader) != true {
		t.Fatalf("should be fresh'")
	}

	reqHeader = createRequestHeader("", "\"foo\"", "")
	resHeader = createResponseHeader("", "W/\"foo\"")
	if Fresh(reqHeader, resHeader) != true {
		t.Fatalf("should be fresh'")
	}

	// when ETags mismatch
	reqHeader = createRequestHeader("", "\"foo\"", "")
	resHeader = createResponseHeader("", "\"bar\"")
	if Fresh(reqHeader, resHeader) != false {
		t.Fatalf("should be stale'")
	}

	// when at least one matches
	reqHeader = createRequestHeader("", " \"bar\" , \"foo\"", "")
	resHeader = createResponseHeader("", "\"foo\"")
	if Fresh(reqHeader, resHeader) != true {
		t.Fatalf("should be fresh")
	}

	// when etag is missing
	reqHeader = createRequestHeader("", "\"foo\"", "")
	resHeader = createResponseHeader("", "")
	if Fresh(reqHeader, resHeader) != false {
		t.Fatalf("should be stale'")
	}

	// when ETag is weak
	reqHeader = createRequestHeader("", "W/\"foo\"", "")
	resHeader = createResponseHeader("", "W/\"foo\"")

	if Fresh(reqHeader, resHeader) != true {
		t.Fatalf("should be fresh on exact match")
	}
	resHeader = createResponseHeader("", "\"foo\"")
	if Fresh(reqHeader, resHeader) != true {
		t.Fatalf("should be fresh on strong match")
	}

	// when ETag is strong
	reqHeader = createRequestHeader("", "\"foo\"", "")
	resHeader = createResponseHeader("", "\"foo\"")

	if Fresh(reqHeader, resHeader) != true {
		t.Fatalf("should be fresh on exact match")
	}
	resHeader = createResponseHeader("", "W/\"foo\"")
	if Fresh(reqHeader, resHeader) != true {
		t.Fatalf("sshould be fresh on weak match")
	}

	// when * is given
	reqHeader = createRequestHeader("", "*", "")
	resHeader = createResponseHeader("", "\"foo\"")
	if Fresh(reqHeader, resHeader) != true {
		t.Fatalf("should be fresh")
	}

	reqHeader = createRequestHeader("", "*, \"bar\"", "")
	if Fresh(reqHeader, resHeader) != false {
		t.Fatalf("should get ignored if not only value")
	}

	// when modified since the date
	reqHeader = createRequestHeader("Sat, 01 Jan 2000 00:00:00 GMT", "", "")
	resHeader = createResponseHeader("Sat, 01 Jan 2000 01:00:00 GMT", "")
	if Fresh(reqHeader, resHeader) != false {
		t.Fatalf("should be stale")
	}

	// when unmodified since the date
	reqHeader = createRequestHeader("Sat, 01 Jan 2000 01:00:00 GMT", "", "")
	resHeader = createResponseHeader("Sat, 01 Jan 2000 00:00:00 GMT", "")

	if Fresh(reqHeader, resHeader) != true {
		t.Fatalf("should be fresh")
	}

	// when Last-Modified is missing
	reqHeader = createRequestHeader("Sat, 01 Jan 2000 01:00:00 GMT", "", "")
	resHeader = createResponseHeader("", "")

	if Fresh(reqHeader, resHeader) != false {
		t.Fatalf("should be stale")
	}

	// with invalid If-Modified-Since date
	reqHeader = createRequestHeader("foo", "", "")
	resHeader = createResponseHeader("Sat, 01 Jan 2000 00:00:00 GMT", "")
	if Fresh(reqHeader, resHeader) != false {
		t.Fatalf("should be stale")
	}

	// with invalid Last-Modified date
	reqHeader = createRequestHeader("Sat, 01 Jan 2000 00:00:00 GMT", "", "")
	resHeader = createResponseHeader("foo", "")
	if Fresh(reqHeader, resHeader) != false {
		t.Fatalf("should be stale")
	}

	// when requested with If-Modified-Since and If-None-Match

	// both match
	log.Println("both match")
	reqHeader = createRequestHeader("Sat, 01 Jan 2000 00:00:00 GMT", "\"foo\"", "")
	resHeader = createResponseHeader("Sat, 01 Jan 2000 00:00:00 GMT", "\"foo\"")
	if Fresh(reqHeader, resHeader) != true {
		t.Fatalf("should be fresh")
	}
	// when only ETag matches
	reqHeader = createRequestHeader("Sat, 01 Jan 2000 00:00:00 GMT", "\"foo\"", "")
	resHeader = createResponseHeader("Sat, 01 Jan 2000 01:00:00 GMT", "\"foo\"")
	if Fresh(reqHeader, resHeader) != false {
		t.Fatalf("should be stale")
	}
	// when only Last-Modified matches
	reqHeader = createRequestHeader("Sat, 01 Jan 2000 00:00:00 GMT", "\"foo\"", "")
	resHeader = createResponseHeader("Sat, 01 Jan 2000 00:00:00 GMT", "\"bar\"")
	if Fresh(reqHeader, resHeader) != false {
		t.Fatalf("should be stale")
	}

	// when none match
	reqHeader = createRequestHeader("Sat, 01 Jan 2000 00:00:00 GMT", "\"foo\"", "")
	resHeader = createResponseHeader("Sat, 01 Jan 2000 01:00:00 GMT", "\"bar\"")
	if Fresh(reqHeader, resHeader) != false {
		t.Fatalf("should be stale")
	}

	// when requested with Cache-Control: no-cache
	reqHeader = createRequestHeader("", "", "no-cache")
	resHeader = createResponseHeader("", "")
	if Fresh(reqHeader, resHeader) != false {
		t.Fatalf("should be stale")
	}
	// when ETags match
	reqHeader = createRequestHeader("", "\"foo\"", "no-cache")
	resHeader = createResponseHeader("", "\"foo\"")
	if Fresh(reqHeader, resHeader) != false {
		t.Fatalf("should be stale")
	}

	// when unmodified since the date
	reqHeader = createRequestHeader("Sat, 01 Jan 2000 00:00:00 GMT", "", "no-cache")
	resHeader = createResponseHeader("Sat, 01 Jan 2000 00:00:00 GMT", "\"foo\"")
	if Fresh(reqHeader, resHeader) != false {
		t.Fatalf("should be stale")
	}
}

func BenchmarkFresh(b *testing.B) {
	b.ResetTimer()
	reqHeader := createRequestHeader("Sat, 01 Jan 2000 00:00:00 GMT", "\"foo\"", "")
	resHeader := createResponseHeader("Sat, 01 Jan 2000 00:00:00 GMT", "\"foo\"")
	for i := 0; i < b.N; i++ {
		Fresh(reqHeader, resHeader)
	}
}
