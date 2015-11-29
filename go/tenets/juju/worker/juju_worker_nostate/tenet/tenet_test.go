// Copyright 2015 Jesse Meek.
// Licensed under the AGPLv3, see LICENCE file for details.

package tenet_test

import (
	"testing"

	jc "github.com/juju/testing/checkers"
	nostate "github.com/lingo-reviews/tenets/go/tenets/juju/worker/juju_worker_nostate/tenet"
	gc "gopkg.in/check.v1"

	"github.com/lingo-reviews/tenets/go/dev/tenet"
	tt "github.com/lingo-reviews/tenets/go/dev/tenet/testing"
)

func Test(t *testing.T) {
	gc.TestingT(t)
}

type noStateSuite struct {
	tt.TenetSuite
}

var _ = gc.Suite(&noStateSuite{})

func (s *noStateSuite) SetUpSuite(c *gc.C) {
	t := nostate.New()

	s.Tenet = t
}

func (s *noStateSuite) TestInfo(c *gc.C) {
	t := s.Tenet
	expected := tenet.Info{
		Language:    "golang",
		Name:        "worker_nostate",
		Usage:       "workers should not access state directly",
		SearchTags:  []string{"juju", "worker"},
		Description: "If you're passing a \\*state.State into your worker, you are almost certainly doing it wrong. The layers go worker->apiserver->state, and any attempt to skip past the apiserver layer should be viewed with *extreme* suspicion.",
	}
	c.Assert(*t.Info(), jc.DeepEquals, expected)
}

func (s *noStateSuite) TestExampleFiles(c *gc.C) {

	files := []string{
		"example/worker.go",
		"example/worker2.go",
		"example/worker3.go",
		"example/worker4.go",
	}

	expectedIssues := []tt.ExpectedIssue{
		{
			Filename: "example/worker.go",
			Text:     "func New(st *state.State, params *HistoryPrunerParams) worker.Worker {",
			Comment: `
I see you've imported state. A worker shouldn't need it. Best practice for writing workers: 
https://github.com/juju/juju/wiki/Guidelines-for-writing-workers
`[1:],
		},
		{
			Filename: "example/worker2.go",
			Text:     "func New(params *HistoryPrunerParams) worker.Worker {",
			Comment: `
Please don't pass in a state object here. Workers should use the API.
More info here: https://github.com/juju/juju/wiki/Guidelines-for-writing-workers
`[1:],
		},
	}

	s.CheckFiles(c, files, expectedIssues...)
}
