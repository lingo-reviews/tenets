package tenet

import (
	"go/ast"

	"golang.org/x/tools/go/ast/astutil"

	"github.com/lingo-reviews/tenets/go/dev/tenet"
)

type noStateTenet struct {
	tenet.Base
	maybeIssue      string
	mostLikelyIssue string
	confidence      func(interface{}) tenet.RaiseIssueOption
}

func New() *noStateTenet {
	t := &noStateTenet{}
	t.SetInfo(tenet.Info{
		Name:        "worker_nostate",
		Usage:       "workers should not access state directly",
		Description: "If you're passing a \\*state.State into your worker, you are almost certainly doing it wrong. The layers go worker->apiserver->state, and any attempt to skip past the apiserver layer should be viewed with *extreme* suspicion.",
		Language:    "golang",
		SearchTags:  []string{"juju", "worker"},
	})

	// We register any issues, metrics and tags that we'll be using.
	t.confidence = t.RegisterMetric("confidence")
	t.RegisterTag("observablity")

	t.maybeIssue = t.RegisterIssue("imports_state_returns_worker",
		tenet.AddComment(`
I see you've imported state. A worker shouldn't need it. Best practice for writing workers: 
https://github.com/juju/juju/wiki/Guidelines-for-writing-workers
`[1:]),
	)

	t.mostLikelyIssue = t.RegisterIssue("func_takes_state_returns_worker",
		tenet.AddComment(`
Please don't pass in a state object here. Workers should use the API.
More info here: https://github.com/juju/juju/wiki/Guidelines-for-writing-workers
`[1:]),
	)

	// First, let's knock out any file that doesn't import state and worker.
	t.SmellNode(func(r tenet.Review, astFile *ast.File) error {
		if !astutil.UsesImport(astFile, "github.com/juju/juju/state") ||
			!astutil.UsesImport(astFile, "github.com/juju/juju/worker") {
			// This file will no longer be smelt by this tenet.
			r.FileDone()
		} else {
		}
		return nil
	})

	t.smellFuncs()

	return t
}

func (t *noStateTenet) smellFuncs() {
	t.SmellNode(func(r tenet.Review, fnc *ast.FuncType) error {
		// We are only interested in funcs with return values.
		if fnc.Results == nil {
			return nil
		}

		// The name of the issue we're going to raise.
		var maybeFound bool

		// Does the func return a worker?
		for _, field := range fnc.Results.List {
			if sym, ok := field.Type.(*ast.SelectorExpr); ok {
				if sym.Sel.Name == "Worker" {

					// We've got a file which imports state and a func
					// returning a worker. We've found our first code smell.
					maybeFound = true
					break
				}
			}
		}

		// If this func doesn't return a worker, we're no longer interested in it.
		if !maybeFound {
			return nil
		}

		// We've found a func returning a worker. Is it taking a pointer to a
		// State object?
		for _, field := range fnc.Params.List {
			if ptr, ok := field.Type.(*ast.StarExpr); ok {
				if sym, ok := ptr.X.(*ast.SelectorExpr); ok {

					// We are just checking the slector here. You could get
					// the package name from sym.X and assert that it's
					// "state". But it may have been renamed. Unlikely, but in
					// that case the issue would be missed. mostLikelyIssue's
					// confidence is set to 80%, indicating the *rare* false
					// positive.
					if sym.Sel.Name == "State" {
						r.RaiseNodeIssue(
							t.maybeIssue,
							fnc,
							t.confidence(0.3),
							tenet.AddTag("observability"),
						)

						// No need to keep checking as we've found our issue for this func.
						return nil
					}
				}
			}
		}

		// None of the func params imported state, but it is returning a
		// worker. So raise an issue with less confidence.
		r.RaiseNodeIssue(
			t.mostLikelyIssue,
			fnc,
			t.confidence(0.8),
			// same as: tenet.SetMetric("confidence", 0.8),
			tenet.AddTag("observability"),
		)

		return nil
	})
}
