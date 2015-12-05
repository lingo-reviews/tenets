package testing

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"

	jt "github.com/juju/testing"
	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"

	"github.com/lingo-reviews/tenets/go/dev/api"
	"github.com/lingo-reviews/tenets/go/dev/tenet"
	"github.com/lingo-reviews/tenets/go/dev/tenet/log"
)

type TenetSuite struct {
	jt.CleanupSuite
	Tenet  tenet.Tenet
	Review tenet.Review
}

// TODO(waigani) we should also be testing issue.Position
type ExpectedIssue struct {
	Filename string
	Text     string
	Comment  string
	Metrics  map[string]interface{}
	Tags     []string
}

func (s *TenetSuite) baseTenet() tenet.BaseTenet {
	return s.Tenet.(tenet.BaseTenet)
}
func (s *TenetSuite) baseReview() tenet.BaseReview {
	return s.Review.(tenet.BaseReview)
}

func (s *TenetSuite) SetUpTest(c *gc.C) {
	// Tenet needs to be set before using this suite.
	c.Assert(s.Tenet, gc.NotNil)
	b := s.baseTenet()
	s.Review = b.NewReview()
	s.AddCleanup(func(c *gc.C) {
		s.Review.(tenet.BaseReview).Close()
	})
}

// SetCfgOption passes in an option that would normally be set on the CLI or
// in .lingo.
func (s *TenetSuite) SetCfgOption(c *gc.C, name, value string) {
	b := s.baseTenet()
	opts := []*api.Option{{
		Name:  name,
		Value: value,
	}}

	c.Assert(b.MixinConfigOptions(opts), jc.ErrorIsNil)
}

func (s *TenetSuite) CheckFiles(c *gc.C, files []string, expectedIssues ...ExpectedIssue) {
	log.Println("CheckFiles")

	br := s.baseReview()
	br.StartReview()
	s.sendFiles(files...)

	s.AssertExpectedIssues(c, ReadAllIssues(c, br), expectedIssues...)
}

func (s *TenetSuite) CheckSRC(c *gc.C, src string, expectedIssues ...ExpectedIssue) {
	fName := s.TmpFile(c, src)
	s.CheckFiles(c, []string{fName}, expectedIssues...)
}

func (s *TenetSuite) AssertExpectedIssues(c *gc.C, issues []*tenet.Issue, expectedIssues ...ExpectedIssue) {
	c.Assert(issues, gc.HasLen, len(expectedIssues))

	tmpDir, err := s.baseReview().TMPDIR()
	c.Assert(err, jc.ErrorIsNil)

	for i, issue := range issues {
		expected := expectedIssues[i]
		c.Assert(issue.Comment, gc.Equals, expected.Comment)
		c.Assert(strings.TrimPrefix(tmpDir, issue.Filename()), gc.Equals, strings.TrimPrefix(tmpDir, expected.Filename))
		c.Assert(issue.LineText, gc.Equals, expected.Text)
		c.Assert(issue.Metrics, jc.DeepEquals, expected.Metrics)
		c.Assert(issue.Tags, jc.SameContents, expected.Tags)
	}
}

func (s *TenetSuite) AssertExpectedSRCIssues(c *gc.C, issues []*tenet.Issue, expectedIssues ...ExpectedIssue) {
	c.Assert(issues, gc.HasLen, len(expectedIssues))

	for i, issue := range issues {
		expected := expectedIssues[i]
		c.Assert(issue.Comment, gc.Equals, expected.Comment)
		c.Assert(issue.LineText, gc.Equals, expected.Text)
	}
}

func (s *TenetSuite) AssertErrs(c *gc.C, errs []error) {
	for _, err := range errs {
		c.Assert(err, jc.ErrorIsNil)
	}
}

func (s *TenetSuite) CheckErrs(c *gc.C, errs []error) {
	for _, err := range errs {
		c.Check(err, jc.ErrorIsNil)
	}
}

var fileNameCounter int

func (s *TenetSuite) TmpFile(c *gc.C, src string) string {
	br := s.baseReview()
	fileNameCounter++
	// write a tmp file to check
	dir, err := br.TMPDIR()
	c.Assert(err, jc.ErrorIsNil)

	fName := filepath.Join(dir, fmt.Sprintf("test_src%d.go", fileNameCounter))
	err = ioutil.WriteFile(fName, []byte(src), 0644)
	c.Assert(err, jc.ErrorIsNil)
	return fName
}

// sendFiles should be called after br.StartReview and before ReadAllIssues
func (s *TenetSuite) sendFiles(files ...string) {
	go func() {
		br := s.baseReview()
		defer br.EndReview()
		for _, f := range files {

			br.SendFile(&api.File{Name: f})
		}
	}()
}

func ReadAllIssues(c *gc.C, r tenet.BaseReview) []*tenet.Issue {
	log.Println("reading all issues")
	var issues []*tenet.Issue
l:
	for {
		select {
		case issue, ok := <-r.Issues():
			if !ok {
				break l
			}
			issues = append(issues, issue)
		case <-time.After(10 * time.Second):
			c.Fatal("timed out waiting for issues")
			break l
		}
	}
	return issues
}
