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
	// set current file being reviewed.
	r.file = f
	b := r.baseTenet()

	r.fileDone = func() {
		r.fileDoneMap[f.Filename()] = true
	}
	// first walk all ast nodes.
	if len(b.astVisitors) > 0 {
		v := r.getVisitor()
		ast.Walk(v, f.AST())
	}

	if len(b.lineVisitors) > 0 {
		// then check all src lines
		r.visitLines()
	}

	return nil
}

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

func (r *review) getVisitor() *astVisitor {
	b := r.baseTenet()
	v := &astVisitor{fileDone: map[string]bool{}}
	v.visit = func(node ast.Node) (w ast.Visitor) {
		visitors := b.astVisitors
		file := r.File()
		fName := file.Filename()

		// Stop walking the AST tree if we don't have any visitors or this
		// tenet is done reviewing this file.
		if len(visitors) == 0 || r.isFileDone(fName) {
			return nil
		}
		if node == nil {
			return w
		}

		visitNode := func(visitor astVisitor) {
			funcType := reflect.TypeOf(visitor.smellNode)
			nodeType := funcType.In(1)
			if nodeType == nil {
				panic("AST Visitor function signature does not have the right format")
			}
			expectedType := fmt.Sprintf("%T", node)
			obtainedType := nodeType.String()

			// If the type of node ast.Walk is visiting matches the type of node
			// in the astVisitor func, call the func with the node.
			if obtainedType == expectedType {

				if !nodeInDiff(file.(BaseFile), node) {
					return
				}

				// set review funcs
				r.smellDoneWithFile = func() {
					visitor.fileDone[r.File().Filename()] = true
				}

				r.smellDone = func() {
					visitor.done = true
				}

				// why doesn't type conversion work? Does it work in later versions of Go?
				// wish this worked: visitorFunc.(func(Review, ast.Node))(r, node)
				f := reflect.ValueOf(visitor.smellNode)
				rV := reflect.ValueOf(r)
				nodeV := reflect.ValueOf(node)
				r := f.Call([]reflect.Value{rV, nodeV})
				if len(r) > 0 {
					if err := r[0].Interface(); err != nil {
						b.SendError(err.(error))
					}
				}
				return
			}
			return
		}
		for _, visitor := range visitors {
			if visitor.isSmellDone() || visitor.isSmellDoneWithFile(fName) {
				continue
			}
			visitNode(visitor)
		}
		return v
	}
	return v
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
		if lineInDiff(diff, int64(start.Line)) {
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
	r.setContext(issue)

	if err := r.setContextualComment(issue); err != nil {
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

// SetContext applies the context to this issue and updates its internal state of
// context to apply to the next issue. It assumes issues will come in
// synchronously.
func (r *review) setContext(issue *Issue) {
	if r.issueOrder == nil {
		r.issueOrder = newIssueOrder()
	}
	o := r.issueOrder
	o.increment(issue)

	overallCtx := overallOrderToContext[o.overall[issue.Name]]
	fileCtx := fileOrderToContext[o.file[issue.Name][issue.Filename()]]

	issue.Context = overallCtx | fileCtx
}

// isContextFull returns true if all comments for all issues have been used.
func (r *review) areAllContextsMatched() bool {
	b := r.baseTenet()
	for _, issue := range b.registeredIssues {
		commSet := issue.comments()
		if len(commSet.commentsForContext(DefaultComment)) > 0 {
			return false
		}

		for _, comm := range commSet.Comments {
			// If an issue has a comment for a context that has not yet been
			// matched, return false.
			if comm.Matched == false {
				return false
			}
		}
	}
	return true
}

func (r *review) setContextualComment(issue *Issue) error {
	commSet := issue.comments()
	comments := commSet.commentsForContext(issue.Context)
	for _, comm := range comments {
		comm.Matched = true
	}

	// TODO(waigani) allow user to limit number of default comments.
	if len(comments) == 0 {
		comments = commSet.commentsForContext(DefaultComment)
	}

	// build comments with template args
	t := template.New("comment template")
	// default message if no comment set
	commentTemplate := "Issue Found"
	if len(comments) > 0 {
		commentTemplate = comments[0].Template
	}
	ct, err := t.Parse(commentTemplate) // TODO(waigani) This only returns the first comment for each context.
	if err != nil {
		return err
	}
	var buf bytes.Buffer
	if err = ct.Execute(&buf, issue.CommVars); err != nil {
		return err
	}
	issue.Comment = buf.String()

	// set the comment as used
	// comm.Used

	return nil
}
