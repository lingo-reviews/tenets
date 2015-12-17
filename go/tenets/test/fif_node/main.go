package main

import (
	"go/ast"

	"github.com/lingo-reviews/tenets/go/dev/server"
	"github.com/lingo-reviews/tenets/go/dev/tenet"
)

type commentTenet struct {
	tenet.Base
}

// a comment
func main() {
	c := &commentTenet{}
	c.SetInfo(tenet.Info{
		Name:     "fif_node",
		Usage:    "take issue with the first comment in every file",
		Language: "Go",
		Description: `
This tenet is part of lingo's internal testing suite. It should raise an issue for the
first comment encountered in every file, using SmellNode.
`,
	})

	issue := c.RegisterIssue("comment", tenet.AddComment("first comment - node", tenet.FirstComment, tenet.InEveryFile))
	c.SmellNode(func(r tenet.Review, comment *ast.Comment) error {
		r.RaiseNodeIssue(issue, comment)
		return nil
	})

	server.Serve(c)
}
