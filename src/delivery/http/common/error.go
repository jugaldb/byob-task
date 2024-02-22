package common

import (
	errorsDom "jugaldb.com/byob_task/src/internal/domain/errors"
	"jugaldb.com/byob_task/src/utils"
)

type httpErrorInfo struct {
	statusCode int
	message    string
	code       string
}

const (
	ERROR_UNAUTHORIZED = "unauthorized"
	ERROR_UNKNOWN      = "unknown"
	ERROR_BAD_REQ      = "bad_request"
	NPCI_COMMON_ERROR  = "npci_common_error"
)

var httpErrorMap map[string]*httpErrorInfo

type HttpError struct {
	StatusCode int
	AppError   *utils.AppError
	Message    string
}

func (he *HttpError) Error() string {
	if he.Message != "" {
		return he.Message
	}
	return he.AppError.GetMessage()

}

func (he *HttpError) Unwrap() error {
	return he.AppError
}

func (he *HttpError) GetStatusCode() int {
	return he.StatusCode
}

func (he *HttpError) GetErrorCode() string {
	return he.AppError.GetCode()
}

func getErrorMap() map[string]*httpErrorInfo {
	if len(httpErrorMap) == 0 {
		httpErrorMap = make(map[string]*httpErrorInfo)
		httpErrorMap[ERROR_UNAUTHORIZED] = &httpErrorInfo{statusCode: 401, message: "Unauthorized", code: ERROR_UNAUTHORIZED}
		httpErrorMap[ERROR_UNKNOWN] = &httpErrorInfo{statusCode: 503, message: "Something went wrong! please try again later.", code: ERROR_UNKNOWN}
		httpErrorMap[ERROR_BAD_REQ] = &httpErrorInfo{statusCode: 400, message: "Bad Request", code: ERROR_BAD_REQ}
		httpErrorMap[errorsDom.PANIC_ERROR] = &httpErrorInfo{statusCode: 500, message: "Something went wrong! please try again later.", code: errorsDom.PANIC_ERROR}
		httpErrorMap[errorsDom.RATE_LIMIT_EXCEEDED] = &httpErrorInfo{statusCode: 429, code: errorsDom.RATE_LIMIT_EXCEEDED}
		httpErrorMap[errorsDom.INVALID_BODY] = &httpErrorInfo{statusCode: 400, code: errorsDom.INVALID_BODY}
	}
	return httpErrorMap
}

func getErrorMapForStatusCodes() map[string]int {
	errorScMap := make(map[string]int)
	errorScMap[errorsDom.PANIC_ERROR] = 500
	return errorScMap
}

func resolveCode(eCode string) *httpErrorInfo {
	eMap := getErrorMap()
	val, ok := eMap[eCode]
	if !ok {
		val = eMap[ERROR_UNKNOWN]
	}
	return val
}

func GetErrorWithMessage(eCode string, msg string) error {
	hEInfo := resolveCode(eCode)
	aErr := utils.AppErrWithCode(utils.AppErrCode(hEInfo.code), msg)
	return NewHttpErrFromAppErr(hEInfo.statusCode, aErr)
}

func GetErrorFromCode(eCode string) error {
	hEInfo := resolveCode(eCode)
	aErr := utils.AppErrWithCode(utils.AppErrCode(hEInfo.code), hEInfo.message)
	return NewHttpErrFromAppErr(hEInfo.statusCode, aErr)
}

func GetHttpErrorFromAppError(ae *utils.AppError) *HttpError {
	hEInfo := resolveCode(ae.GetCode())
	return NewHttpErrFromAppErr(hEInfo.statusCode, ae)
}

func NewHttpErrFromAppErr(sc int, ae *utils.AppError) *HttpError {
	return &HttpError{
		StatusCode: sc,
		AppError:   ae,
	}
}
