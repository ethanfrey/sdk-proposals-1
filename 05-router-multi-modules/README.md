## Routing multiple modules

Here we start to look at composing larger applications.  When you want to
split functionality, or import multiple pieces from other repos, we need a
way to configure them to work together.

The first item to look at is the router, and registering so txs are routed
to the module that can handle them (and is registered for them).

In doing so we also implement a counter that can be incremented or
decremented by one per tx. Each tx must be signed, but no more permission
is needed.

We also look at the genesis file parsing here to provide an initial state.

We show how different modules have their own spaces and can interact without
conflicts by using prefixes.
