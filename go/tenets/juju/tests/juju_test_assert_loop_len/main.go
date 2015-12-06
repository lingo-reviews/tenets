package main

import (
	"github.com/lingo-reviews/tenets/go/dev/server"
	"github.com/lingo-reviews/tenets/go/tenets/juju/tests/juju_test_assert_loop_len/tenet"
)

func main() {
	server.Serve(tenet.New())
}
