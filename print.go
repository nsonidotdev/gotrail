package gotrail

import (
	"fmt"
	"strings"
)

type printOptions struct {
	isRoot      bool
	isLastChild bool
	prefix      string
}

const (
	emptyPrefix = "    "
	midPrefix   = "├── "
	endPrefix   = "└── "
	downPrefix  = "│   "
)

func PrintTrace(trace *Trace) {
	if !isInitialized.Load() {
		return
	}

	fmt.Printf("TRACE: %s\n", trace.ID)
	recursivePrint(trace.root, printOptions{isRoot: true, prefix: ""})
}

func PrintSpanTree(span *Span) {
	if !isInitialized.Load() {
		return
	}

	fmt.Printf("SPAN: %s\n", span.ID)
	recursivePrint(span, printOptions{isRoot: true, prefix: ""})
}

func recursivePrint(s *Span, opts printOptions) {
	lines := formatSpan(s)

	if !opts.isRoot {
		printSpanLine(opts.prefix, downPrefix, "")

		for i, line := range lines {
			var firstPrefix string
			var defaultPrefix string

			if opts.isLastChild {
				firstPrefix = endPrefix
				defaultPrefix = emptyPrefix
			} else {
				firstPrefix = midPrefix
				defaultPrefix = downPrefix
			}

			if i == 0 {
				printSpanLine(opts.prefix, firstPrefix, line)
			} else {
				printSpanLine(opts.prefix, defaultPrefix, line)
			}
		}
	} else {
		for _, line := range lines {
			fmt.Println(line)
		}
	}

	if len(s.children) == 0 {
		return
	}

	for i, childSpan := range s.children {
		isLastChild := i == len(s.children)-1

		var newPrefix string
		if opts.isRoot {
			newPrefix = ""
		} else if !opts.isLastChild {
			newPrefix = opts.prefix + downPrefix
		} else {
			newPrefix = opts.prefix + emptyPrefix
		}

		recursivePrint(childSpan, printOptions{
			isRoot:      false,
			isLastChild: isLastChild,
			prefix:      newPrefix,
		})
	}
}

func formatSpan(s *Span) []string {
	lines := make([]string, 0, 3+len(s.Attributes))

	lines = append(
		lines,
		s.Name,
		fmt.Sprintf("ID: %s", s.ID),
		fmt.Sprintf("status: %s | duration: %dms", strings.ToUpper(string(s.status)), s.Duration.Milliseconds()),
	)

	if s.reason != "" {
		lines = append(lines, fmt.Sprintf("reason: %s", s.reason))
	}

	if len(s.Attributes) != 0 {
		metaLines := formatAttrs(s)
		lines = append(lines, metaLines...)
	}

	return lines
}

func formatAttrs(s *Span) []string {
	lines := make([]string, 0, len(s.Attributes))
	for key, value := range s.Attributes {
		lines = append(lines, fmt.Sprintf("%s: %s", key, value))
	}

	return lines
}

func printSpanLine(prevPrefix string, newPrefix string, line string) {
	formattedLine := fmt.Sprintf("%s%s%s", prevPrefix, newPrefix, line)
	fmt.Println(formattedLine)
}
