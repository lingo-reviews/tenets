package tenet

import (
	"go/ast"
	"strings"

	"github.com/juju/errors"
	"github.com/lingo-reviews/tenets/go/dev/tenet"
	"github.com/lingo-reviews/tenets/go/dev/tenet/astutil"
)

type unusedArgTenet struct {
	tenet.Base
}

func New() *unusedArgTenet {
	t := &unusedArgTenet{}
	t.SetInfo(tenet.Info{
		Name:        "exhaustive_switch",
		Usage:       "ensure that a switch has a case for every possible value of a variable",
		Description: "Ensure that a switch has a case for every possible value of a variable",
		SearchTags:  []string{"switch"},
		Language:    "golang",
	})

	exhaustTag := t.RegisterOption("exhaust_tag", "lingo:exhaustive", "If this tag is found in a comment above a switch statement, that switch will be treated as an exhaustive switch.")
	issue := t.RegisterIssue("missing_case",
		tenet.AddComment(`The following cases are missing from this switch:{{.cases}}.`),
	)

	// map of a switched variable to all case values.
	switchedVarSets := map[int]map[string][]ast.Expr{}

	t.SmellNode(func(r tenet.Review, file *ast.File) error {

		for _, com := range file.Comments {
			if strings.Contains(com.Text(), *exhaustTag) {
				switchedVarSets[r.File().Fset().Position(com.End()).Line+1] = map[string][]ast.Expr{}
			}
		}

		return nil
	})

	// First find switched vars
	t.SmellNode(func(r tenet.Review, swt *ast.SwitchStmt) error {

		// Did we find a tag for this switch?
		swtLine := r.File().Fset().Position(swt.Pos()).Line
		if _, ok := switchedVarSets[swtLine]; !ok {
			return nil
		}

		var switchedVar string
		if ident, ok := swt.Tag.(*ast.Ident); ok {

			var err error
			switchedVar, err = astutil.TypeOf(ident)
			if err != nil {
				return errors.Trace(err)
			}
		} else {
			return errors.Errorf("found switch, but could not get variable name: %#v", swt.Tag)
		}

		for _, stm := range swt.Body.List {
			if c, ok := stm.(*ast.CaseClause); ok {
				for _, l := range c.List {

					switchedVarSets[swtLine][switchedVar] = append(switchedVarSets[swtLine][switchedVar], l)
				}
			}
		}

		return nil
	})

	missingCases := map[int][]*ast.Ident{}

	// Then find genDecls with the switched type in it and make sure all
	// GenDecls of that type are switched on.
	t.SmellNode(func(r tenet.Review, genDec *ast.GenDecl) error {

		if len(switchedVarSets) == 0 {
			r.FileDone()
		}

		for switchLine, switchSet := range switchedVarSets {
			for switchType := range switchSet {

				var inGenDecl bool
				// first see if a switched type is in this genDecl
				for _, s := range genDec.Specs {

					if valSpec, ok := s.(*ast.ValueSpec); ok {

						for _, n := range valSpec.Names {

							typ, err := astutil.TypeOf(n)
							if err != nil {
								// TODO(waigani) log
								continue
							}

							// Does the type in this GenDecl match that of the
							// switch?
							if typ == switchType {
								inGenDecl = true
							}

						}

					}
				}

				if !inGenDecl {
					continue
				}

				// Now make sure there is a case for each GenGecl
				for _, s := range genDec.Specs {
					if valSpec, ok := s.(*ast.ValueSpec); ok {
						for _, vName := range valSpec.Names {

							// Does our switch contain the var from  genDecl?
							if !contains(vName, switchedVarSets[switchLine][switchType]) {
								missingCases[switchLine] = append(missingCases[switchLine], vName)
							}
						}
					}
				}
			}
		}

		return nil
	})

	// Do another sweep over switches, this time raising an issue for any in the missingCases map
	t.SmellNode(func(r tenet.Review, swt *ast.SwitchStmt) error {

		if len(missingCases) == 0 {
			// no switches with missing cases were found. Don't smell any more
			// nodes.
			r.SmellDone()
			return nil
		}

		switchLine := r.File().Fset().Position(swt.Pos()).Line

		if cases, ok := missingCases[switchLine]; ok {
			var missing string
			for _, m := range cases {
				missing += ", " + m.String()
			}

			r.RaiseNodeIssue(issue, swt, tenet.CommentVar("cases", missing[1:]))
		}
		return nil
	})

	return t
}

func contains(needle *ast.Ident, haystack []ast.Expr) bool {
	for _, x := range haystack {
		if ident, ok := x.(*ast.Ident); ok {

			// TODO(waigani) return on more than name.
			if ident.Name == needle.Name {
				return true
			}
		}
	}
	return false
}
