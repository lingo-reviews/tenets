package server

import (
	"fmt"
	"net"
	"os"
	"runtime"

	"google.golang.org/grpc"

	"github.com/juju/errors"
	"github.com/lingo-reviews/tenets/go/dev/api"
	"github.com/lingo-reviews/tenets/go/dev/tenet"
	"github.com/lingo-reviews/tenets/go/dev/tenet/log"
)

// Serve starts an RPC server hosting the api methods. It will first return
// its socket address.
func Serve(t tenet.Tenet) {

	lis, err := listener()
	if err != nil {
		log.Fatal("listen error:", err)
	}
	defer lis.Close()

	s := grpc.NewServer()
	api.RegisterTenetServer(s, newAPI(t))

	b := t.(tenet.BaseTenet)
	b.Init()

	// log non-fatal errors. TODO(waigani) make non-fatal errors avaliable to
	// client and do something a litte more interesting with them.
	go func() {
		errorsc := b.Errors()
		for err := range errorsc {
			log.Printf("%v", err)
		}
	}()
	s.Serve(lis)
}

func newAPI(t tenet.Tenet) *server {
	return &server{
		tenet: t,
	}
}

func listener() (net.Listener, error) {
	if os.Getenv("LINGO_CONTAINER") != "" {
		return net.Listen("tcp", ":8000")
	}

	switch runtime.GOOS {
	case "darwin":
		return localTcpListener()

	case "linux", "freebsd":
		return localUnixListener()

	default:
		panic("Unsupported OS.")
	}
}

func localTcpListener() (net.Listener, error) {
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return nil, errors.Trace(err)
	}

	fmt.Println(lis.Addr().String())

	return lis, nil
}

func localUnixListener() (net.Listener, error) {
	laddr := net.UnixAddr{Net: "unix"} // Name: "@001fc" use this to debug.
	lis, err := net.ListenUnix("unix", &laddr)
	if err != nil {
		return nil, errors.Trace(err)
	}
	socketAddr := lis.Addr().String()

	// Print the socket address so the client knows where to connect to.
	fmt.Println(socketAddr)

	return lis, nil
}
