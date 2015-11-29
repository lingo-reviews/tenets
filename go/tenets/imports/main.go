package main

import (
	"github.com/lingo-reviews/tenets/go/dev/server"
	"github.com/lingo-reviews/tenets/go/tenets/imports/tenet"
)

func main() {

	// Serve up the tenet.
	server.Serve(tenet.New())
}
