package tenet

// This file maps objects to api structs that can be passed on the wire.

import (
	"fmt"
	"go/token"

	"github.com/lingo-reviews/tenets/go/dev/api"
)

func APIInfo(i *Info) *api.Info {
	apiInfo := &api.Info{
		Name:        i.Name,
		Usage:       i.Usage,
		Description: i.Description,
		Version:     i.Version,
		Tags:        i.tags,
		Metrics:     i.metrics,
		Language:    i.Language,
	}

	apiInfo.Options = make([]*api.Option, len(i.Options))
	for i, o := range i.Options {
		apiInfo.Options[i] = &api.Option{
			Name:  o.name,
			Usage: o.usage,
			Value: *o.value,
		}
	}

	return apiInfo
}

func APIIssue(i *Issue) *api.Issue {
	issue := &api.Issue{
		Name:      i.Filename(),
		Position:  apiIssueRange(i.Position),
		Comment:   i.Comment,
		CtxBefore: i.CtxBefore,
		LineText:  i.LineText,
		CtxAfter:  i.CtxAfter,
		Link:      i.Link,
		Metrics:   apiMetrics(i.Metrics),
		Tags:      i.Tags,
		NewCode:   i.NewCode,
		Patch:     i.Patch,
	}
	if i.Err != nil {
		issue.Err = i.Err.Error()
	}
	return issue
}

// TODO(waigani) can protobuf support interfaces?
func apiMetrics(oldMap map[string]interface{}) map[string]string {
	newMap := make(map[string]string)
	for k, v := range oldMap {
		newMap[k] = fmt.Sprintf("%v", v)
	}
	return newMap
}

func apiIssueRange(r *issueRange) *api.IssueRange {
	return &api.IssueRange{
		Start: apiPosition(r.Start),
		End:   apiPosition(r.End),
	}
}

func apiPosition(p token.Position) *api.Position {
	return &api.Position{
		Filename: p.Filename,
		Offset:   int64(p.Offset),
		Line:     int64(p.Line),
		Column:   int64(p.Column),
	}
}
