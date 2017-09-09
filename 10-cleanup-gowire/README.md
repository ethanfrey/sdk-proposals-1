# Cleanup go-wire

While working out the early examples, and trying to see how they would work
for a programmer new to our stack, I realized go-wire (and go-data) is a
big source of confusion, and needs to be made more accessible.

I agree with the functionality, but it needs to be easy to use.  And while
I would prefer a more standard encoding scheme, nothing else does seem
to meet the needs (please prove me wrong), but if this is our standard,
it should be well spec'd and well documented, and easy to implement in
various languages, including the type-byte stuff.

I opened an issue here: https://github.com/tendermint/go-wire/issues/23

But let us sketch out how we could introduce this to a novice user.
Also some other WTF cases that need to be clarified and fixed, like this:

```
func main() {
    acct := &Account{Name: "foo"}
    bs := wire.BinaryBytes(acct)
    var input Account
    err := wire.ReadBinaryBytes(bs, &input)
    // returns an error as we pass in *Account, not **Account
    // what other lib does this????
}
```
