package gotrail

import "context"

func Skip(ctx context.Context, reason string) {
	finish(ctx, finishOptions{
		status: statusSuccess,
		reason: reason,
	})
}

func (s *Span) Skip(reason string) {
	s.finish(finishOptions{
		status: statusSuccess,
		reason: reason,
	})
}
