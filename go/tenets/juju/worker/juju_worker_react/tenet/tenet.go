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
		Name:        "worker_react",
		Usage:       "reactive worker to use Notify or Strings worker",
		Description: "If you want to react to watcher events, you should probably use worker.NewNotifyWorker or worker.NewStringsWorker.",
		Language:    "golang",
		SearchTags:  []string{"juju", "worker"},
	})
	return t
}
