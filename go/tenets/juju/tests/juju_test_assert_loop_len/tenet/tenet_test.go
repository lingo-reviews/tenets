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
		"example/network_test.go",
	}

	expectedIssues := []tt.ExpectedIssue{
		{
			Filename: "example/network_test.go",
			Text:     "\tfor p, aval := range server.Ports() {",
			Comment: `
Even if you assert the length of the result of this call before iterating
over it, you cannot guarantee the result will be the same each time you call
it. Thus, you cannot be sure that the asserts within the for loop will get
run. Please assign the result of the call to a variable, assert the expected
length of the variable and then loop over it.`[1:]}, {
			Filename: "example/network_test.go",
			Text:     "\tfor _, p := range portSlice {",
			Comment:  "The asserts within this loop may never get run. The length of the collection being looped needs to be asserted.",
		}, {
			Filename: "example/network_test.go",
			Text:     "\tfor _, aval := range server.Ports() {",
			Comment:  "Here also, you can't gurantee the length of the call's result.",
		},
	}

	s.CheckFiles(c, files, expectedIssues...)
}
