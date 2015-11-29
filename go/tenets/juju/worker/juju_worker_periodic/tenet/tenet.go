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
		Name:        "worker_periodic",
		Usage:       "use the PeriodicWoker",
		Description: "If you just want to do something every <period>, use worker.NewPeriodicWorker",
		Language:    "golang",
		SearchTags:  []string{"juju", "worker"},
	})
	return t
}
