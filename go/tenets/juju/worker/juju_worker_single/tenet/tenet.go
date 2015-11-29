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
		Name:        "worker_single",
		Usage:       "don't use worker/singular",
		Description: "If you're using worker/singular, you are quite likely to be doing it wrong, because you've written a worker that breaks when distributed. Things like provisioner and firewaller only work that way because we weren't smart enough to write them better; but you should generally be writing workers that collaborate correctly with themselves, and eschewing the temptation to depend on the funky layer-breaking of singular.",
		Language:    "golang",
		SearchTags:  []string{"juju", "worker"},
	})
	return t
}
