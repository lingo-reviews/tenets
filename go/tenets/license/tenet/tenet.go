// Copyright 2015 Jesse Meek.
// Licensed under the AGPLv3, see LICENCE file for details.

package tenet

import (
	"fmt"
	"strings"

	"github.com/lingo-reviews/tenets/go/dev/tenet"
)

type licenseTenet struct {
	tenet.Base
}

func New() *licenseTenet {
	t := &licenseTenet{}
	t.SetInfo(tenet.Info{
		Name:        "license",
		Usage:       "Each file should contain the appropriate license header.",
		Description: "Ensure that each file in the project begins with a license: \"{{.header}}\"",
		SearchTags:  []string{"license", "comment", "doc-comment"},
		Language:    "go",
	})

	confidence := t.RegisterMetric("confidence")
	issue := t.RegisterIssue("incorrect_header", tenet.AddComment(`
Each file should start with a license header of the following format:

{{.header}}
`[1:], tenet.FirstComment),
		tenet.AddComment("This file also needs the correctly formatted license header.", tenet.SecondComment),
		tenet.AddComment("And so on, license header again.", tenet.ThirdComment),
	)

	// RegisterOption returns a pointer to a string. The string will be
	// updated by the time it is used in the smell.
	headerPnt := t.RegisterOption("header", "Copyright Me", "the license header to check for")

	t.SmellLine(func(r tenet.Review, n int, line []byte) error {
		headerVal := *headerPnt
		lineMatchesHeader := func() bool {
			var i int
			header := strings.Split(headerVal, "\n")
			for i = 1; i <= len(header); i++ {
				if n == i && string(line) != header[i-1] {
					// It doesn't match, no need to keep sniffing this file.
					r.SmellDoneWithFile()
					return false
				}
			}
			if n == i {
				fmt.Println(n, i, string(line), headerVal)
				// Every line of the header matches. Don't check again.
				r.SmellDoneWithFile()
			}
			return true
		}

		if !lineMatchesHeader() {
			r.RaiseLineIssue(issue, 1, 1, confidence(0.7), tenet.CommentVar("header", headerVal))
		}

		return nil
	})

	return t
}
