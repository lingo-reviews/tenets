package tenet

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/token"
	"io/ioutil"
	"os"
	"reflect"
	"time"

	"github.com/juju/errors"

	"text/template"

	"github.com/lingo-reviews/tenets/go/dev/api"
	"github.com/lingo-reviews/tenets/go/dev/tenet/log"
)

type review struct {
	tenet   Tenet
	waitc   chan struct{}
	filesc  chan *api.File
	issuesc chan *Issue

	// file is the file currently under review.
	file       File
	issueOrder *issueOrder

	// a scratch dir for artefacts while reviewing.
	tmpdir string

	fileDone func()

	// SmellDone is set before each line/node is visited. If called, that
	// visitor no longer visits that file.
	smellDone func()

	smellDoneWithFile func()

	fileDoneMap map[string]bool
}

// StartReview listens for files sent to r.SendFile(filename) and reviews them.
func (r *review) StartReview() {
	go func() {
		defer r.Close()
		log.Println("started review")
		b := base(r.tenet)

		// check files synchronously to ensure correct ordering and that we stop
		// after the context is full.
		fset := token.NewFileSet()
		log.Println("reading off filesc")

		for {
			select {
			case file, ok := <-r.filesc:
				if !ok && file == nil {
					log.Println("all files reviewed.")
					return
				}

				f, err := buildFile(file.Name, "", fset, file.Lines)
				if err != nil {
					log.Println("could not build file")
					b.SendError(errors.Annotatef(err, "could not find file: %q", file))
					continue
				}
				log.Println("checking file", file)
				err = r.check(f)
				b.addErrOnErr(err, f, 0)
				log.Println("finished checking file", file)
			case <-time.After(3 * time.Second):
				b.errorsc <- errors.New("timed out waiting for file")
				return
			}
		}
	}()
}

func (r *review) SendFile(file *api.File) {
	r.filesc <- file
}

func (r *review) EndReview() {
	close(r.filesc)
}

// Closes a review. Can be called more than once.
func (r *review) Close() {
	if r.waitc != nil {
		log.Println("closing review and waitc")
		close(r.waitc)
		os.RemoveAll(r.tmpdir)
		r.waitc = nil
	}
}

func (r *review) IsClosed() bool {
	return r.waitc == nil
}

func (r *review) sendIssue(issue *Issue) {
	r.issuesc <- issue
}

func (r *review) Issues() chan *Issue {
	return r.issuesc
}

// check cannot be called async as we set r.file on the review pointer. This
// has to remain until we've finished walking.
func (r *review) check(f File) error {
	log.Printf("checking file %s", f.Filename())
	// set current file being reviewed.
	r.file = f
	b := r.baseTenet()

	r.fileDone = func() {
		r.fileDoneMap[f.Filename()] = true
	}

	// TODO(waigani) this should be r.recursiveASTWalk()
	for _, visitor := range b.astVisitors {
		r.walkAST(&visitor)
	}

	// first walk all ast nodes.
	// r.recursiveASTWalk(b, f)

	// TODO(waigani) support recursive line visits.
	if len(b.lineVisitors) > 0 {
		// then check all src lines
		r.visitLines()
	}

	return nil
}

// TODO(waigani) Get this working. We need to group astVisitors by filename.
// Top level smells will be added to b.astVisitor.
// Nested smells will be added with r.SmellNode
// They will be added to the file collection of visitors: r.astVisitors[filename]
// That collection will be used in the recursiveASTWallk func
//
// Slices are only read once at the beginning of a loop, so we need to use a
// recursive func to update the collection. At the end of the walks, the
// original range is removed. If new visitors were added during the walks,
// they are then run.
// func (r *review) recursiveASTWalk() {
// 	xxx.Print(r.File().Filename())
// 	l := len(b.astVisitors)
// 	for _, visitor := range visitors {
// 		r.walkAST(&visitor)
// 	}

// 	b.astVisitors.deleteRange(0, l-1)
// 	if len(b.astVisitors) > 0 {
// 		r.recursiveASTWalk(b, f)
// 	}
// }

func (r *review) baseTenet() *Base {
	return r.tenet.(BaseTenet).base()
}

func (r *review) TMPDIR() (_ string, err error) {
	if r.tmpdir == "" {
		r.tmpdir, err = ioutil.TempDir(os.TempDir(), "tenet_review_"+RandString(5))
	}
	return r.tmpdir, err
}

// --- visitor methods ---

func (r *review) SmellDoneWithFile() {
	r.smellDoneWithFile()
}

