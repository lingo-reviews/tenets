package main

import (
	"github.com/lingo-reviews/tenets/go/dev/server"
	"github.com/lingo-reviews/tenets/go/tenets/license/tenet"
)

func main() {
	server.Serve(tenet.New())
}
