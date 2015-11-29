package tenet

// information about the tenet
type Info struct {
	Name        string
	Usage       string
	Description string
	Language    string
	SearchTags  []string
	// Tags and metrics need to be registered and cannot be set here directly.
	tags    []string
	metrics []string
	Options []*option
	Version string
}

func (b *Base) Info() *Info {
	return b.info
}

func (b *Base) SetInfo(i Info) Tenet {
	b.info = &i
	return b
}
