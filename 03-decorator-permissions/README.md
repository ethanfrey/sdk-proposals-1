## Decorators and Permissions

We show how we can use decorators, by implementing a custom tx type that
provides all info for the various decorators, and then wrapping them
in a stack around the actual tx (set, update, and delete).

We modify the logic, so on creation we set an "owner", if the data is set,
only the owner can update and delete it.

We also show how to build the corresponding tx from command line flags on
the client side. Still using no proofs.

We explain how permissions work in general, but just concentrate on
the auth permission for now.
