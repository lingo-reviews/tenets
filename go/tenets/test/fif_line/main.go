package main

import (
	"regexp"

	"github.com/lingo-reviews/tenets/go/dev/server"
	"github.com/lingo-reviews/tenets/go/dev/tenet"
)

type commentTenet struct {
	tenet.Base
}

// a comment
func main() {
	c := &commentTenet{}
	c.SetInfo(tenet.Info{
		Name:     "fif_line",
		Usage:    "take issue with the first comment in every file",
		Language: "Go",
		Description: `
This tenet is part of lingo's internal testing suite. It should raise an issue for the
first comment encountered in every file, using SmellLine.
`,
	})

	issue := c.RegisterIssue("comment", tenet.AddComment("first comment - line", tenet.FirstComment, tenet.InEveryFile))
	c.SmellLine(func(r tenet.Review, n int, line []byte) error {
		if regexp.MustCompile(`\/\/.*`).Match(line) {
			r.RaiseLineIssue(issue, n, n)
		}
		return nil
	})

	server.Serve(c)
}
