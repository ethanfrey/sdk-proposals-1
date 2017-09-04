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



## Remove magic
