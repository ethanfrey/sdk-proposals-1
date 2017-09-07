package two

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

func NewApp(dbPath string, logger log.Logger) App {
	return App{
		IAVLApp: abci.NewIAVLApp(dbPath, logger),
		handler: NewHandler(),
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

func init() {
	wire.RegisterInterface(
		Tx{},
		wire.ConcreteType{SetTx{}, 0x1},
		wire.ConcreteType{RemoveTx{}, 0x2},
	)
}

// LoadTx handles parsing the binary format, here is a simple intro to go-wire
func LoadTx(txBytes []byte) (Tx, error) {
	var tx TxWrapper
	err := data.FromWire(txBytes, &tx)
	return tx.Tx, err
}

// TODO: show basic tests

// This is the more general DeliverTx, using the errors packager
// and delegating to a handler.  Now, we can separate all logic in the
// handler, and this becomes a quite simple adaptor no matter how complex
// the logic
func (a *App) DeliverTx(txBytes []byte) abci.Result {
	tx, err := LoadTx(txBytes)
	if err != nil {
		return errors.Result(err)
	}

	// TODO: what is in a context?
	ctx := sdk.NewContext(a.logger)
	res, err := a.handler.DeliverTx(ctx, a.IAVLApp.DeliverState(), tx)
	if err != nil {
		return errors.Result(err)
	}
	return sdk.ToABCI(res)
}

// A more complete CheckTx to mirror the DeliverTx
func (a *App) CheckTx(txBytes []byte) abci.Result {
	tx, err := LoadTx(txBytes)
	if err != nil {
		return errors.Result(err)
	}

	// TODO: what is in a context?
	ctx := sdk.NewContext(a.logger)
	res, err := a.handler.CheckTx(ctx, a.IAVLApp.CheckState(), tx)
	if err != nil {
		return errors.Result(err)
	}
	return sdk.ToABCI(res)
}
