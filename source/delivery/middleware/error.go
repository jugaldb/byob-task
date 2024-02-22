package httpMiddleware

import (
	"fmt"
	"jugaldb.com/byob_task/src/delivery/http/common"
	"net/http"
)

func (mw *Middlewares) ErrorMW() HttpAdapter {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					common.HandleError(r.Context(), w, fmt.Errorf("%v", rec))
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
