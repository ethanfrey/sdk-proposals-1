## Unsafe Module Interaction

Building on the last topic, we have multiple modules, this time, we want
to allow them to interact.

To have this make sense, we introduce the coin module from the standard
library as the base module.  We also decide to add "pay to vote" functionality
that subtracts from your account and adds to some pot.

We can then also show taking fees in the middleware as another form of
cross-module interaction.

All modules need to prefix their data, and they all directly modify the
internal data of each other.  It is clear we can combine code this way, but
it is fragile, and difficult to reason about all possible states or the
invariants as any module can modify them (no longer have a fixed total tokens).
