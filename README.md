# LINGO REVIEWS ALPHA

## House keeping

Welcome to the closed Alpha! There will be bugs. Please help us catch them and
open an issue. We also welcome feature requests and feedback on lingo's
usefulness and potential.

Okay, with that out of the way, let's get started.

## Install

To install, run:

```bash
wget http://lingo.reviews/lingo.zip; unzip lingo.zip
```

Then place lingo in your PATH:

```bash
cp lingo /usr/local/bin/
```

## First Run

### Docker Quick Start

If you have Docker installed:

```bash
# Find some source code to review.
cd go/tenets/license/tenet/example

# Review the code.
lingo review

# Read the tenet documentation for this code.
lingo docs

```

When lingo reviews, it looks for a .lingo file in the current or parent
directories. If those tenets use a docker driver (default) and no local docker
image is found, lingo goes and gets it. The first time you pull a docker
tenet, it will pull the tenet base images. This means future tenet pulls will
be much quicker.


Next, start without a .lingo file:

```bash
cd go/tenets/simpleseed/example

# This will write a .lingo file.
lingo init

# List avaliable tenets on hub.docker.com:
docker search lingoreviews

# Add the simpleseed example:
lingo add lingoreviews/simpleseed

# Pull down the images from hub.docker.com:
lingo pull

# If you didn't pull, review will do it for you.

# Review the code, this time we'll keep some output at the end:
lingo review --output-format --json-pretty

```

Notes: Tenets can be pulled from any docker repository. A better tenet search
UI is in the pipeline.

Lingo will prompt you to open each issue. Supported editors are: vi, vim,
emacs, nano and subl. To skip the confirm steps, use --keep-all.

### Binary Quick Start

All the other example folders under go/tenets use the binary driver. To build
them all at once, cd into the root of go/tenets and run:

```bash
lingo build binary --all
```

You'll see the following output:

```bash
$ lingo build binary --all
Building Go binary: [~/.lingo_home/tenets/lingoreviews/juju_nosingle]
Building Go binary: [~/.lingo_home/tenets/lingoreviews/imports]
...
Building Go binary: [~/.lingo_home/tenets/lingoreviews/unused_arg]
Building Go binary: [~/.lingo_home/tenets/lingoreviews/juju_worker_periodic]
binary 17 / 17 [========================================================] 100.00 % 12s
Success! All binary tenets built.
```

`cd` into any example folder and run `lingo review`. In a similar fashion, you
can `lingo build docker --all` to build local copies of all the docker
tenets.To add the binary drivers, we need to specify the driver when we add
it:

```bash
lingo add lingoreviews/simpleseed --driver binary
```

Otherwise, the driver will default to "docker". By default, binary tenets are
installed in ~/.lingo_home/tenets/[owner]/[name]. This can be overridden with
the LINGO_BIN environment variable.

## Bash Auto-Complete

Run `lingo --generate-bash-completion` to enable commands to auto-complete.
Commands such as `add` and `info` will auto-complete with the names of built
binary tenets.

## Options

Some tenets take options. To view their available options run:

```bash
lingo info <tenet-name>
```

The imports tenet, for example, takes a blacklist_regex option, here's an
example of setting it:

```bash
lingo add lingoreviews/imports --options blacklist_regex=".*/State"
```

## Writing a Tenet

# Go

Start [here](https://github.com/lingo-reviews/tenets/tree/master/go/dev). The
`go/tenets` directory also has a variety of examples of tenets in Go. Copy
`go/tenets/simpleseed` to get started.

# Other languages

The api.proto file in the root of this repoistory can be used to generate the
tenet API libs in C, C++, Java, Go, Node.js, Python, Ruby, Objective-C, PHP
and C#. Visit grpc.io to learn more.


## LAAS - Lingo As As Service

Go to www.lingo.reviews/dashboard to hook lingo up to your github repoistory.
Add .lingo files to your repoistory and the lingo bot will review every new
pull request.
