package errorsDom

import (
	"fmt"
	"jugaldb.com/byob_task/src/utils"
)

const (
	INVALID_CONTEXT_VALUE = "INVALID_CTX_VAL"
	RATE_LIMIT_EXCEEDED   = "RATE_LIMIT_EXCEEDED"
	PANIC_ERROR           = "PANIC_ERROR"
	INVALID_REQUEST       = "INVALID_REQUEST"
	INVALID_BODY          = "INVALID_BODY"
)

func InvalidContextVal(val string) error {
	return utils.AppErrWithCode(INVALID_CONTEXT_VALUE, fmt.Sprintf("Invalid value in context: %v", val))
}

func PanicError(err error) error {
	return utils.AppErrWithError(PANIC_ERROR, "Something went wrong: "+err.Error(), err)
}

func InvalidRequest() error {
	return utils.AppErrWithCode(INVALID_REQUEST, "The request is invalid")
}

func InvalidBody(message string) error {
	return utils.AppErrWithCode(INVALID_BODY, "The "+message+" is invalid")
}
