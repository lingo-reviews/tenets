package main

import (
	"github.com/lingo-reviews/tenets/go/dev/server"
	"github.com/lingo-reviews/tenets/go/tenets/juju/worker/juju_worker_tomb/tenet"
)

func main() {

	// Serve up the tenet.
	server.Serve(tenet.New())
}
