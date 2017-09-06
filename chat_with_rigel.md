# Notes from Chatting with Rigel

* Explain flattening tx - okay

* Discuss one account object
    * Both like state space separation

* Context and Permissions (Actors)
    * Hard to use different permissioning systems in different modules
    * Clean up actors to be more generic/extensible?

* Easier to start up
    * No router, no separation, just one handler
    * No worries about permissions or state-space applications
    * Show an optional use without permission system, isolation as intro example
    * Separation as advanced features

* Make it easier to ramp up from easy to advanced
    * Super simple, one handler, no decorators (eyes)
    * Add decorators (signing, nonce, etc.)
    * Add decorators that call into other modules (fees -> coin)
    * Add dispatcher for separate handlers (completely isolated)
    * Add permissions between modules

* Routing to use less magic
    * Not using the json name of the tx prefixed by module. ugh
        * On the router, explicity register all tx to each module/handler
        * One type for all tx for a given module, this is registered explicitly
        * In the module, it can store an interface in the tx and switch to determine which specific tx it is
    * Module name -> state-space prefix and permission name
        * Should explicitly wrap isolation on stuff,not done automagically?


```Go
// explicitly route the tx types
// no state-space isolation or permissioning
Router(
  {CoinTx{}, coin.Handler()},
  {StakeTx{}, stake.Handler()},
  {CounterTx{}, counter.Handler()},
)

// Isolate them all
Router(
  {CoinTx{}, Isolate("coin", coin.SafeHandler())},
  {StakeTx{}, Isolate("stake", stake.SafeHandler())},
  {CounterTx{}, Isolate("cntr", counter.SafeHandler())},
)
```

bottom-level modules don't need to care
if a module sets permissions and modifies other modules, it must
know how to do it - edit state space directly, or use IPC calls

can we set this in compile time?
eg. coin.Handler() will prefix itself, coin.SafeHandler() will let the Isolate prefix it

```
package coin

func Handler(isolated bool) sdk.Handler { ... }

func Handler() sdk.Handler { ... }
// do we need a different interface for handlers designed to work in isolated environments?
func SafeHandler() sdk.SafeHandler { ... }
```

## IPC

In non-isolated, they will just edit state-space directly.
In isolated environments, designed to scale more, the modules can trigger tx
on other modules.
For ease of reasoning, we consider all these to be done synchronously,
calling the router that then directs them.
Or is there another way to do it?
