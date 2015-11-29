// Copyright 2015 Jesse Meek.
// Licensed under the AGPLv3, see LICENCE file for details.

package tenet

import (
	"bytes"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"strings"

	"github.com/juju/errors"
)

// TODO(waigani) spit up interface. Only expose what we want the tenet author
// to see.

// TODO(waigani) PREALPHA clean this up. There are several methods that are not used and several that should not be exported.
// TODO(waigani) currently we need to pass in File to the ast.Vistor. Can we not? Can we have just the node?

// gofile implements File
type gofile struct {
	fset     *token.FileSet
	ast      *ast.File
	lines    [][]byte
	filename string
}

func (f *gofile) AST() *ast.File {
	return f.ast
}

func (f *gofile) Fset() *token.FileSet {
	return f.fset
}

func (f *gofile) Filename() string {
	return f.filename
}

func (f *gofile) linePosition(line int) token.Position {
	return token.Position{
		Filename: f.Filename(),
		Offset:   f.Fset().Base(), // offset, starting at 0
		Line:     line,            // line number, starting at 1
		// TODO(matt, waigani) implement. We'll have to update the AddIssue methods. Needs to be done before stable v1.
		// Column   int    // column number, starting at 1 (character count)
	}
}

func (f *gofile) newIssueRange(start, end int) *issueRange {
	return &issueRange{f.linePosition(start), f.linePosition(end)}
}

func (f *gofile) newIssueRangeFromNode(n ast.Node) *issueRange {
	s := f.Fset().Position(n.Pos())
	e := f.Fset().Position(n.End())
	return &issueRange{
		s,
		e,
	}
}

// PosLine returns the complete src line at p, including the terminating newline.
func (f *gofile) posLine(p token.Pos) []byte {
	return f.Line(f.Fset().Position(p).Line)
}

func (f *gofile) Line(i int) []byte {
	for n, b := range f.Lines() {
		if n+1 == i {
			return b
		}
	}
	return nil
}

func (f *gofile) Lines() [][]byte {
	return f.lines
}

func (f *gofile) setLines(lines [][]byte) {
	f.lines = lines
}

func (f *gofile) IsMain() bool {
	if f.AST().Name.Name == "main" {
		return true
	}
	return false
}

func (f *gofile) IsTest() bool { return strings.HasSuffix(f.Filename(), "_test.go") }

func buildFile(path, src string, fset *token.FileSet) (File, error) {
	var srcBytes []byte
	if src == "" {
		var err error
		// TODO(matt) TECHDEBT use "go/scanner".Scanner instead of loading all bytes to memory.
		srcBytes, err = ioutil.ReadFile(path)
		if err != nil {
			return nil, errors.Trace(err)
		}
	} else {
		srcBytes = []byte(src)
	}

	f, err := parser.ParseFile(fset, path, srcBytes, parser.ParseComments)
	if err != nil {
		return nil, errors.Trace(err)
	}

	// type info
	file := &gofile{
		filename: path,
		ast:      f,
		fset:     fset,
	}

	file.setLines(bytes.Split(srcBytes, []byte("\n")))
	return file, nil
}
