# LINGO TENETS

This is a repository of tenets used by the [Lingo tool](https://github.com/lingo-reviews/lingo) 
to manage code and product quality. Each tenet captures a piece of best
practice and ensures documentation and implementation align with it.


## Writing a Tenet

### Go

Currently Lingo has best support for Go. Start [here](https://github.com/lingo-reviews/tenets/tree/master/go/dev). The
`go/tenets` directory has a variety of examples of tenets in Go. Copy
`go/tenets/simpleseed` to get started.

### Python

Support for Python tenets is under way. Get in touch (hello@lingo.reviews) if
you would like to help.

### Other languages

The api.proto file in the root of this repoistory can be used to generate the
tenet API libs in C, C++, Java, Go, Node.js, Python, Ruby, Objective-C, PHP
and C#. Visit grpc.io to learn more.

While each language has its own idiosyncrasies, there are three components
that every tenet will need: an RPC API (auto generated), a server to serve up
that API and helper libs to analyse code and send back results via the API. If
you would like to add support for one of these languages, please open an issue
to let us know that you are doing so - so we do not double up.
