// Copyright 2015 Jesse Meek.
// Licensed under the AGPLv3, see LICENCE file for details.

package tenet_test

import (
	"testing"

	tt "github.com/lingo-reviews/tenets/go/dev/tenet/testing"
	gc "gopkg.in/check.v1"

	imports "github.com/lingo-reviews/tenets/go/tenets/imports/tenet"
)

func Test(t *testing.T) {
	gc.TestingT(t)
}

type importsSuite struct {
	tt.TenetSuite
}

var _ = gc.Suite(&importsSuite{})

func (s *importsSuite) SetUpTest(c *gc.C) {
	s.Tenet = imports.New()
	s.TenetSuite.SetUpTest(c)
}

func (s *importsSuite) TestExampleFiles(c *gc.C) {

	s.SetCfgOption(c, "blacklist_regex", ".*/state")

	files := []string{
		"example/worker.go",
	}

	expectedIssues := []tt.ExpectedIssue{
		{
			Filename: "example/worker.go",
			Text:     "\t\"github.com/juju/juju/state\"",
			Comment:  `This package should not be bringing in "github.com/juju/juju/state"`,
		},
	}

	s.CheckFiles(c, files, expectedIssues...)
}
