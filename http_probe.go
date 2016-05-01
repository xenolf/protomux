package protomux

import "bytes"
import "net/http"

var httpMethods = []string{
	http.MethodGet,
	http.MethodHead,
	http.MethodPost,
	http.MethodPut,
	http.MethodPatch,
	http.MethodDelete,
	http.MethodConnect,
	http.MethodOptions,
	http.MethodTrace,
}

type httpProbe struct {
}

func (p httpProbe) apply(data []byte) bool {

	// If the buffer contains "HTTP" we have HTTP/1.1 with version
	if bytes.Contains(data, []byte("HTTP")) {
		return true
	}

	// If it's not matching, we could have HTTP/1.0 without a version. Try to probe with HTTP methods.
	for _, method := range httpMethods {
		if bytes.HasPrefix(data, []byte(method)) {
			return true
		}
	}

	return false
}
