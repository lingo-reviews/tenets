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
		Name:        "worker_tomb",
		Usage:       "custom worker should use tomb.Tomb",
		Description: `If you're writing a custom worker, and not using a tomb.Tomb, you are almost certainly doing it wrong. Read the blog post, or just the code -- it's less than 200 lines and it's about 50\% comments.`,
		Language:    "golang",
		SearchTags:  []string{"juju", "worker"},
	})
	return t
}
