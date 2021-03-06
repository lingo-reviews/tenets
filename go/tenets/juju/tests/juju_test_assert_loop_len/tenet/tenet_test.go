// Copyright 2015 Jesse Meek.
// Licensed under the AGPLv3, see LICENCE file for details.

// TODO(waigani) These tests demonstrate the use of several test helpers.
// Do the same when writing tenet/seed.

package tenet_test

import (
	"testing"

	// TODO(matt, waigani) I've ended up calling a lot of packages "tenet".
	// This will lead to confusion. Once the dust settles, let's think of some
	// sane naming.
	loop "github.com/lingo-reviews/tenets/go/tenets/juju/tests/juju_test_assert_loop_len/tenet"

	tt "github.com/lingo-reviews/tenets/go/dev/tenet/testing"
	gc "gopkg.in/check.v1"
)

func Test(t *testing.T) {
	gc.TestingT(t)
}

type suite struct {
	tt.TenetSuite
}

var _ = gc.Suite(&suite{})

func (s *suite) SetUpTest(c *gc.C) {
	s.Tenet = loop.New()
	s.TenetSuite.SetUpTest(c)
}

func (s *suite) TestForLoop(c *gc.C) {

	files := []string{
		"example/bad_test.go",
	}

	expectedIssues := []tt.ExpectedIssue{
		{
			Filename: "example/bad_test.go",
			Text:     "\tfor _, s := range list() {",
			Comment: `
Even if you assert the length of the result of this call before iterating
over it, you cannot guarantee the result will be the same each time you call
it. You cannot be sure that the asserts within the for loop will get
run. Please assign the result of list() to a variable, assert the expected
length of the variable and then loop over that.`[1:],
		}, {
			Filename: "example/bad_test.go",
			Text:     "\tfor _, cont := range strings.Split(string(content), \"\\n\") {",
			Comment:  `Here also, you can't gurantee the length of Split()'s result.`,
		}, {
			Filename: "example/bad_test.go",
			Text:     "\tfor _, s := range a {",
			Comment:  `The asserts within this loop may never get run. The length of "a" needs to be asserted.`,
		},
		{
			Filename: "example/bad_test.go",
			Text:     "\tfor _, s := range a {",
			Comment:  `The length of "a" needs to be asserted also`,
		}, {
			Filename: "example/bad_test.go",
			Text:     "\tfor _, s := range a {",
			Comment:  `The length of "a" needs to be asserted also`,
		}, {
			Filename: "example/bad_test.go",
			Text:     "\tfor _, s := range a {",
			Comment:  `The length of "a" needs to be asserted also`,
		},
	}

	s.CheckFiles(c, files, expectedIssues...)
}
