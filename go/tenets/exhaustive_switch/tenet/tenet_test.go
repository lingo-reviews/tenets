package tenet_test

import (
	"testing"

	tt "github.com/lingo-reviews/tenets/go/dev/tenet/testing"
	gc "gopkg.in/check.v1"

	exhaustiveSwitch "github.com/lingo-reviews/tenets/go/tenets/exhaustive_switch/tenet"
)

func Test(t *testing.T) {
	gc.TestingT(t)
}

type exhaustiveSwitchSuite struct {
	tt.TenetSuite
}

var _ = gc.Suite(&exhaustiveSwitchSuite{})

func (s *exhaustiveSwitchSuite) SetUpSuite(c *gc.C) {
	l := exhaustiveSwitch.New()
	s.Tenet = l
}

func (s *exhaustiveSwitchSuite) TestExampleFiles(c *gc.C) {
	files := []string{
		"example/example.go",
	}

	expectedIssues := []tt.ExpectedIssue{
		{
			Filename: "example/example.go",
			Text: `
	switch p {
	case onefish:
	default:`[1:],
			Comment: "The following cases are missing from this switch: twofish, redfish, bluefish.",
		}, {
			Filename: "example/example.go",
			Text: `
	switch p {
	case onefish, twofish:
	default:`[1:],
			Comment: "The following cases are missing from this switch: redfish, bluefish.",
		}, {
			Filename: "example/example.go",
			Text: `
	switch p {
	case redfish:
	case bluefish:
	default:`[1:],
			Comment: "The following cases are missing from this switch: onefish, twofish.",
		}, {
			Filename: "example/example.go",
			Text: `
	switch p {
	case onefish, twofish:
	case redfish:
	default:`[1:],
			Comment: "The following cases are missing from this switch: bluefish.",
		}, {
			Filename: "example/example.go",
			Text: `
	switch p {
	case onefish:
	default:`[1:],
			Comment: "The following cases are missing from this switch: twofish, redfish, bluefish.",
		},
	}

	s.CheckFiles(c, files, expectedIssues...)
}