func (r *review) SmellDone() {
	r.smellDone()
}

func (r *review) FileDone() {
	r.fileDone()
}

type lineVisitor struct {
	fileDone map[string]bool
	done     bool
	visit    smellLineFunc
}

type smellLineFunc func(r Review, n int, line []byte) error

func (l *lineVisitor) Visit(r Review, n int, line []byte) error {
	r.(*review).smellDoneWithFile = func() {
		l.fileDone[r.File().Filename()] = true
	}
	r.(*review).smellDone = func() {
		l.done = true
	}
	return l.visit(r, n, line)
}

func (l *lineVisitor) isSmellDoneWithFile(filename string) bool {
	return l.fileDone[filename]
}

func (l *lineVisitor) isSmellDone() bool {
	return l.done
}

func lineInDiff(diff []int64, lineNo int64) bool {
	for _, l := range diff {
		if lineNo == l {
			return true
		}
	}
	return false
}

func (r *review) visitLines() {
	b := r.baseTenet()
	f := r.File()
	fName := f.Filename()

	for _, v := range b.lineVisitors {
		for i, line := range f.Lines() {

			diff := f.(BaseFile).diff()
			if len(diff) > 0 && !lineInDiff(diff, int64(i+1)) {
				continue
			}

			if v.isSmellDoneWithFile(fName) || v.isSmellDone() || r.isFileDone(fName) {
				break
			}
			n := i + 1
			b.addErrOnErr(v.Visit(r, n, line), f, n)
		}
	}
}

func (r *review) isFileDone(filename string) bool {
	return r.fileDoneMap[filename]
}

// astVisitor walks through each AST node.
type astVisitor struct {
	done      bool
	fileDone  map[string]bool
	visit     func(node ast.Node) (w ast.Visitor)
	smellNode smellNodeFunc
}

// TODO(waigani) should be of type: func(Review, ast.Node) error
type smellNodeFunc interface{}

func (v *astVisitor) Visit(node ast.Node) (w ast.Visitor) {
	return v.visit(node)
}

func (v *astVisitor) isSmellDone() bool {
	return v.done
}

func (v *astVisitor) isSmellDoneWithFile(filename string) bool {
	return v.fileDone[filename]
}

func (r *review) walkAST(v *astVisitor) {
	v.visit = func(node ast.Node) (w ast.Visitor) {
		file := r.File()
		// TODO(waigani) quick hack to get diff working. Come back and work out what's going on with diff?
		if node == nil || v.isSmellDone() || !nodeInDiff(file.(BaseFile), node) {
			// Keep walking other nodes.
			return v
		}

		fName := file.Filename()
		if v.isSmellDoneWithFile(fName) || r.isFileDone(fName) {
			// Stop walking all nodes.
			return nil
		}

		funcType := reflect.TypeOf(v.smellNode)
		nodeType := funcType.In(1)
		if nodeType == nil {
			panic("AST Visitor function signature does not have the right format")
		}
		expectedType := fmt.Sprintf("%T", node)
		obtainedType := nodeType.String()

		// If the type of node ast.Walk is visiting matches the type of node
		// in the astVisitor func, call the func with the node.
		if obtainedType == expectedType {

			// set review funcs
			r.smellDoneWithFile = func() {
				v.fileDone[r.File().Filename()] = true
			}

			r.smellDone = func() {
				v.done = true
			}

			// why doesn't type conversion work? Does it work in later versions of Go?
			// wish this worked: visitorFunc.(func(Review, ast.Node))(r, node)
			f := reflect.ValueOf(v.smellNode)
			rV := reflect.ValueOf(r)
			nodeV := reflect.ValueOf(node)
			refV := f.Call([]reflect.Value{rV, nodeV})
			if len(refV) > 0 {
				if err := refV[0].Interface(); err != nil {
					b := r.baseTenet()
					b.SendError(err.(error))
				}
			}
			return
		}
		return v
	}
	ast.Walk(v, r.File().AST())
}

func nodeInDiff(f BaseFile, node ast.Node) bool {
	start := f.(*gofile).fset.Position(node.Pos())
	end := f.(*gofile).fset.Position(node.End())

	diff := f.diff()
	if len(diff) == 0 {
		// skip diff check
		return true
	}

	for i := start.Line; i <= end.Line; i++ {
		if lineInDiff(diff, int64(i)) {
			return true
		}
	}

	return false
}

func (r *review) File() File {
	return r.file
}

