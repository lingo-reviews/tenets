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
		Name:        "worker_callback",
		Usage:       "no methods outside worker.Worker interface for callback workers",
		Description: "If your worker has any methods outside the worker.Worker interface, DO NOT use any of the above callback-style workers. Those methods, that need to communicate with the main goroutine, need to know that goroutine's state, so that they don't just hang forever.",
		Language:    "golang",
		SearchTags:  []string{"juju", "worker"},
	})
	return t
}
