# Cosmos SDK proposals

Here are some ideas for improving the sdk, especially after feedback from Jae.

## Naming

Do we want to do this?

* modules -> domains
* middleware -> decorators (wrappers)?

## Simplify layered tx

Right now one tx embeds another tx, embeds another.  For any givne app, there
is only one correct ordering, but it can be confusing to determine this.  There
was also negative feedback on the size and complexity of this system.

Instead of the current system with wrapper tx for "middleware", every app can
define a custom message type, which includes all data for the middleware it is
compiled to use, as well as a slot for one tx (SendTx, CreateRole, StakeToken)
that executes business logic after the checks.

This means that the parsing code in DeliverTx and CheckTx in the abci app
bindings must be overriden per application (not too hard), but does it also
mean the middleware need to know about the exact type to use it?

We can change the Middleware interface and have it accept `msg interface{}`
or something like that, which it will then check holds the appropriate data
for this Middleware `GetFee() Coins`. Every app can implement their own top
level message types, and then just implement these interfaces for the middleware
they want to use.

Example:

```
type BasecoinMsg struct {
    Fee types.Coin
    Nonce uint64
    Tx Tx
}

func (m BasecoinMessage) GetFee() types.Coin {
    return m.Fee
}

type Handler interface {
    DeliverTx(ctx Context, msg interface{}, store SimpleDB) (Result, error)
}

type Middleware interface {
    DeliverTx(ctx Context, msg interface{}, store SimpleDB, next Handler) (Result, error)
}

type FeeGetter interface {
    GetFee() types.Coin
}

func (FeeMiddleware)DeliverTx(ctx Context, msg interface{},
    store SimpleDB, next Handler) (Result, error) {

    tx, ok := msg.(FeeGetter)
    if !ok {
        return Result{}, errors.New("Message must implement FeeGetter")
    }
    fee := tx.GetFee()
    // .... work with the fee as normal
}

```

Another benefit to this is that the basecoin message can make some nice
optimizations.  Like, let's say we can specify which account (or multisig) the
nonce belongs to, but if there is one signer the nonce can be safely assumed
to refer to that account.  We can handle this easily in the `BasecoinMessage`
as it has the Nonce and Signature info stored and type-safely accessible and
could implement the `GetNonce` to do this optimization invisibly to the rest
of the system.

```
type Nonce struct {
    Sequence uint64
    Address []byte
}

type NonceGetter interface {
    GetNonce() Nonce
}

func (m BasecoinMsg) GetNonce() Nonce {
    res := m.Nonce
    if len(res.Address) == 0 {
        signers := m.GetSigners()
        if len(signers) == 1 {
            res.Address = signers[0]
        }
    }
    return res
}
```

(TODO: work out example code in the simplify_tx package)

## Make context more generic

The current context we push into the application handler looks like this:

```
// Context is an interface, so we can implement "secure" variants that
// rely on private fields to control the actions
type Context interface {
    log.Logger
    WithPermissions(perms ...Actor) Context
    HasPermission(perm Actor) bool
    GetPermissions(chain, app string) []Actor
    IsParent(ctx Context) bool
    Reset() Context
    ChainID() string
    BlockHeight() uint64
}
```

Yes, it is ugly and just growing.  How about we break it into a few pieces,
use context.Context from the go standard library and some helper functions.

[context.WithValue](https://godoc.org/context#example-WithValue) will allow us
to add plenty of information in an immutible way. Now, let's imagine we replace
the first argument with a simple `context.Context`, how might this look?

```
// look mom, no collisions in the key space.
type logKeyType int
const logKey logKeyType = 1

func SetLogger(ctx context.Context, logger log.Logger) context.Context {
    return context.WithValue(ctx, logKey, logger)
}

func GetLogger(ctx context.Context) logger log.Logger {
    return ctx.Value(logKey)
}
```

`Reset` and `IsParent` should be added to some general package (base?).
`SetChainID`, `GetChainID`, `GetBlockHeight`... should be added to another
app package.

We could imagine all permisison stuff to be done in one `auth` package that
provided those interfaces. Any module/domain that wanted to use that convention
could use it. Anyone that wanted a different convention could use their own.
You would just need to make sure all modules in your app used the same
convention (so it may evolve into two or three different ways of building apps
that are slightly incompatible, but this is flexibility)

Thoughts on this?  Maybe there is a different way (like adding more interfaces
to context)?  But this seems more like standard lib.  If you like it, I will
sketch this (and the above section) out more.

## Remove more magic

Although Frey loves magic, decorators, reflection, and meta-classes... Jae is
most likely correct that they do not make an easy-to-use, and accessible
framework.

Let's look at the magic that is used and see how to address it:

* State sandbox based on module name
* Permission sandbox based on module name (and some ibc flag)
* Routing of tx based on parsing the go-data name (this one is ugly)

Also there but being dealt with:

* Tx layering (addressed above)
* Viper sprinkled everywhere (to be addressed in develop soon anyway)

Thinking about these....
