package server

import (
	"errors"
	"io"
	"time"

	"golang.org/x/net/context"

	"github.com/lingo-reviews/tenets/go/dev/api"
	"github.com/lingo-reviews/tenets/go/dev/tenet"
	"github.com/lingo-reviews/tenets/go/dev/tenet/log"
)

type server struct {
	// TODO(waigani) This is a little silly. Break the tenet concept down
	// more, and rename components.
	tenet tenet.Tenet
}

func (s *server) GetInfo(_ context.Context, _ *api.Nil) (*api.Info, error) {
	i := s.tenet.(tenet.BaseTenet).Info()
	return tenet.APIInfo(i), nil
}

// Options are passed in via .lingo or on the CLI.
func (s *server) Configure(_ context.Context, cfg *api.Config) (*api.Nil, error) {
	s.tenet.(tenet.BaseTenet).MixinConfigOptions(cfg.Options)
	return &api.Nil{}, nil
}

func (s *server) APIVersion(_ context.Context, _ *api.Nil) (*api.SchemaVersion, error) {
	return &api.SchemaVersion{}, nil
}

// Review reviews each file streamed from the client in sync and streams back
// all issues found.
func (s *server) Review(stream api.Tenet_ReviewServer) error {
	b := s.tenet.(tenet.BaseTenet)
	r := b.NewReview()
	r.StartReview()
	log.Println("review started")

	go func() {
		for {
			file, err := stream.Recv()
			if err == io.EOF {
				r.EndReview()
				log.Println("file stream closed. filesc closed.")
				// read done.
				return
			}
			if err != nil {
				log.Fatalf("failed to receive a file: %v", err)
			}
			log.Printf("server got file: %s", file.Name)
			if r.IsClosed() {
				log.Printf("review is closed. Not reviewing %s", file.Name)
				return
			}

			r.SendFile(file.Name)
		}
	}()

	// Read from our local issuesc and sent to client.
l:
	for {
		select {
		case issue, ok := <-r.Issues():
			if !ok && issue == nil {
				// issuesc is closed. We are done.
				break l
			}
			if err := stream.Send(tenet.APIIssue(issue)); err != nil {
				log.Fatalf("Failed to send an issue: %v", err)
			}

		case <-time.After(20 * time.Second):
			// it has taken too long to find an issue, something's wrong, end
			// the review.
			return errors.New("timed out waiting for issues")
		}
	}

	return nil
}
