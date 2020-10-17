# fresh

[![Build Status](https://travis-ci.org/vicanso/fresh.svg?branch=master)](https://travis-ci.org/vicanso/fresh)

HTTP response freshness testing，it is copied from [fresh](https://github.com/jshttp/fresh) by golang.

## API

### Fresh

- `RequestHeader`
- `ResponseHeader`

```go
req := httptest.NewRequest("GET", "/users/me", nil)
resp := httptest.NewRecorder()
// true
Fresh(req.Header, resp.Header)
```

### Check

- `modifiedSince` IfNoneMatch of request header field

- `noneMatch` IfNoneMatch of request header field

- `cacheControl` Cache-Control of request header field

- `lastModified` LastModified of response header field

- `eTag` ETag of response header field


```go
Check([]byte("Sat, 01 Jan 2000 00:00:00 GMT"), []byte("\"foo\""), nil, []byte("Sat, 01 Jan 2000 00:00:00 GMT"), []byte("\"foo\""))
```
