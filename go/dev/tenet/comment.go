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
	InFirstFile
	InSecondFile
	InThirdFile
	InFourthFile
	InFifthFile
	InEveryFile
	InOverall
)

var fileContext = map[int]CommentContext{
	1: InFirstFile,
	2: InSecondFile,
	3: InThirdFile,
	4: InFourthFile,
	5: InFifthFile,
}

var commContext = map[int]CommentContext{
	1: FirstComment,
	2: SecondComment,
	3: ThirdComment,
	4: FourthComment,
	5: FifthComment,
}

type comment struct {

	// the comment template.
	Template string

	// the context in which this comment should be used.
	commentContexts []CommentContext

	// the file context to which the commentContext is scoped.
	fileContexts []CommentContext

	// a map of each context the comment was found in.
	matches map[CommentContext]bool
}

func (c *comment) addCommentCtx(ctx CommentContext) {
	c.commentContexts = append(c.commentContexts, ctx)
}

func (c *comment) addFileCtx(ctx CommentContext) {
	c.fileContexts = append(c.fileContexts, ctx)
}

// add a context in which this comment was matched.
func (c *comment) addMatch(ctx CommentContext) {
	c.matches[ctx] = true
}

func (c *comment) allContextsMatched() bool {
	for _, ctx := range c.allContexts() {
		var matched bool
		for matchedCtx := range c.matches {

			// If we matched on a default, don't count it.
			if ctx&(DefaultComment|InEveryFile) != 0 {
				continue
			}

			if ctx&matchedCtx == ctx {
				matched = true
				break
			}
		}
		if !matched {
			return false
		}
	}
	return true
}

// returns a slice of each commentContext scoped to a file context.
func (c *comment) allContexts() []CommentContext {
	var ctxs []CommentContext
	for _, fCtx := range c.fileContexts {
		for _, commCtx := range c.commentContexts {
			ctxs = append(ctxs, fCtx|commCtx)
		}
	}
	return ctxs
}

func (c *comment) matchesContext(ctx CommentContext) bool {
	for _, commCtx := range c.allContexts() {
		if commCtx&ctx == commCtx {
			return true
		}
	}
	return false
}

func isFileContext(ctx CommentContext) bool {

	// Every file is a special case not mapped in fileContext
	if ctx&InEveryFile != 0 {
		return true
	}

	for _, fCtx := range fileContext {
		if fCtx&ctx != 0 {
			return true
		}
	}
	return false
}
