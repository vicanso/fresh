# fresh

[![Build Status](https://travis-ci.org/vicanso/fresh.svg?branch=master)](https://travis-ci.org/vicanso/fresh)

HTTP response freshness testing，it is copied from [fresh](https://github.com/jshttp/fresh) by golang.

## API

### Fresh

- `RequestHeader`
- `ResponseHeader`

```go
reqHeader = &RequestHeader{
  IfNoneMatch:     []byte("\"foo\""),
  IfModifiedSince: []byte("Sat, 01 Jan 2000 00:00:00 GMT"),
}
resHeader = &ResponseHeader{
  ETag:         []byte("\"foo\""),
  LastModified: []byte("Sat, 01 Jan 2000 00:00:00 GMT"),
}
// true
Fresh(reqHeader, resHeader)
```

### Check

- `modifiedSince` IfNoneMatch of requset header field

- `noneMatch` IfNoneMatch of request header field

- `cacheControl` Cache-Control of request header field

- `lastModified` LastModified of response header field

- `etag` ETag of response header field


```go
Check([]byte("Sat, 01 Jan 2000 00:00:00 GMT"), []byte("\"foo\""), nil, []byte("Sat, 01 Jan 2000 00:00:00 GMT"), []byte("\"foo\""))
```
