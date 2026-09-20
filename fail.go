package gotrail

import "context"

func Fail(ctx context.Context, reason string) {
	finish(ctx, finishOptions{
		status: statusFail,
		reason: reason,
	})
}

func (s *Span) Fail(reason string) {
	s.finish(finishOptions{
		status: statusFail,
		reason: reason,
	})
}
