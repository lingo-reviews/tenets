// Copyright 2015 Jesse Meek.
// Licensed under the AGPLv3, see LICENCE file for details.

package tenet

import (
	"go/ast"
	"regexp"

	"github.com/juju/errors"
	"github.com/lingo-reviews/tenets/go/dev/tenet"
)

type importsTenet struct {
	tenet.Base
}

func New() *importsTenet {
	t := &importsTenet{}
	t.SetInfo(tenet.Info{
		Name:        "imports",
		Usage:       "police the imports of a package",
		Description: "imports matching the following regex are restricted: \"{{.blacklist_regex}}\"",
		SearchTags:  []string{"import"},
		Language:    "go",
	})

	blacklisted := t.RegisterOption("blacklist_regex", "", "a regex to filter imports against")
	issue := t.RegisterIssue("blacklisted_import",
		tenet.AddComment(`This package should not be bringing in {{.importName}}`),
	)
	// Metrics and tags could also be registered at this point to help track
	// this type of issue.

	t.SmellNode(func(r tenet.Review, imp *ast.ImportSpec) error {
		importName := imp.Path.Value
		m, err := regexp.MatchString(*blacklisted, importName)
		if err != nil {
			return errors.Trace(err)
		}

		if m {
			r.RaiseNodeIssue(issue, imp, tenet.CommentVar("importName", importName))
		}
		return nil
	})

	return t
}
