package three

import (
	wire "github.com/tendermint/go-wire"
	"github.com/tendermint/go-wire/data"
	"github.com/tendermint/tmlibs/log"

	sdk "github.com/cosmos/cosmos-sdk"
	"github.com/cosmos/cosmos-sdk/abci"
	"github.com/cosmos/cosmos-sdk/errors"
)

// App represents an ABCI app build on the persistence model
type App struct {
	// IAVLApp looks like sdk.app.Basecoin without the DeliverTx and CheckTx
	// methods, but handles all queries, commit, handshaking, etc and
	// provides a presistent IAVLTree as backing data
	*abci.IAVLApp
	handler sdk.Handler
	logger  log.Logger
}

// DecoratedHandler returns the same handler, but surrounded by a number
// of decorators,
func DecoratedHandler() sdk.Handler {
	// define a decorator list here to be applied together
	// combines many decorators into one
	sdk.ChainDecorators(
		// these read no data, but add loging and panic recovery
		logger.Decorator,
		recovery.Decorator,
		// these just inspect the packet, possibly return error or modify context
		auth.Decorator,
		chain.Decorator,
		// this actually checks the state space and modifies it (isolated?)
		nonce.Decorator,
	).WithHandler(
		// and they wrap the same handler in the end
		NewHandler(),
	)

}

func NewApp(dbPath string, logger log.Logger) App {
	return App{
		IAVLApp: abci.NewIAVLApp(dbPath, logger),
		handler: DecoratedHandler(),
		logger:  logger,
	}
}

// Tx is an interface to group all transactions we support.
// simplest way to multiplex with go-wire
// doesn't support json, but yeah, should work for binary only
type Tx interface {
	ValidateBasic() error
}

// TxWrapper is a struct to hold an interface, so we have a place to store
// it when using the serialization library, can be ignored otherwise
type TxWrapper struct {
	Tx `json:"unwrap"`
}

// TODO: introduce go-data and json magic?
func init() {
	wire.RegisterInterface(
		Tx{},
		wire.ConcreteType{SetTx{}, 0x1},
		wire.ConcreteType{RemoveTx{}, 0x2},
	)
}

// LoadTx handles parsing the binary format, here is a simple intro to go-wire
func LoadTx(txBytes []byte) (SignedMsg, error) {
	var msg SignedMsg
	err := data.FromWire(txBytes, &msg)
	return msg, err
}

// Msg is all data we want to have signed (for Tx and decorators)
// All decorators are here on one level and statically type-checked
// in the construction phase
type Msg struct {
	// use TxWrapper here for better serialization
	// how does the handler receive this?  via GetTx()????
	// do we embed it and have TxWrapper expose some `Unwrap() Tx` method??
	// other ideas
	Tx TxWrapper `json:"tx"`

	Chain chain.Data `json:"chain"` // this is just a string for the id
	Nonce nonce.Data `json:"nonce"`
	// Fee is a bit more complex as it modifies state, will add in later demo
	// Fee   fee.Data   `json:"fee"`
}

// GetTx is maybe used by the handler to get the info?? Or a better way???
// Ideally the Handler can receive the desired tx and parse it regardless
// of whether it was decorated or not.
//
// maybe the ChainDecorator()...WithHandler() method will call GetTx on the
// message and then pass that into the handler????
func (m Msg) GetTx() Tx {
	return m.Tx
}

// GetNonce is used by the nonce.Middleware to check the data
func (m Msg) GetNonce() nonce.Data {
	return m.Nonce
}

// GetChain is used by the chain.Middleware to check the data
// (verify we are on the same chain)
func (m Msg) GetChain() chain.Data {
	return m.Chain
}

// SignedMsg wraps the Msg with the signature,
// this is separated out, so we can easily determine binary format
// of the message that needs to be signed, as well as the message
// together with it's signature
type SignedMsg struct {
	// this is embedded, so we eaisly expose GetFee, GetNonce, Tx etc.
	Msg  `json:"msg"`
	Sigs auth.Signature `json:"sigs"`
}

// SignBytes is what we sign, right?
func (s SignedMsg) SignBytes() []byte {
	return wire.BinaryBytes(s.Msg)
}

// Bytes is the total info to store in the
func (s SignedMsg) Bytes() []byte {
	return wire.BinaryBytes(s)
}

// Validate makes sure it is signed properly
func (s SignedMsg) Validate() error {
	return nil
}

// DeliverTx calls all decorators, then the handler
func (a *App) DeliverTx(txBytes []byte) abci.Result {
	msg, err := LoadTx(txBytes)
	if err != nil {
		return errors.Result(err)
	}

	// first I had the decorators totally separate from the handler,
	// but they need to really wrap the call, so they can act on the
	// input and the output

	ctx := sdk.NewContext(a.logger)
	res, err := a.handler.DeliverTx(ctx, a.IAVLApp.DeliverState(), tx)
	if err != nil {
		return errors.Result(err)
	}
	return sdk.ToABCI(res)
}

// A more complete CheckTx to mirror the DeliverTx
func (a *App) CheckTx(txBytes []byte) abci.Result {
	// TODO
}
