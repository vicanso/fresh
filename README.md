# fresh

HTTP response freshness testing

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

