// Copyright 2015 Jesse Meek.
// Licensed under the AGPLv3, see LICENCE file for details.

package tenet

import (
	"go/ast"
	"go/token"

	"github.com/lingo-reviews/tenets/go/dev/api"
)

// --- tenet author interfaces ---

// As a tenet author, Tenet, Review and File are the tools of your trade.

// Tenet defines what the tenet is about and sets anything needed before a
// review. Tenet should never be called inside SmellNode and SmellLine.
type Tenet interface {

	// Returns information on the tenet.
	Info() *Info

	// Sets information about the tenet used by the lingo client.
	SetInfo(Info) Tenet

	// Register a metric which can be added to a raised issue.
	RegisterMetric(key string) func(interface{}) RaiseIssueOption

	// Register a tag which can be added to a raised issue
	RegisterTag(tag string) RaiseIssueOption

	// Registers an issue that this tenet can raise. The string returned is
	// the the name of the issue which is used as the first argument to
	// Review.RaiseNodeIssue or Review.RaiseLineIssue
	RegisterIssue(issueName string, opts ...RegisterIssueOption) string

	// Returns the value of the option.
	RegisterOption(name, value, usage string) *string

	// SmellNode will smell every node that matches the type in smellNodeFunc.
	SmellNode(f smellNodeFunc) Tenet

	// SmellLine will smell every line of every file.
	SmellLine(f smellLineFunc) Tenet
}

// Review should only be used inside SmellNode and SmellLine.
type Review interface {

	// RaiseLineIssue sends the named issue (issueName) to lingo, along with
	// the start and end lines of the issue and metadata from opts.
	RaiseLineIssue(issueName string, start, end int, opts ...RaiseIssueOption) Review

	// RaiseNodeIssue sends the named issue (issueName) to lingo, along with metadata from n and opts.
	RaiseNodeIssue(issueName string, n ast.Node, opts ...RaiseIssueOption) Review

	// File is the current file being reviewed.
	File() File

	// The current smell will no longer be called at all.
	SmellDone()

	// The current smell will no longer be called for the current file.
	SmellDoneWithFile()

	// The current file will not be smelt by this tenet again. This should
	// only be used if it is not logical to keep looking. If you just want to
	// limit the number of times an issue is raised, use comment contexts. e.g.
	// tenet.FirstCommentInFile
	FileDone()
}

// File represents the current file being reviewed.
type File interface {

	// Returns the source code a line i.
	Line(i int) []byte

	// Returns all source code for this file.
	Lines() [][]byte

	// Returns the ast.File node of this file.
	AST() *ast.File

	// The name of the file currently being reviewed.
	Filename() string

	// Returns the FileSet this file is a member of.
	Fset() *token.FileSet
}

// --- system interfaces ----

// As a tenet author, you can safely ignore these.

// BaseReview is only for use by the system.
type BaseReview interface {

	// These methods are only exported because other system packages need
	// them.

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

	// A chan of issues found in the review.
	Issues() chan *Issue

	// have we found an issue for every context?
	areAllContextsMatched() bool

	// TODO(waigani) move this Tenet interface as it is useful for tenet authors.
	// A temporary working directory that the review can use.
	TMPDIR() (dirpath string, err error)
}

// BaseReview is only for use by the system.
type BaseTenet interface {

	// These methods are only exported because other system packages need
	// them.

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

// File represents a file being checked. It is intended for use by the system.
type BaseFile interface {
	newIssueRange(start, end int) *issueRange
	newIssueRangeFromNode(n ast.Node) *issueRange
	linePosition(line int) token.Position
	posLine(p token.Pos) []byte
	setLines([][]byte)
	diff() []int64
}
