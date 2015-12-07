// Copyright 2015 Jesse Meek.
// Licensed under the AGPLv3, see LICENCE file for details.

package tenet

import (
	"errors"
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
		tenet.AddComment(`The asserts within this loop may never get run. The length of "{{.looped}}" needs to be asserted.`, tenet.FirstComment),
		tenet.AddComment(`The length of "{{.looped}}" needs to be asserted also`, tenet.DefaultComment),
	)

	// TODO(waigani) good canditate for a patch fix.
	rangeCallExpIssue := t.RegisterIssue("range_over_call_exp",
		tenet.AddComment(`
Even if you assert the length of the result of this call before iterating
over it, you cannot guarantee the result will be the same each time you call
it. You cannot be sure that the asserts within the for loop will get
run. Please assign the result of {{.looped}} to a variable, assert the expected
length of the variable and then loop over that.`[1:], tenet.FirstComment),
		tenet.AddComment(`
Here also, you can't gurantee the length of {{.looped}}'s result.`[1:], tenet.SecondComment),
		tenet.AddComment(`
Again, need to assert result of {{.looped}} first.`[1:], tenet.DefaultComment),
	)

	// // First, knock out any file that isn't a test
	t.SmellNode(func(r tenet.Review, _ *ast.File) error {
		if !strings.HasSuffix(r.File().Filename(), "_test.go") {
			r.FileDone()
		}
		return nil
	})

	// All nodes that have been asserted in a loop.
	var ranged possibleBadRange

	// Check if any range body contains an assert or check.
	t.SmellNode(func(r tenet.Review, rngStmt *ast.RangeStmt) error {

		// TODO(waigani) need to return inc statements and check they are
		// asserted after loop.
		if containsIncStmt(rngStmt.Body.List) {
			return nil
		}

		if containsCheckOrAssert(rngStmt.Body.List) {
			switch n := rngStmt.X.(type) {
			case *ast.CallExpr:
				looped := "the call"
				switch x := n.Fun.(type) {
				case *ast.Ident:
					looped = x.Name + "()"
				case *ast.SelectorExpr:
					looped = x.Sel.Name + "()"
				}
				r.RaiseNodeIssue(rangeCallExpIssue, n, tenet.CommentVar("looped", looped))
			case *ast.Ident:
				ranged.add(n)
			default:
				// TODO(waigani) log unknown range symbol
			}
		}
		return nil
	})

	// Find any idents that have been constructed in this file.
	t.SmellNode(func(r tenet.Review, ident *ast.Ident) error {

		if ranged.empty() {
			// Nothing was found to be ranged over. No need to keep smelling.
			r.FileDone()
		}

		// This ident has not been ranged over, so we're not interested in it.
		if !ranged.has(ident) {
			return nil
		}

		if declaredWithCompLit(ident) {
			// The var has been constructed in this file, so it is clear
			// to see it's length. We are no longer interested in it.
			ranged.remove(ident)
		}

		return nil
	})

	// Check that any remaining ranged vars not constructed in this file have had their lengths asserted.
	t.SmellNode(func(r tenet.Review, callExpr *ast.CallExpr) error {
		if fun, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
			if fun.Sel.String() != "Assert" && fun.Sel.String() != "Check" {
				return nil
			}
			if args := callExpr.Args; len(args) == 3 {

				switch n := args[0].(type) {
				case *ast.Ident:

					// we have an assert on the collection.
					if sel, ok := args[1].(*ast.SelectorExpr); ok {
						if sel.Sel.String() == "HasLen" {

							// The length has been asserted, wer're no longer
							// interested in this var. Remove it if we added it.
							ranged.remove(n)
							return nil
						}
					}
				case *ast.CallExpr:
					if f, ok := n.Fun.(*ast.Ident); ok {
						if f.Name == "len" {
							if x, ok := n.Args[0].(*ast.Ident); ok {
								// The length has been asserted, wer're no longer
								// interested in this var. Remove it if we added it.
								ranged.remove(x)
							}
						}
					}
				}

			}
		}
		return nil
	})

	// Do one more run over range statements.
	t.SmellNode(func(r tenet.Review, rngStmt *ast.RangeStmt) error {

		if ranged.empty() {
			// There are no possible vars of unasserted len in this file.
			// Don't run any more smells.
			r.FileDone()
		}

		// Find our assert ranges again
		if containsCheckOrAssert(rngStmt.Body.List) {
			switch n := rngStmt.X.(type) {
			case *ast.Ident:
				// we're only interested in idents this time
				if ranged.has(n) {
					r.RaiseNodeIssue(assertLoopIssue, n, tenet.CommentVar("looped", n.Name))
				}
			}
		}
		return nil
	})

	return t
}

