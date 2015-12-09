package network_test

import (
	"strings"

	gc "gopkg.in/check.v1"
)

// no issue
func (s *suite) TestTable(c *gc.C) {

	a := []string{
		"a",
		"a",
	}

	for _, s := range a {
		c.Assert(s, gc.Equals, "a")
	}
}

// no issue
func (s *suite) TestTable2(c *gc.C) {

	var a = []string{
		"a",
		"a",
	}

	for _, s := range a {
		c.Assert(s, gc.Equals, "a")
	}
}

// no issue
func (s *suite) TestTable3(c *gc.C) {

	var a []string = []string{
		"a",
		"a",
	}

	for _, s := range a {
		c.Assert(s, gc.Equals, "a")
	}
}

// issue
func (s *suite) TestTable4(c *gc.C) {
	a := list()

	for _, s := range a {
		c.Assert(s, gc.Equals, "a")
	}
}

// issue
func (s *suite) TestTable5(c *gc.C) {

	var a []string = list()

	for _, s := range a {
		c.Assert(s, gc.Equals, "a")
	}
}

// issue
func (s *suite) TestTable6(c *gc.C) {

	var a = list()

	for _, s := range a {
		c.Assert(s, gc.Equals, "a")
	}
}

// issue
func (s *suite) TestTable7(c *gc.C) {

	var a []string
	a = list()

	for _, s := range a {
		c.Assert(s, gc.Equals, "a")
	}
}

// no issue
func (s *suite) TestMake(c *gc.C) {

	a := make([]string, 2)
	a = list()

	for _, s := range a {
		c.Assert(s, gc.Equals, "a")
	}
}

// no issue
func (s *suite) TestLen(c *gc.C) {

	a := list()
	c.Assert(len(a), gc.Equals, 2)

	for _, s := range a {
		c.Assert(s, gc.Equals, "a")
	}
}

// issue
func (s *suite) TestFunc(c *gc.C) {

	c.Assert(list(), gc.HasLen, 2)
	for _, s := range list() {
		c.Assert(s, gc.Equals, "a")
	}
}

// no issue
func (s *suite) TestInc(c *gc.C) {

	a := list()
	var i int
	for _, s := range a {
		i += 1
		c.Assert(s, gc.Equals, "a")
	}
}

// issue
func (s *suite) TestFunc2(c *gc.C) {
	content := ""
	for _, cont := range strings.Split(string(content), "\n") {
		c.Assert(cont, gc.Equals, "")
	}
}

func list() []string {
	return []string{
		"a",
		"a",
	}
}
