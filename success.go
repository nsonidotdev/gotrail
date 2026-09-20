package gotrail

import "context"

func Success(ctx context.Context) {
	end(ctx, endOptions{
		status: statusSuccess,
	})
}

func (s *Span) Success() {
	s.end(endOptions{
		status: statusSuccess,
	})
}
