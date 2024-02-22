package httpMiddleware

import (
	"jugaldb.com/byob_task/src/utils"
	"net/http"
)

func ApplyMiddlewares(h http.Handler, mw *Middlewares, config *utils.Config) http.Handler {
	return Adapt(h, mw.LogMW(), mw.ErrorMW())
}
