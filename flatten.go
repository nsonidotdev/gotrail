package gotrail

func flattenSpans(t *Trace) []*Span {
	flatSpans := make([]*Span, 0)
	flatSpans = append(flatSpans, t.root)
	return flattenSpansRecursive(t.root, flatSpans)
}

func flattenSpansRecursive(s *Span, flatSpans []*Span) []*Span {
	for _, child := range s.children {
		flatSpans = append(flatSpans, child)
		flatSpans = flattenSpansRecursive(child, flatSpans)
	}

	return flatSpans
}
