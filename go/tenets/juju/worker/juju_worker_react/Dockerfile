FROM golang
# FROM golang:onbuild # <---- This line is all that will be needed.

ENV LINGO_CONTAINER true

## ---------
# The following is only needed while lingo libs are privately hosted on
# bitbucket. Once they are published, 'FROM golang:onbuild' is all we need
# here. But for now we need to manually checkout the repos into the paths
# copied below.

COPY . /go/src/github.com/lingo-reviews
COPY tenets/go/tenets/juju/worker/juju_worker_react /go/src/app
WORKDIR /go/src/app
RUN go get -v -d
RUN go install -v
ENTRYPOINT /go/bin/app
## ----------

# This info is used for searching for tenet images.
LABEL reviews.lingo.name="lingoreviews/juju_worker_react" \
