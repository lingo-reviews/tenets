// Package tenet provides a library to write a tenet in Go.
//
// The tenet's API is served as an RPC plugin (github.com/lingo-
// reviews/dev/plugin) consumed by the lingo executable (github.com/lingo-
// reviews/lingo). The following is a fully functional tenet plugin (See
// github.com/lingo-reviews/tenets/go/tenets for more examples):
//
//
// package main
//
// import (
// 	"go/ast"
//
// 	"github.com/lingo-reviews/tenets/go/dev/plugin"
// 	"github.com/lingo-reviews/tenets/go/dev/tenet"
// )
//
// type commentTenet struct {
// 	tenet.Base
// }
//
// func main() {
// 	c := &commentTenet{}
// 	c.SetInfo(tenet.Info{
// 		Name:     "simpleseed",
// 		Usage:    "every comment should be awesome",
// 		Language: "golang",
// 	})
// 	c.SmellNode(func(file *tenet.File, comment *ast.Comment) error {
// 		if comment.Text != "most awesome comment ever" {
// 			c.AddNodeIssue(file, comment, 0.9)
// 		}
// 		return nil
// 	})
// 	c.AddComment("this comment could be more awesome")
// 	plugin.Serve(c)
// }
//
package tenet

// TODO(waigani) PREALPHA update this.
