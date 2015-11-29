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
	Context   CommentContext         // Used to select the correct comment
	NewCode   bool                   // When checking a diff, this indicates if the issue was found in existing or new code.
	Err       error                  // Any err encounted while building the issue.
	Metrics   map[string]interface{} // Any metrics that this issue was raised with
	Tags      []string               // Any tags this issue was raised with.

	commentSet *commentSet
	file       File // TODO(waigani) get File out of the issue struct.
	filename   string

	// TODO(matt, waigani) Implement this. Possibly use github.com/waigani/diffparser and github.com/waigani/astnode.
	// The idea is:
	// - issue.DiffFix() returns a diff patch to fix the issue.
	// - run Lingo with --fix. If issue.CanFix, Lingo prompts the user to keep/discard the patch.
	// - Lingo assembles patchs into one diff and, depending on flags, either applies the patch or just saves the diff to file.
	Patch string // A diff patch resolving the issue.
}

// TODO(waigani) These maps should not be globals. They need to hang off base.
var fileOrderToContext = map[int]CommentContext{
	1: FirstCommentInFile,
	2: SecondCommentInFile,
	3: ThirdCommentInFile,
	4: FourthCommentInFile,
	5: FifthCommentInFile,
}

var overallOrderToContext = map[int]CommentContext{
	1: FirstComment,
	2: SecondComment,
	3: ThirdComment,
	4: FourthComment,
	5: FifthComment,
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

func (issue *Issue) comments() *commentSet {
	if issue.commentSet == nil {
		issue.commentSet = &commentSet{}
	}
	return issue.commentSet
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
	file    map[string]map[string]int
	overall map[string]int
}

func newIssueOrder() *issueOrder {
	return &issueOrder{
		map[string]map[string]int{},
		map[string]int{},
	}
}

func (o *issueOrder) increment(issue *Issue) {
	if _, ok := o.file[issue.Name]; !ok {
		o.file[issue.Name] = map[string]int{}
	}
	o.file[issue.Name][issue.Filename()]++
	o.overall[issue.Name]++
}

// TODO(waigani) below is no longer used as we are reviewing files in sync.
// When the --keep-all flag is used, we should review all files async and use
// below.

//order by file name
// type byFile []*Issue

// func (issues byFile) Len() int {
// 	return len(issues)
// }
// func (issues byFile) Swap(i, j int) {
// 	issues[i], issues[j] = issues[j], issues[i]
// }
// func (issues byFile) Less(i, j int) bool {
// 	return issues[i].Filename() < issues[j].Filename()
// }

// // order by line
// type byLine []*Issue

// func (issues byLine) Len() int {
// 	return len(issues)
// }
// func (issues byLine) Swap(i, j int) {
// 	issues[i], issues[j] = issues[j], issues[i]
// }
// func (issues byLine) Less(i, j int) bool {
// 	li := issues[i].Position.Start.Line
// 	lj := issues[j].Position.Start.Line

// 	if li != lj {
// 		return li < lj
// 	}
// 	// if on same line, sort by column
// 	return issues[i].Position.Start.Column < issues[j].Position.Start.Line
// }
