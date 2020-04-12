package main

import (
	"regexp"
	"strings"

	"github.com/mitchellh/go-wordwrap"
)

type formatter struct {
	Substitute string
	WrapLength int
}

func formatParagraphs(paragraphs []string, options formatter) string {
	fmted := make([]string, len(paragraphs), len(paragraphs))
	for i, p := range paragraphs {
		var s string
		if options.Substitute != "" {
			s = strings.ReplaceAll(p, "%s", options.Substitute)
		} else {
			s = p
		}
		fmted[i] = wordwrap.WrapString(
			strings.TrimSpace(
				strings.ReplaceAll(
					regexp.MustCompile(`\s+`).ReplaceAllString(s, " "),
					"\n", " ")),
			80)
	}
	return strings.Join(fmted, "\n\n") + "\n"
}
