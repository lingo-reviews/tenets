# Write a Go Tenet!

A tenet is a micro-service which talks to lingo over RPC. In this directory you'll find:


### api/

The RPC api. This is the low level transport code generated from api.proto which enables lingo to talk to your tenet. You can basically ignore it.

### server/

There is only one method from the server package that you should need:

```go
server.Serve(t tenet.Tenet)
```

This should be called at the end of your main function. It takes a tenet object and serves it up as an RPC server that lingo can talk to.

### tenet/

The tenet package is a helper library for writing a tenet. Here is the minimum code to get a tenet running:

```go
package main

import (
	"github.com/lingo-reviews/tenets/go/dev/server"
	"github.com/lingo-reviews/tenets/go/dev/tenet"
)

func main() {
server.Serve(&tenet.Base{})
}
```

You could build and add this tenet to lingo, but it would do nothing. We should define what this tenet's all about:

```go
t  := &tenet.Base{}

	t.SetInfo(tenet.Info{
		// This information will be shown in `lingo info <tenet-name>`
		Name:     "example_tenet",
		Usage:    "every comment should be awesome",
		Language: "Go",
		
		// Tenets are about documenting rules as much as catching violations
		// of them. This description gets composed with that of other tenets
		// to document a package.
		Description: `
simpleseed is a demonstration tenet showing the structure required
to write a tenet in Go. When reviewing code with simpleseed it will be suggested
that all comments could be more awesome.
`,
	})


server.Serve(t)
```

Let's make it smell all comments in a go file:

```go
t  := &tenet.Base{}

 ...

	// First register the issue that this tenet will look for. It returns the
	// name of the issue which you'll use to raise the issue after smelling
	// the comment.
	// You can register as many issues as you like.
	issue := t.RegisterIssue("sucky_comment")

	// You can smell as many nodes as you like.
	t.SmellNode(func(r tenet.Review, commentNode *ast.Comment) error {
		if comment.Text != "// most awesome comment ever" {
			r.RaiseNodeIssue(issue, commentNode)
		}
		return nil
	})

server.Serve(t)
```

This will raise an issue for every non-awesome comment, with the default message "Issue Found". To set the message:

```go
t.RegisterIssue("sucky_comment", tenet.AddComment("this comment could be more awesome"))
```

What if we don't want to raise every issue, but just enough to point out the problem?

```go

issue := t.RegisterIssue("sucky_comment",

		// You can add as may comments as you like. Though note only the
		// comment matching the context will be used.
		tenet.AddComment("comments really should be awesome", tenet.FirstComment),
		tenet.AddComment("the comment in this file should also be more awesome", tenet.FirstCommentInFile),
	)

```

The first time we see the issue, lingo will comment "comments really should be awesome". Then, once for every file the issue is found in lingo will comment "the comment in this file should also be more awesome".

Let's set a variable in the comment:

```go
issue := t.RegisterIssue("sucky_comment",
		tenet.AddComment("comments really should be {{.myvar}}"),
		)

// then in our smell
r.RaiseNodeIssue(issue, commentNode, tenet.CommentVar({{.myvar}}, "awesome"))

```

Let's get that variable from the user:

```go
	// t.RegisterOption(name, default, usage)
	commentType := t.RegisterOption("comment_type", "awesome", "set the type of comment")

	// then in our smell
	r.RaiseNodeIssue(issue, commentNode, tenet.CommentVar({{.myvar}}, *commentType))	

```

When the user runs `lingo info <tenet>` they'll see "comment_type" as an option they can set. t.RegisterOption returns a pointer to a string. The value of that pointer is updated with the user's setting by the time it is used in the smell.

What if we've found an smell, but are not sure if it's really an issue? You can register custom metrics an tags to handle this.


```go
	confidence := t.RegisterMetric("confidence")
	style := t.RegisterTag("style")

	// then in our smell
	r.RaiseNodeIssue(issue, commentNode, style, confidence(8))

```

Now the issue has been raised with a confidence score of 8 and a tag of "style". To use these when reviewing:

```bash
lingo review --tags style,someOtherTag --metrics-higher-than confidence=5 --metrics-lower-than confidence=9
```

This enables lingo to monitor a code base with fine grained control. You could, for example, encode the connascence princples (http://connascence.io).

NOTE: In the closed Alpha you can register and raise metrics and tags - and they'll appear in the json output - but you cannot yet filter a review with them.