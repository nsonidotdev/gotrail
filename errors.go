package gotrail

import "errors"

var (
	errCtxSpanMistyped  = errors.New("could not parse span")
	errCtxSpanNotFound  = errors.New("span not found")
	errSpanCtxCancelled = errors.New("span cancelled by context")
	errProcTerminated   = errors.New("process terminated")
	errNoSpans          = errors.New("no spans in the span tree")
)
