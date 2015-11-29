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
		Name:        "worker_nonakedchan",
		Usage:       "never make a naked channel send/recieive",
		Description: "Basically never do a naked channel send/receive. If you're building a structure that makes you think you need them, you're most likely building the wrong structure.",
		Language:    "golang",
		SearchTags:  []string{"juju", "worker"},
	})
	return t
}
