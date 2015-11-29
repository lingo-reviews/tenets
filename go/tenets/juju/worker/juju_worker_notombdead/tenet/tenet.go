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
		Name:        "worker_notombdead",
		Usage:       "don't use .tomb.Dead()",
		Description: "If you're using .tomb.Dead(), you are very probably doing it wrong -- the only reason (that I'm aware of) to select on that .Dead() rather than on .Dying() is to leak inappropriate information to your clients. They don't care if you're dying or dead; they care only that the component is no longer functioning reliably and cannot fulfil their requests. Full stop. Whatever started the component needs to know why it failed, but that parent is usually not the same entity as the client that's calling methods.",
		Language:    "golang",
		SearchTags:  []string{"juju", "worker"},
	})
	return t
}
