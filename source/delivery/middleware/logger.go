package httpMiddleware

import (
	"jugaldb.com/byob_task/src/utils"
	"net/http"
)

func (mw *Middlewares) LogMW() HttpAdapter {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			url := r.URL.String()
			if url == "/health" || url == "/metrics" {
				h.ServeHTTP(w, r)
			} else {
				utils.GetAppLogger().Printf("%s %s %s %s\n", r.RemoteAddr, r.Method, r.URL, r.UserAgent())
				h.ServeHTTP(w, r)
			}
		})
	}
}
