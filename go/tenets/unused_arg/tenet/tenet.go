package tenet

import (
	"go/ast"
	"strings"

	"github.com/lingo-reviews/tenets/go/dev/tenet"
)

type unusedArgTenet struct {
	tenet.Base
}

func New() *unusedArgTenet {
	t := &unusedArgTenet{}
	t.SetInfo(tenet.Info{
		Name:        "unused_arg",
		Usage:       "catch funcs that don't use an argument",
		Description: "Ensure a function's arguments are used in the body of the function.",
		SearchTags:  []string{"function", "method"},
		Language:    "golang",
	})

	confidence := t.RegisterMetric("confidence")
	issue := t.RegisterIssue("unused_func_arg",
		// TODO(waigani) "Argument{s} {args} {is_are} not used in the function's body."
		tenet.AddComment(`{{.args}} used`),
	)

	t.SmellNode(func(r tenet.Review, fnc *ast.FuncDecl) error {
		args := fnc.Type.Params.List
		if len(args) == 0 {
			return nil
		}

		v := &visitor{
			args: map[string]bool{},
		}
		for _, arg := range args {
			for _, ident := range arg.Names {
				v.args[ident.Name] = true
			}
		}
		ast.Walk(v, fnc.Body)
		if len(v.args) > 0 {
			names := make([]string, 0, len(v.args))
			for k := range v.args {
				names = append(names, k)
			}
			unused := `"` + strings.Join(names, `", "`)

			if len(v.args) > 1 {
				unused += `" aren't`
			} else {
				unused += `" isn't`
			}

			r.RaiseNodeIssue(issue, fnc, confidence(0.5), tenet.CommentVar("args", unused))
		}

		return nil
	})
	return t
}

type visitor struct {
	args map[string]bool
}

func (v *visitor) Visit(n ast.Node) ast.Visitor {
	if ident, ok := n.(*ast.Ident); ok {
		if v.args[ident.Name] {
			delete(v.args, ident.Name)
		}
	}

	return v
}
