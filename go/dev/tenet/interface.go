// Copyright 2015 Jesse Meek.
// Licensed under the AGPLv3, see LICENCE file for details.

package tenet

import (
	"go/ast"
	"go/token"

	"github.com/lingo-reviews/tenets/go/dev/api"
)

// Tenet is intended for use by tenet authors.
type Tenet interface {
	Info() *Info
	SetInfo(Info) Tenet

	RegisterMetric(key string) func(interface{}) RaiseIssueOption

	RegisterTag(tag string) RaiseIssueOption

	// Returns the name of the issue.
	RegisterIssue(issueName string, opts ...RegisterIssueOption) string

	// Returns the value of the option.
	RegisterOption(name, value, usage string) *string

	SmellNode(f smellNodeFunc) Tenet

	SmellLine(f smellLineFunc) Tenet
}

// Review is intended for use by tenet authors inside smell funcs.
type Review interface {

	// RaiseLineIssue sends the named issue (issueName) to lingo, along with
	// the start and end lines of the issue and metadata from opts.
	RaiseLineIssue(issueName string, start, end int, opts ...RaiseIssueOption) Review

	// RaiseNodeIssue sends the named issue (issueName) to lingo, along with metadata from n and opts.
	RaiseNodeIssue(issueName string, n ast.Node, opts ...RaiseIssueOption) Review

	// File is the current file being reviewed.
	File() File

	// The current smell will no longer be called.
	SmellDone()

	// The current smell will no longer be called for the current file.
	SmellDoneWithFile()

	// The current file will not be smelt by this tenet again. This should
	// only be used if it is not logical to keep looking. If you just want to
	// limit the number of times an issue is raised, use comment contexts.
	FileDone()
}

// BaseReview is only for use by the system.
type BaseTenet interface {
	Init()
	NewReview() *review
	Info() *Info
	MixinConfigOptions(opts []*api.Option) error

	SendError(error)
	Errors() chan error

	// This is a convinence method that gives us access to the base struct
	// when it is composed within a tenet.Tenet.
	base() *Base
}

// BaseReview is only for use by the system.
type BaseReview interface {

	// Starts the review. This should be called in a goroutine before sending
	// files to the review.
	StartReview()

	// Closes a review
	Close()

	// Send files to a review.
	SendFile(string)

	// Call this when you've finished sending all the files to review.
	EndReview()

	// File currently being reviewed.
	File() File

	// a chan of issues found in the review
	Issues() chan *Issue

	// have we found an issue for every context?
	areAllContextsMatched() bool

	// A temporary working directory that the review can use.
	TMPDIR() (dirpath string, err error)
}

// File represents a file being checked. It is intended for use by tenet authors.
type File interface {
	Line(i int) []byte
	Lines() [][]byte
	AST() *ast.File
	Fset() *token.FileSet
	Filename() string
}

// File represents a file being checked. It is intended for use by the system.
type BaseFile interface {
	newIssueRange(start, end int) *issueRange
	newIssueRangeFromNode(n ast.Node) *issueRange

	linePosition(line int) token.Position
	posLine(p token.Pos) []byte
	setLines([][]byte)
}
