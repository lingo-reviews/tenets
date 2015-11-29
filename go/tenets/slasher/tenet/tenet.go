package tenet

import (
	"regexp"

	"github.com/lingo-reviews/tenets/go/dev/tenet"
)

type slasherTenet struct {
	tenet.Base
}

func New() *slasherTenet {
	t := &slasherTenet{}
	t.SetInfo(tenet.Info{
		Name:        "slasher",
		Usage:       "Comments should have a space after '//' forward slashes",
		Description: "Comments should have a space after '//' forward slashes",
		SearchTags:  []string{"format", "comment", "doc-comment"},
		Language:    "golang",
	})

	confidence := t.RegisterMetric("confidence")

	issue := t.RegisterIssue("no_space_after_comment",
		tenet.AddComment("You need a space after the '//'", tenet.FirstComment),
		tenet.AddComment("Here needs a space also.", tenet.SecondComment),
		tenet.AddComment("And so on, please always have a space.", tenet.ThirdComment),
	)

	t.SmellLine(func(r tenet.Review, n int, line []byte) error {
		if regexp.MustCompile(`\/\/[^\s]{1}`).Match(line) {
			r.RaiseLineIssue(issue, n, n, confidence(0.9))
		}
		return nil
	})

	return t
}
