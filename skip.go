package gotrail

import "context"

func Skip(ctx context.Context, reason string) {
	end(ctx, endOptions{
		status: statusSuccess,
		reason: reason,
	})
}

func (s *Span) Skip(reason string) {
	s.end(endOptions{
		status: statusSuccess,
		reason: reason,
	})
}
