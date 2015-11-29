// Copyright 2015 Jesse Meek.
// Licensed under the AGPLv3, see LICENCE file for details.

// TODO(waigani) These tests demonstrate the use of several test helpers.
// Do the same when writing tenet/seed.

package tenet_test

import (
	"testing"

	"github.com/lingo-reviews/tenets/go/dev/tenet"

	// TODO(matt, waigani) I've ended up calling a lot of packages "tenet".
	// This will lead to confusion. Once the dust settles, let's think of some
	// sane naming.
	license "github.com/lingo-reviews/tenets/go/tenets/license/tenet"

	tt "github.com/lingo-reviews/tenets/go/dev/tenet/testing"
	gc "gopkg.in/check.v1"
)

const (
	comment1 = `Each file should start with a license header of the following format:

// Copyright 2015 Jesse Meek.
// Licensed under the AGPLv3, see LICENCE file for details.

`
	comment2 = "This file also needs the correctly formatted license header."
	comment3 = "And so on, license header again."
)

func Test(t *testing.T) {
	gc.TestingT(t)
}

type licenseSuite struct {
	tt.TenetSuite
}

var _ = gc.Suite(&licenseSuite{})

func (s *licenseSuite) SetUpTest(c *gc.C) {
	s.Tenet = license.New()
	s.TenetSuite.SetUpTest(c)

	s.SetCfgOption(c, "header", `
// Copyright 2015 Jesse Meek.
// Licensed under the AGPLv3, see LICENCE file for details.
`[1:])

}

func (s *licenseSuite) TestInfo(c *gc.C) {
	t := s.Tenet
	expected := &tenet.Info{
		Name:        "license",
		Usage:       "Each file should contain the appropriate license header.",
		Description: `Ensure that each file in the project begins with a license: "{{.header}}"`,
		SearchTags:  []string{"license", "comment", "doc-comment"},
		Language:    "go",
	}
	i := t.Info()
	c.Assert(i.Name, gc.Equals, expected.Name)
	c.Assert(i.Usage, gc.Equals, expected.Usage)
	c.Assert(i.Description, gc.Equals, expected.Description)
	c.Assert(i.SearchTags, gc.DeepEquals, expected.SearchTags)
	c.Assert(i.Language, gc.Equals, expected.Language)
	c.Assert(i.Options, gc.HasLen, 1)
}

func (s *licenseSuite) TestExampleFiles(c *gc.C) {

	files := []string{
		"example/file1.go",
		"example/file2.go",
		"example/file3.go",
	}

	expectedIssues := []tt.ExpectedIssue{
		{
			Filename: "example/file1.go",
			Text:     "// file1 this should be a license header.",
			Comment:  comment1,
		},
		{
			Filename: "example/file2.go",
			Text:     "// file2 this should be a license header.",
			Comment:  comment2,
		}, {
			Filename: "example/file3.go",
			Text:     "// file3 this should be a license header.",
			Comment:  comment3,
		},
	}

	s.CheckFiles(c, files, expectedIssues...)
}

func (s *licenseSuite) TestSRC(c *gc.C) {
	src := `
// Copyright 2015 Jesse Meek.
// Licensed under the AGPLv3, see LICENCE file for etails.

// package comment
package main
`[1:]

	expected := tt.ExpectedIssue{
		Text:    "// Copyright 2015 Jesse Meek.",
		Comment: comment1,
	}

	s.CheckSRC(c, src, expected)
}

func (s *licenseSuite) TestIssues(c *gc.C) {
	src := `
// not a license header

package main
`[1:]

	s.CheckSRC(c, src, tt.ExpectedIssue{
		Text:    "// not a license header",
		Comment: comment1,
	})
}
