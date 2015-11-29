This Juju Go tenet targets the juju/workers package.

A *state.State should not be needed to construct a worker.

If you're passing a \\*state.State into your worker, you are almost certainly doing it wrong. The layers go worker->apiserver->state, and any attempt to skip past the apiserver layer should be viewed with *extreme* suspicion.