type possibleBadRange []*ast.Ident

func (p *possibleBadRange) has(n *ast.Ident) bool {
	for _, asserted := range *p {

		// TODO(waigani) compare on more than name. We can't use pointer
		// addresses as they change per smell ast walk.
		if sameIdent(asserted, n) {
			return true
		}
	}
	return false
}

func (p *possibleBadRange) add(n *ast.Ident) {
	v := *p
	v = append(v, n)
	*p = v
}

func (p *possibleBadRange) remove(n *ast.Ident) {
	v := *p
	var nv []*ast.Ident
	for _, ident := range v {

		// TODO(waigani) compare on more than name. We can't use pointer
		// addresses as they change per smell ast walk.
		if !sameIdent(ident, n) {
			nv = append(nv, ident)
		}
	}

	*p = nv
}

func (p *possibleBadRange) empty() bool {
	return len(*p) == 0
}

func sameIdent(a, b *ast.Ident) bool {
	// TODO(waigani) Don't rely on name, it could change and still be the same
	// ident.
	if a.String() != b.String() {
		return false
	}

	pa, err := declPos(a)
	if err != nil {
		// TODO(waigani) log error
		return false
	}

	pb, err := declPos(b)
	if err != nil {
		// TODO(waigani) log error
		return false
	}

	if pa != pb {
		return false
	}

	return true
}

// returns the possition the ident was declared
func declPos(n *ast.Ident) (int, error) {
	switch t := n.Obj.Decl.(type) {
	case *ast.AssignStmt:
		return int(t.TokPos), nil
	case *ast.ValueSpec:
		if len(t.Names) == 0 {
			return 0, errors.New("decl statement has no names")
		}
		// Even if this is not the name of the ident, it is in the same decl
		// statement. We are interested in the pos of the decl statement.
		return int(t.Names[0].NamePos), nil
	default:
		// TODO(waigani) log
	}
	return 0, errors.New("could not get declaration position")
}

func declaredWithCompLit(ident *ast.Ident) bool {
	switch n := ident.Obj.Decl.(type) {
	case *ast.AssignStmt:

		// find position of ident on lhs of assignment
		var identPos int
		for i, exp := range n.Lhs {
			if a, ok := exp.(*ast.Ident); ok && a.String() == ident.String() {
				identPos = i
				break
			}
		}

		// get the rhs counterpart
		switch x := n.Rhs[identPos].(type) {
		case *ast.CompositeLit:
			return true
		case *ast.CallExpr:
			if f, ok := x.Fun.(*ast.Ident); ok {
				if f.Name == "make" {
					return true
				}
			}
		default:
			// TODO(waigani) log
		}
	case *ast.ValueSpec:

		if len(n.Values) == 0 {
			return false
		}

		// find position of ident on lhs of assignment
		var identPos int
		for i, name := range n.Names {
			if name.String() == ident.String() {
				identPos = i
			}
		}

		// get the rhs counterpart
		switch n.Values[identPos].(type) {
		case *ast.CompositeLit:
			return true
		}
	default:
		// TODO(waigani) log
	}
	return false
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
		default:
			// TODO(waigani) log
		}
	}

	return false
}

// TODO(waigani) need to check if the inc var is asserted after the loop.
func containsIncStmt(stmts []ast.Stmt) bool {
	for _, stmt := range stmts {
		switch n := stmt.(type) {
		case *ast.IncDecStmt:
			return true
		case *ast.AssignStmt:
			if n.Tok.String() == "+=" {
				return true
			}
		default:
			// TODO(waigani) log
		}
	}
	return false
}
