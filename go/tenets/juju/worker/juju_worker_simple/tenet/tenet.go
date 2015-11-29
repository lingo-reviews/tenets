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
		Name:        "worker_simple",
		Usage:       "a basic worker should use the SimpleWorker",
		Description: "If you really just want to run a dumb function on its own goroutine, use worker.NewSimpleWorker.",
		Language:    "golang",
		SearchTags:  []string{"juju", "worker"},
	})
	return t
}
