package gotrail

import "context"

func Fail(ctx context.Context, reason string) {
	end(ctx, endOptions{
		status: statusFail,
		reason: reason,
	})
}

func (s *Span) Fail(reason string) {
	s.end(endOptions{
		status: statusFail,
		reason: reason,
	})
}
