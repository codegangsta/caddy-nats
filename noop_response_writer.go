package caddynats

import "net/http"

type noopResponseWriter struct {
	headers http.Header
}

func (noopResponseWriter) Write(p []byte) (int, error) {
	return len(p), nil
}

func (noopResponseWriter) WriteHeader(statusCode int) {
	//noop
}

func (n noopResponseWriter) Header() http.Header {
	if n.headers == nil {
		n.headers = http.Header{}
	}
	return n.headers
}
