package utils

import (
	"fmt"
)

type AppErrCode string

// StructuredError - uniform interface for transmitting errors across the system.
type AppError struct {
	code      AppErrCode
	message   string
	tags      map[string]string
	cause     error
	extraData map[string]any
}

func (ae *AppError) initTags() {
	if ae.tags == nil {
		ae.tags = make(map[string]string)
	}
}

func (ae *AppError) initExtraData() {
	if ae.extraData == nil {
		ae.extraData = make(map[string]any)
	}
}

func (ae *AppError) Error() string {
	return fmt.Sprintf("%v : %v", ae.code, ae.message)
}

func (ae *AppError) Unwrap() error {
	return ae.cause
}

func (ae *AppError) GetCode() string {
	return string(ae.code)
}

func (ae *AppError) GetMessage() string {
	return ae.message
}

func (ae *AppError) AddTag(key string, val string) {
	ae.initTags()
	ae.tags[key] = val
}

func (ae *AppError) AddExtraData(key string, val any) {
	ae.initExtraData()
	ae.extraData[key] = val
}

func (ae *AppError) GetExtraData() map[string]any {
	return ae.extraData
}

func (ae *AppError) GetTags() map[string]string {
	return ae.tags
}

func (ae *AppError) AddTags(tags map[string]string) map[string]string {
	for k, v := range tags {
		ae.AddTag(k, v)
	}
	return ae.tags
}

//func (ae *AppError) GetSentryEvent() (*sentry.Event, map[string]any) {
//	sEvent, extraDetails := errors.BuildSentryReport(ae)
//	if sEvent.Tags == nil {
//		sEvent.Tags = make(map[string]string)
//	}
//	for k, v := range ae.tags {
//		sEvent.Tags[k] = v
//	}
//	return sEvent, extraDetails
//}

func AppErrWithCode(code AppErrCode, message string) *AppError {
	err := &AppError{
		code:    code,
		message: message,
	}
	return err
}

func AppErrWithError(code AppErrCode, message string, e error) *AppError {
	err := &AppError{
		code:    code,
		message: message,
	}
	return err
}

const (
	ERR_ITEM_NOT_FOUND AppErrCode = "ITEM_NOT_FOUND"
	ERR_INVALID_INPUT  AppErrCode = "INVALID_INPUT"
)

func ErrInvalidInput(message string) *AppError {
	return AppErrWithCode(ERR_INVALID_INPUT, message)
}
