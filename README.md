# fresh

[![Build Status](https://travis-ci.org/vicanso/fresh.svg?branch=master)](https://travis-ci.org/vicanso/fresh)

HTTP response freshness testing，it is copied from [fresh](https://github.com/jshttp/fresh) by golang.

## API

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

