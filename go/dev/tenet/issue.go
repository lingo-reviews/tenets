// Copyright 2015 Jesse Meek.
// Licensed under the AGPLv3, see LICENCE file for details.

package tenet

import (
	"go/token"
	"strings"
)

// TODO(waigani) do we need this?
// "honnef.co/go/importer"

// Problem represents a problem in some source code.
// Borrows from problem struct from https://github.com/golang/lint/blob/master/lint.go
type Issue struct {
	Name      string                 // Name is the the name of the checker that added the issue
	Position  *issueRange            // position in source file
	Comment   string                 // The rendered comment for this issue.
	CommVars  map[string]interface{} // key/value pairs for use in comment template variables e.g. {{.somevarname}}
	CtxBefore string                 // source lines before the problem line(s)
	LineText  string                 // the source line(s)
	CtxAfter  string                 // source lines after the problem line(s)
	Link      string                 // (optional) the link to the style guide for the problem
	NewCode   bool                   // When checking a diff, this indicates if the issue was found in existing or new code.
	Err       error                  // Any err encounted while building the issue.
	Metrics   map[string]interface{} // Any metrics that this issue was raised with
	Tags      []string               // Any tags this issue was raised with.
	comments  []*comment             // A slice of possible comments for this issue.
	file      File                   // TODO(waigani) get File out of the issue struct.
	filename  string

	// TODO(matt, waigani) Implement this. Possibly use github.com/waigani/diffparser and github.com/waigani/astnode.
	// The idea is:
	// - issue.DiffFix() returns a diff patch to fix the issue.
	// - run Lingo with --fix. If issue.CanFix, Lingo prompts the user to keep/discard the patch.
	// - Lingo assembles patchs into one diff and, depending on flags, either applies the patch or just saves the diff to file.
	Patch string // A diff patch resolving the issue.
}

func (issue *Issue) Filename() string { // TODO(waigani) Remove this and File from issue and just use issue.filename
	return issue.Position.Start.Filename
}

func (issue *Issue) setSource(iRange *issueRange) *Issue {
	issue.Position = iRange
	issue.CtxBefore, issue.CtxAfter = getLineContext(issue.file, iRange.Start.Line, iRange.End.Line)

	start := iRange.Start.Line
	issue.LineText = string(issue.file.Line(start))
	for i := start + 1; i < iRange.End.Line; i++ {
		issue.LineText += "\n" + string(issue.file.Line(i))
	}
	return issue
}

func (issue *Issue) addComment(commentTemplate string, contexts ...CommentContext) {
	com := &comment{
		Template: commentTemplate,
		matches:  map[CommentContext]bool{},
	}

	// Split contexts into file and comment contexts.
	for _, ctx := range contexts {
		if isFileContext(ctx) {
			com.addFileCtx(ctx)
		} else {
			com.addCommentCtx(ctx)
		}
	}

	// Set defaults
	if len(com.commentContexts) == 0 {
		com.addCommentCtx(DefaultComment)
		if len(com.fileContexts) == 0 {
			com.addFileCtx(InEveryFile)
		}
	} else if len(com.fileContexts) == 0 {
		com.addFileCtx(InOverall)
	}

	issue.comments = append(issue.comments, com)
}

// TODO(waigani) rename getStringContext
func getLineContext(f File, issueStartLine, issueEndLine int) (ctxBefore, ctxAfter string) {
	buffer := 4 // TODO(waigani) make this a config.
	beforeStart := issueStartLine - buffer
	if beforeStart < 0 {
		beforeStart = 0
	}
	afterEnd := issueEndLine + buffer
	end := len(f.Lines())
	if afterEnd > end {
		afterEnd = end
	}

	var before, after []string
	for i := beforeStart; i < issueStartLine; i++ {
		before = append(before, string(f.Line(i)))
	}

	for i := issueEndLine + 1; i <= afterEnd; i++ {
		after = append(after, string(f.Line(i)))
	}

	return strings.Join(before, "\n"), strings.Join(after, "\n")
}

// TODO(waigani) NextComment - moves through each comment in context

type issueRange struct {
	Start token.Position
	End   token.Position
}

type RegisterIssueOption func(*Issue)
type RaiseIssueOption func(*Issue)

// CommentVar stores key/value pairs used to populate the rule's comment
// templates e.g.
// r.AddInfo("spaces", 2)
// r.Comments = []string{"You have {{.spaces}} after a period, when the styleguide specifies that there should only be one."
// CommentVar returns an issueOption which sets a variable to be used in the comments.
func CommentVar(key string, value interface{}) RaiseIssueOption {
	return func(issue *Issue) {
		issue.CommVars[key] = value
	}
}

// --- Issue ordering ---

// IssueOrder is used to keep track of when each issue is reported. It is used
// to set the comment context for each issue.
// file: issueName:fileName:issueCount
// overall: issueName:issueCount
type issueOrder struct {

	// map of issue name to number of times issue has been found.
	issueCount map[string]int //map[issueName]issueCount

	// map of issue name to number of times an issue has been found in a file.
	issueInFileCount map[string]map[string]int //map[issueName]map[fileCount]issueCount

	// map of issue name to when a file with this issue was first found.
	fileOrder map[string]map[string]int
}

func newIssueOrder() *issueOrder {
	return &issueOrder{
		issueCount:       map[string]int{},
		issueInFileCount: map[string]map[string]int{},
		fileOrder:        map[string]map[string]int{},
	}
}

func (o *issueOrder) increment(issueName, fileName string) {
	// Setup map.
	if o.issueInFileCount[issueName] == nil {
		o.issueInFileCount[issueName] = map[string]int{}
	}
	if o.fileOrder[issueName] == nil {
		o.fileOrder[issueName] = map[string]int{}
	}

	// Keep track of the first n files this issue was found in, where n == number of file contexts.
	if o.fileOrder[issueName][fileName] == 0 {
		o.fileOrder[issueName][fileName] = len(o.fileOrder[issueName]) + 1
	}

	// Keep track of the first n times this issue was found, where n == number of comment contexts.
	// if o.issueCount[issueName] < len(commContext) {
	o.issueCount[issueName]++
	// }

	// Keep track of the first n times this issue was found in this file, where n == number of comment contexts.
	// if o.issueInFileCount[issueName][fileName] < len(commContext) {
	o.issueInFileCount[issueName][fileName]++
	// }
}
