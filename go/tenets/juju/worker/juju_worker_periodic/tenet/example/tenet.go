package tenet

import "github.com/lingo-reviews/tenets/go/dev/tenet"

type workerTenet struct {
	tenet.Base
}

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
