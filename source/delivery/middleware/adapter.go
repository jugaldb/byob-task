package httpMiddleware

import "net/http"

type HttpAdapter func(http.Handler) http.Handler

func Adapt(h http.Handler, adapters ...HttpAdapter) http.Handler {
	for _, adapter := range adapters {
		h = adapter(h)
	}
	return h
}
