package common

import (
	"context"
	"encoding/json"
	"github.com/getsentry/sentry-go"
	"jugaldb.com/byob_task/src/internal/domain/httpReqRes"
	"jugaldb.com/byob_task/src/utils"
	"net/http"
)

func SendRawJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func SendJson(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	resp := &httpReqRes.HttpResponse{
		Data: data,
		Code: httpReqRes.RESPONSE_SUCCESS,
	}
	json.NewEncoder(w).Encode(resp)
}

func SendFailureJson(w http.ResponseWriter, data any) {
	resp := &httpReqRes.HttpResponse{
		Data: data,
		Code: httpReqRes.RESPONSE_FAILURE,
	}
	json.NewEncoder(w).Encode(resp)
}

func RateLimitExceeded(ctx context.Context, w http.ResponseWriter, retryAfter string) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Retry-After", retryAfter)
	w.WriteHeader(429)
}

func HandleError(ctx context.Context, w http.ResponseWriter, e error) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")

	var custErr *HttpError
	var appErr *utils.AppError
	switch err := e.(type) {
	case *HttpError:
		custErr = err
	case *utils.AppError:
		appErr = err
		custErr = GetHttpErrorFromAppError(err)
	default:
		custErr = GetErrorFromCode(ERROR_UNKNOWN).(*HttpError)
	}

	// add error log
	utils.GetAppLogger().Error(e)

	// log request to sentry...
	if hub := sentry.GetHubFromContext(ctx); hub != nil {
		hub.WithScope(func(scope *sentry.Scope) {
			if appErr != nil {
				scope.SetExtras(appErr.GetExtraData())
				scope.SetTags(appErr.GetTags())
			}
			hub.CaptureException(e)
		})
	}
	w.WriteHeader(custErr.GetStatusCode())
	jsonPayload := httpReqRes.HttpResponse{
		Code:    custErr.GetErrorCode(),
		Message: custErr.Error(),
	}
	json.NewEncoder(w).Encode(jsonPayload)
}
