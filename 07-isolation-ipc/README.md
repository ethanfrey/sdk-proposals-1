# Scalable multi-modules and IPC

Here we show how to do the last example in a safer method. We introduce
state-space isolation, as well as permission isolation. We show how each
module can send txs to the other modules, with custom payloads and possibly
with custom permissions, to trigger specific actions.

We will re-architecture the previous example to do so "safely", we notice
how we decouple the pieces and make reasoning about large programs easier.
We also mention that this is a bit more complex and can be skipped for
the first mvp, when there are few modules and few tx types. But should be
used for larger systems, and especially when importing third-party modules
to perform actions, like staking.
