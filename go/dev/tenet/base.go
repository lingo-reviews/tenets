// Copyright 2015 Jesse Meek.
// Licensed under the AGPLv3, see LICENCE file for details.

package tenet

import (
	"fmt"
	"go/token"
	"strings"

	"github.com/juju/errors"
	"github.com/lingo-reviews/tenets/go/dev/api"
	"github.com/lingo-reviews/tenets/go/dev/tenet/log"
)

// Base implements the Tenet interface and is intended to be composed
// with 3rd party tenets.
type Base struct {
	// all issues this tenet looks for.
	registeredIssues map[string]*Issue

	// all issues this tenet found.
	issuesc chan *Issue

	// errors holds non-fatal errors. // TODO(waigani) currently the server
	// just logs these. We need to make logs visable to user and provide
	// options for the user around errors.
	errorsc chan error

	// tmpdir is the dir for tenets to work in.
	tmpdir string

	astVisitors  []astVisitor
	lineVisitors []lineVisitor

	info *Info
}

// base allows us to access the base struct when it's embeded in another
// struct with a Tenet interface type.
func base(t Tenet) *Base {
	return t.(BaseTenet).base()
}

func (b *Base) base() *Base {
	return b
}

func (b *Base) Init() {
	b.errorsc = make(chan error, 1)
}

func (b *Base) NewReview() *review {
	r := &review{
		tenet:       b,
		issuesc:     make(chan *Issue),
		filesc:      make(chan *api.File),
		waitc:       make(chan struct{}),
		fileDoneMap: map[string]bool{},
	}
	go func() {
		<-r.waitc
		if r.issuesc != nil {
			log.Println("closing issuesc")
			close(r.issuesc)
		}
		r.waitc = nil
	}()

	return r
}

func (b *Base) SendError(err error) {
	log.Println("sending error", err)
	log.Println(b.errorsc == nil)
	select {
	case b.errorsc <- err:
	default:
	}
}

func (b *Base) Errors() chan error {
	return b.errorsc
}

func (b *Base) SmellNode(f smellNodeFunc) Tenet {

	b.astVisitors = append(b.astVisitors, astVisitor{
		smellNode: f,
		fileDone:  map[string]bool{},
	})
	return b
}

func (b *Base) SmellLine(f smellLineFunc) Tenet {
	b.lineVisitors = append(b.lineVisitors, lineVisitor{
		visit:    f,
		fileDone: map[string]bool{},
	})
	return b
}

func (b *Base) MixinConfigOptions(opts []*api.Option) error {
	for _, opt := range opts {
		if err := b.setOpt(opt); err != nil {
			return errors.Trace(err)
		}
	}
	return nil
}

func (b *Base) setOpt(opt *api.Option) error {
	if b.info == nil {
		return errors.New("tenet info is nil")
	}
	for _, bOpt := range b.info.Options {
		if bOpt.name == opt.Name {
			*bOpt.value = opt.Value
			return nil
		}
	}
	return errors.Errorf("tenet has no option %q", opt.Name)
}

func AddComment(comment string, ctx ...CommentContext) RegisterIssueOption {
	return func(issue *Issue) {

		if issue.commentSet == nil {
			issue.commentSet = &commentSet{}
		}
		issue.commentSet.AddComment(comment, ctx...)
	}
}

type option struct {
	name  string
	value *string
	usage string
}

// TODO(waigani) support interface values
func (b *Base) RegisterOption(name string, value string, usage string) *string {

	// toml doesn't support "-"
	blacklistChars := "- "
	if strings.ContainsAny(name, blacklistChars) {
		// Yes panic, this is a developer error.
		msg := fmt.Sprintf("option name %q cannot contain any of the following characters: %q", name, blacklistChars)
		panic(msg)
	}

	if b.info == nil {
		// Yes panic, this is a developer error.
		panic("options cannot be registered before tenet.Info is set, please do that first.")
	}
	v := &value
	b.info.Options = append(b.info.Options, &option{
		name:  name,
		value: v,
		usage: usage,
	})

	// return the value from the pointer for this option which will either be the
	// default value passed in or one updated by the user.
	return v
}

// RegisterMetric registers a metric key name that can be used when raising an issue.
func (b *Base) RegisterMetric(key string) func(val interface{}) RaiseIssueOption {
	b.info.metrics = append(b.info.metrics, key)

	return func(val interface{}) RaiseIssueOption {
		return func(issue *Issue) {

			if issue.Metrics == nil {
				issue.Metrics = map[string]interface{}{}
			}

			issue.Metrics[key] = val
		}
	}
}

// RegisterTag registers a tag name that can be used when registering an issue.
func (b *Base) RegisterTag(tag string) RaiseIssueOption {
	b.info.tags = append(b.info.tags, tag)

	return func(issue *Issue) {
		issue.Tags = append(issue.Tags, tag)
	}
}

func (b *Base) RegisterIssue(issueName string, opts ...RegisterIssueOption) string {
	issue := &Issue{
		Name:       issueName,
		commentSet: &commentSet{},
		CommVars:   map[string]interface{}{},
	}

	for _, opt := range opts {
		opt(issue)
	}
	if b.registeredIssues == nil {
		b.registeredIssues = map[string]*Issue{}
	}
	b.registeredIssues[issueName] = issue

	return issueName
}

type errWithContext struct {
	err     error
	errLine *token.Position
}

func (e *errWithContext) Error() string {
	return e.err.Error()
}

func (e *errWithContext) Line() string {
	return e.errLine.String()
}

// TODO(waigani) call this handleTenetError and make a TenetError type - only
// those can be passed in.
// posOfErr is the position of the node/line that was being parsed when the
// error occoured.
func (b *Base) addErrOnErr(err error, f File, posOfErr interface{}) bool {
	if err != nil {
		// TODO(waigani) this log is a quick hack. We should read all the errs off errorsc.
		log.Println(err.Error())
		errCtx := &errWithContext{err: err}
		switch p := posOfErr.(type) {
		case token.Pos:
			fset := f.Fset()
			pos := fset.Position(p)
			errCtx.errLine = &pos
		case int:
			line := f.(BaseFile).linePosition(p)
			errCtx.errLine = &line
		default:
			panic(fmt.Sprintf("unknown posOfErr type: %T", posOfErr))
		}

		b.errorsc <- errCtx
		return true
	}
	return false
}
