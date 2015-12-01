# LINGO REVIEWS ALPHA

## House keeping

Welcome to the closed Alpha! There will be bugs. Please help us catch them and open an issue. We also welcome feature requests and feedback on lingo's usefulness and potential.

Okay, with that out of the way, let's get started.

## First Run

### Docker Quick Start

If you have Docker installed, try the following:

```bash
# Find some source code to review.
cd go/tenets/license/tenet/example

# Review the code.
lingo review

# Read the tenet documentation for this code.
lingo docs

```

When lingo reviews, it looks for a .lingo file in the current or parent directories. If those tenets use a docker driver (default) and no local docker image is found, lingo goes and gets it. The first time you pull a docker tenet, it will pull the tenet base images. This means future tenet pulls will be much quicker.


This time, let's start without a .lingo:

```bash
cd go/tenets/simpleseed/example

# This will write a .lingo file.
lingo init

# List avaliable tenets on hub.docker.com:
docker search lingoreviews

# Add the license example:
lingo add lingoreviews/simpleseed

# Pull down the images from hub.docker.com:
lingo pull

# If you didn't pull, review will do it for you.

# Review the code, this time we'll keep some output at the end:
lingo review --output-format --json-pretty

```

Notes: Tenets can be pulled from any docker repository. A better tenet search UI is in the pipeline.

Have a play. You'll see lingo prompts you to open the issue. Supported editors are: vi, vim, emacs, nano and subl. If you want to skip the confirm steps, use --keep-all.

### Binary Quick Start

All the other example folders under go/tenets use the binary driver. To build them all at once, cd into the root of go/tenets and run:

```bash
lingo build binary --all
```

You'll see the following output:

```bash
$ lingo build binary --all
Building Go binary: [/home/jesse/.lingo_home/tenets/lingoreviews/juju_nosingle]
Building Go binary: [/home/jesse/.lingo_home/tenets/lingoreviews/imports]
...
Building Go binary: [/home/jesse/.lingo_home/tenets/lingoreviews/unused_arg]
Building Go binary: [/home/jesse/.lingo_home/tenets/lingoreviews/juju_worker_periodic]
binary 17 / 17 [==============================================================================================================] 100.00 % 12s
Success! All binary tenets built.
```

Now, cd into any example folder and run `lingo review`. In a similar fashion, you can `lingo build docker --all` to build local copies of all the docker tenets.


The following is a fully functional tenet. You'll find it in go/tenets/simpleseed. 

```go

package main

import (
	"go/ast"

	"github.com/lingo-reviews/go/tenets/dev/server"
	"github.com/lingo-reviews/go/tenets/dev/tenet"
)

type commentTenet struct {
	tenet.Base
}

func main() {
	c := &commentTenet{}
	c.SetInfo(tenet.Info{
		Name:     "simpleseed",
		Usage:    "every comment should be awesome",
		Language: "Go",
		Description: "description",
	})

	issue := c.RegisterIssue("sucky_comment", tenet.AddComment("this comment could be more awesome"))
	c.SmellNode(func(r tenet.Review, comment *ast.Comment) error {
		if comment.Text != "// most awesome comment ever" {
			r.RaiseNodeIssue(issue, comment)
		}
		return nil
	})

	server.Serve(c)
}

```

To build, cd into go/tenets/simpleseed and run `lingo build`. You can now add simpleseed to any project. Because it's not docker, we need to specify the driver when we add it:

```bash
lingo add lingoreviews/simpleseed --driver binary
```

If you run `lingo --generate-bash-completion` commands such as `add` and `info` will autocomplete with all the built binary names. 

`lingo build` looks for a .lingofile for instructions on how to build the tenet. This is the simpleseed .lingofile:

```toml
language = "go"
owner = "lingoreviews"
name = "simpleseed"

[docker]
  build=false
  overwrite_dockerfile=true
```

By default, binary tenets are installed in ~/.lingo_home/tenets/[owner]/[name]. This can be overridden with the LINGO_BIN environment variable.

You'll note .lingofile sets docker build to false. If you remove that and run build again, lingo will build both a docker image and a binary. Until github.com/lingo-reviews/tenets is published, you'll have to manually git clone it into ~/go/src/github.com/lingo-reviews/tenets before building a docker image.

As you'll likely be writing and updating lots of tenets (we hope!), lingo build provides a --all flag. To build all the binary tenets in go/tenets in one hit (recommended):

```bash
cd go/tenets
lingo build binary --all
```

This looks for all .lingofiles under the current directory and attempts to build all the tenets.

## Options

Some tenets take options. To view their available options run:

```bash
lingo info
```

The imports tenet takes a blacklist_regex option, here's an example of setting it:

```bash
lingo add lingoreviews/imports --options blacklist_regex=".*/State"
```


## Bash completion

To enable bash completion, run the following:

```bash
lingo --generate-bash-completion
```

Commands and some arguments will now autocomplete, in particular installed binary tenet names.


## Writing a Tenet

# Go
Start with go/dev/README.md. The `go/tenets` directory also has a variety of examples of tenets in Go. Copy one of those to get started.

# Other languages
The api.proto file in the root of this file can be used to generate the tenet API libs in C, C++, Java, Go, Node.js, Python, Ruby, Objective-C, PHP and C#. Visit grpc.io to learn more.


### LAAS - Lingo As As Service

Go to www.lingo.reviews/dashboard to hook lingo up your github repoistory. Add a .lingo file to your repo and the lingo bot will review every new pull request.