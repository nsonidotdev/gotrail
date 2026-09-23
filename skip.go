package gotrail

import "context"

func Skip(ctx context.Context, reason string) {
	end(ctx, endOptions{
		status: statusSkip,
		reason: reason,
	})
}

func (s *Span) Skip(reason string) {
	s.end(endOptions{
		status: statusSkip,
		reason: reason,
	})
}
