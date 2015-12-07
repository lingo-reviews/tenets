package tenet_test

import (
	slasher "github.com/lingo-reviews/tenets/go/tenets/slasher/tenet"

	"testing"

	tt "github.com/lingo-reviews/tenets/go/dev/tenet/testing"
	gc "gopkg.in/check.v1"
)

func Test(t *testing.T) {
	gc.TestingT(t)
}

type slasherSuite struct {
	tt.TenetSuite
}

var _ = gc.Suite(&slasherSuite{})

func (s *slasherSuite) SetUpSuite(c *gc.C) {
	l := slasher.New()
	s.Tenet = l
}

func (s *slasherSuite) TestExampleFiles(c *gc.C) {
	files := []string{
		"example/demo.go",
	}

	expectedIssues := []tt.ExpectedIssue{
		{
			Filename: "example/demo.go",
			Text:     "//first comment",
			Comment:  "You need a space after the '//'",
			Metrics:  map[string]interface{}{"confidence": 0.9},
		},
		{
			Filename: "example/demo.go",
			Text:     "//second comment",
			Comment:  "Here needs a space also.",
			Metrics:  map[string]interface{}{"confidence": 0.9},
		}, {
			Filename: "example/demo.go",
			Text:     "//third comment",
			Comment:  "And so on, please always have a space.",
			Metrics:  map[string]interface{}{"confidence": 0.9},
		},
	}

	s.CheckFiles(c, files, expectedIssues...)
}
