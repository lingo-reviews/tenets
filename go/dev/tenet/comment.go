// Copyright 2015 Jesse Meek.
// Licensed under the AGPLv3, see LICENCE file for details.

package tenet

// CommentContext is a bit operator which dictates when the comment should be
// used. e.g. first, second, last time an issue is encountered.
type CommentContext int

const (
	DefaultComment CommentContext = 1 << iota
	FirstComment
	SecondComment
	ThirdComment
	FourthComment
	FifthComment
	FirstCommentInFile
	SecondCommentInFile
	ThirdCommentInFile
	FourthCommentInFile
	FifthCommentInFile
)

type comment struct {
	// unique id per commentSet.
	ID int

	// string templates.
	Template string

	// the context in which this comment should be used.
	Context CommentContext

	// Has an issue been raised in a context that matches this comment's
	// context? Note, even if true, this comment may not be used if another
	// comment also matched the context. We record matches to know when we
	// have found issues for all contexts and can stop reviewing.
	Matched bool
}

type commentSet struct {
	Comments []*comment

	// comment id incrementor
	idInc int
}

func (c *commentSet) AddComment(commentTemplate string, context ...CommentContext) {
	var finalCtx CommentContext
	if len(context) == 0 {
		context = []CommentContext{DefaultComment}
	}
	for _, ctx := range context {
		finalCtx |= ctx
	}

	c.idInc++
	c.Comments = append(c.Comments, &comment{
		ID:       c.idInc,
		Template: commentTemplate,
		Context:  finalCtx,
	})
}

func (c *commentSet) commentsForContext(con CommentContext) []*comment {
	var comments []*comment
	for _, comm := range c.Comments {
		if comm.Context&con != 0 {
			comments = append(comments, comm)
		}
	}

	return comments
}
