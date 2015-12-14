package tenet_test

import (
	"testing"

	tt "github.com/lingo-reviews/tenets/go/dev/tenet/testing"
	gc "gopkg.in/check.v1"

	unusedArg "github.com/lingo-reviews/tenets/go/tenets/unused_arg/tenet"
)

func Test(t *testing.T) {
	gc.TestingT(t)
}

type unusedArgSuite struct {
	tt.TenetSuite
}

var _ = gc.Suite(&unusedArgSuite{})

func (s *unusedArgSuite) SetUpSuite(c *gc.C) {
	l := unusedArg.New()
	s.Tenet = l
}

func (s *unusedArgSuite) TestExampleFiles(c *gc.C) {
	files := []string{
		"example/example.go",
	}

	metrics := map[string]interface{}{"confidence": 0.5}

	expectedIssues := []tt.ExpectedIssue{
		{
			Filename: "example/example.go",
			Text:     "func saySomething(something string) {\n\tfmt.Println(\"hi\")",
			Comment:  `"something" isn't used`,
			Metrics:  metrics,
		}, {
			Filename: "example/example.go",
			Text:     "func saySomethingOther(something, otherthing string) {\n\tfmt.Println(\"hi\")",
			Comment:  `"something", "otherthing" aren't used`,
			Metrics:  metrics,
		}, {
			Filename: "example/example.go",
			Text:     "func saySomethingElse(something, otherthing string) {\n\tfmt.Println(something)",
			Comment:  `"otherthing" isn't used`,
			Metrics:  metrics,
		},
	}

	s.CheckFiles(c, files, expectedIssues...)
}
