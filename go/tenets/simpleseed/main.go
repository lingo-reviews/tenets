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
		Name:     "simpleseed",
		Usage:    "every comment should be awesome",
		Language: "Go",
		Description: `
simpleseed is a demonstration tenet showing the structure required
to write a tenet in Go. When reviewing code with simpleseed it will be suggested
that all comments could be more awesome.
`,
	})

	issue := c.RegisterIssue("sucky_comment", tenet.AddComment("this comment could be more awesome", tenet.FirstComment))
	c.SmellNode(func(r tenet.Review, comment *ast.Comment) error {
		if comment.Text != "// most awesome comment ever" {
			r.RaiseNodeIssue(issue, comment)
		}
		return nil
	})

	server.Serve(c)
}
