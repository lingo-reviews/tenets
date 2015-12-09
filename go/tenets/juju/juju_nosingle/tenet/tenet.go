package tenet

import "github.com/lingo-reviews/tenets/go/dev/tenet"

// TODO(waigani) This is a stub tenet which only generates documentation, no
// source sniffing. Take a look at juju_worker_nostate for a more interesting
// tenet.

type workerTenet struct {
	tenet.Base
}

func New() *workerTenet {
	t := &workerTenet{}
	t.SetInfo(tenet.Info{
		Name:        "worker_nosingle",
		Usage:       "use the PeriodicWoker",
		Description: "Singletons are basically the same as global variables, except even worse, and if you try to make them responsible for goroutines they become more horrible still.",
		Language:    "golang",
		SearchTags:  []string{"juju", "worker"},
	})
	return t
}