func (r *review) RaiseNodeIssue(issueName string, n ast.Node, opts ...RaiseIssueOption) Review {
	return r.raiseIssue(issueName, r.File().(BaseFile).newIssueRangeFromNode(n), opts)
}

func (r *review) RaiseLineIssue(issueName string, start, end int, opts ...RaiseIssueOption) Review {
	return r.raiseIssue(issueName, r.File().(BaseFile).newIssueRange(start, end), opts)
}

func (r *review) raiseIssue(issueName string, iRange *issueRange, opts []RaiseIssueOption) Review {
	if r.areAllContextsMatched() {
		return r
	}

	b := r.baseTenet()
	// TODO(waigani) error handle this.
	i := b.registeredIssues[issueName]
	if i == nil {
		// Yes panic, this is a developer error.
		msg := fmt.Sprintf("issue %q cannot be raised before it is registered", issueName)
		panic(msg)
	}
	issue := &Issue{}
	x := *i
	x.copyTo(issue)

	// TODO(waigani) This a Go blemish. Is there a nicer way to copy a struct with a map?
	issue.CommVars = map[string]interface{}{}
	for _, opt := range opts {
		opt(issue)
	}

	// TODO(waigani) this is a quick hack. We need to pull File out of *Issue.
	issue.file = r.File()
	issue.setSource(iRange)

	if err := r.setContextualComment(issue); err != nil {
		// If no comment has been set for the context in which this issue was
		// found, don't raise it.
		if err == errNoCommentForContext {
			return r
		}
		issue.Err = err
	}

	log.Println("sending issue")
	r.sendIssue(issue)
	log.Println("not blocked")

	if r.areAllContextsMatched() {

		// This is our last issue raised, close the issue chan. Note: if a
		// tenet reviews async, this will become a race condition.
		r.Close()
	}
	return r
}

func (self *Issue) copyTo(newIssue *Issue) {
	*newIssue = *self
}

// --- issue context ---

func addContextIfMissing(collection CommentContext, contexts ...CommentContext) CommentContext {
	for _, context := range contexts {
		if collection&context == 0 {
			collection |= context
		}
	}
	return collection
}

// isContextFull returns true if all contexts for all comments for all issues have been used.
func (r *review) areAllContextsMatched() bool {
	b := r.baseTenet()

	// TODO(waigani) keep a running tally of matched contexts and only iterate
	// over those that have not been found.
	for _, issue := range b.registeredIssues {
		for _, comm := range issue.comments {
			if !comm.allContextsMatched() {
				return false
			}
		}
	}
	return true
}

var errNoCommentForContext = errors.New("there are no comments for this context")

func (r *review) getIssueOrder() *issueOrder {
	if r.issueOrder == nil {
		r.issueOrder = newIssueOrder()
	}
	return r.issueOrder
}

// setContextualComment applies the contextual comment to this issue and
// updates its internal state of context to apply to the next issue. It
// assumes issues will come in synchronously.
func (r *review) setContextualComment(issue *Issue) error {
	o := r.getIssueOrder()
	issueName := issue.Name
	filename := r.File().Filename()
	o.increment(issueName, filename)

	fileCtx := r.currentFileContext(issue.Name)
	commentInFileCtx := commContext[o.issueInFileCount[issueName][filename]]
	commentInOverallCtx := commContext[o.issueCount[issueName]]

	var foundComments []*comment
	for _, comm := range issue.comments {

		var found bool
		// does the comment match an in-file context?
		foundFileCtx := fileCtx | commentInFileCtx | InEveryFile | DefaultComment
		if comm.matchesContext(foundFileCtx) {
			found = true
			comm.addMatch(fileCtx)
		}
		// does the comment match an overall context?
		foundOverallCtx := InOverall | commentInOverallCtx | DefaultComment
		if comm.matchesContext(foundOverallCtx) {
			found = true
			comm.addMatch(foundOverallCtx)
		}

		if found {
			foundComments = append(foundComments, comm)
		}
	}

	if len(foundComments) == 0 {
		return errNoCommentForContext
	}

	var err error
	issue.Comment, err = buildComment(foundComments[0].Template, issue.CommVars)
	return err
}

func buildComment(commentTemplate string, commentVars map[string]interface{}) (string, error) {
	// Build comments with template args
	t := template.New("comment template")

	ct, err := t.Parse(commentTemplate)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err = ct.Execute(&buf, commentVars); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// returns the current file's context for this issue.
func (r *review) currentFileContext(issueName string) CommentContext {
	o := r.getIssueOrder()
	fOrder := o.fileOrder[issueName][r.File().Filename()]
	return fileContext[fOrder]
}
