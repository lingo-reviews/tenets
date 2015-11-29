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
		Name:        "worker_leakytomb",
		Usage:       "tomb.ErrDying should not leak",
		Description: "If you're letting tomb.ErrDying leak out of your workers to any clients, you are definitely doing it wrong -- you risk stopping another worker with that same error, which will quite rightly panic (because that tomb is not yet dying).",
		Language:    "golang",
		SearchTags:  []string{"juju", "worker"},
	})
	return t
}
