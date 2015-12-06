// Copyright 2015 Jesse Meek.
// Licensed under the AGPLv3, see LICENCE file for details.

package tenet

import (
	"go/ast"
	"strings"

	"github.com/lingo-reviews/tenets/go/dev/tenet"
)

type assertLoopLenTenet struct {
	tenet.Base
}

func New() *assertLoopLenTenet {
	t := &assertLoopLenTenet{}
	t.SetInfo(tenet.Info{
		Name:        "juju_test_assert_loop_len",
		Usage:       "If asserting within a loop, the length of the colleciton being iterated should be asserted",
		Description: "If asserting within a loop, the length of the colleciton being iterated should be asserted",
		SearchTags:  []string{"test", "loop"},
		Language:    "go",
	})

	assertLoopIssue := t.RegisterIssue("loop_len_not_asserted",
		tenet.AddComment("The asserts within this loop may never get run. The length of the collection being looped needs to be asserted.", tenet.FirstComment),
		tenet.AddComment("The length of this loop needs to be asserted also", tenet.DefaultComment),
	)

	// TODO(waigani) good canditate for a patch fix.
	rangeCallExpIssue := t.RegisterIssue("range_over_call_exp",
		tenet.AddComment(`
Even if you assert the length of the result of this call before iterating
over it, you cannot guarantee the result will be the same each time you call
it. Thus, you cannot be sure that the asserts within the for loop will get
run. Please assign the result of the call to a variable, assert the expected
length of the variable and then loop over it.`[1:], tenet.FirstComment),
		tenet.AddComment(`
Here also, you can't gurantee the length of the call's result.`[1:], tenet.SecondComment),
		tenet.AddComment(`
Don't loop over call result.`[1:], tenet.DefaultComment),
	)

	// // First, knock out any file that isn't a test
	t.SmellNode(func(r tenet.Review, _ *ast.File) error {
		if !strings.HasSuffix(r.File().Filename(), "_test.go") {
			r.FileDone()
		}
		return nil
	})

	// Now, the only nodes smelt are within test files.
	// Build up a list of all collections that have an asserted len
	var assertLen []string
	t.SmellNode(func(r tenet.Review, callExpr *ast.CallExpr) error {
		if fun, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
			if fun.Sel.String() != "Assert" && fun.Sel.String() != "Check" {
				return nil
			}
			if args := callExpr.Args; len(args) == 3 {
				if ident, ok := args[0].(*ast.Ident); ok {
					// we have an assert on the collection.
					if sel, ok := args[1].(*ast.SelectorExpr); ok {
						if sel.Sel.String() == "HasLen" {
							assertLen = append(assertLen, ident.String())
							return nil
						}
					}

				}

			}
		}
		return nil
	})

	// Check if any range body contains an assert or check and the
	// collection ranged over has not had an asserted length.
	t.SmellNode(func(r tenet.Review, rngStmt *ast.RangeStmt) error {
		if containsCheckOrAssert(rngStmt.Body.List) {
			switch n := rngStmt.X.(type) {
			case *ast.CallExpr:

				r.RaiseNodeIssue(rangeCallExpIssue, n)
			case *ast.Ident:
				var checked bool
				for _, asserted := range assertLen {
					if n.String() == asserted {
						checked = true
					}
				}
				if !checked {
					r.RaiseNodeIssue(assertLoopIssue, n)
				}
			default:
				// TODO(waigani) log unknown range symbol
			}

		}
		return nil
	})

	return t
}

// TODO(waigani) check for assetCustom(c) type funcs within the loop
func containsCheckOrAssert(stmts []ast.Stmt) bool {
	for _, stmt := range stmts {
		switch exp := stmt.(type) {
		case *ast.ExprStmt:
			switch n := exp.X.(type) {
			case (*ast.CallExpr):

				switch x := n.Fun.(type) {
				case (*ast.SelectorExpr):

					for _, sel := range []string{
						"Assert",
						"Check",
					} {
						if x.Sel.String() == sel {
							return true

						}
					}
				}

			}
		}
	}

	return false
}
