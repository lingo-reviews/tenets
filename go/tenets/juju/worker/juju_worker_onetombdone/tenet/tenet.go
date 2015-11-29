package tenet

import "github.com/lingo-reviews/tenets/go/dev/tenet"

type workerTenet struct {
	tenet.Base
}

// TODO(waigani) This is a stub tenet which only generates documentation, no
// source sniffing. Take a look at juju_worker_nostate for a more interesting
// tenet.

func New() *workerTenet {
	t := &workerTenet{}
	t.SetInfo(tenet.Info{
		Name:        "worker_onetombdone",
		Usage:       "one .tomb.Done() call",
		Description: "If it's possible for your worker to call .tomb.Done() more than once, or less than once, you are definitely doing it very very wrong indeed.",
		Language:    "golang",
		SearchTags:  []string{"juju", "worker"},
	})
	return t
}
