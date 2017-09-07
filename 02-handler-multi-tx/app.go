package two

import (
	"log"

	"github.com/cosmos/cosmos-sdk/abci"
	wire "github.com/tendermint/go-wire"
	"github.com/tendermint/go-wire/data"
)

// App represents an ABCI app build on the persistence model
type App struct {
	// IAVLApp looks like sdk.app.Basecoin without the DeliverTx and CheckTx
	// methods, but handles all queries, commit, handshaking, etc and
	// provides a presistent IAVLTree as backing data
	*abci.IAVLApp
}

func NewApp(dbPath string, logger log.Logger) App {
	return App{
		IAVLApp: abci.NewIAVLApp(dbPath, logger),
	}
}

// SetTx sets the data at the key
type SetTx struct {
	Key   data.Bytes
	Value data.Bytes
}

// RemoveTx deletes the contents at a key
type RemoveTx struct {
	Key data.Bytes
}

// simplest way to multiplex with go-wire
// doesn't support json, but yeah, should work for binary only
type Tx interface{}

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

// Bytes encodes this, must be symetric with LoadTx
// Shows how to encode for the client side
func (s SetTx) Bytes() []byte {
	// note that this type and the recevier must differ by exactly one points
	// eg. if you serialize SetTx, you pass *SetTx to the loader
	// if you serialize *SetTx, you must pass **SetTx to the loader
	// we must also encode it wrapped in order to get the type byte properly
	return wire.BinaryBytes(TxWrapper{s})
}

// Bytes encodes this, must be symetric with LoadTx
func (r RemoveTx) Bytes() []byte {
	return wire.BinaryBytes(TxWrapper{r})
}

// TODO: show switching on tx type
// TODO: show basic tests

// Show the minimal way to handle delivertx
func (a *App) DeliverTx(txBytes []byte) abci.Result {
	tx, err := LoadTx(txBytes)
	if err != nil {
		return abci.NewError(err.Error())
	}
	db := a.IAVLApp.DeliverDB()
	db.Set(s.Key, s.Value)
	return abci.Result{}
}

// Show the minimal way to handle checktx
func (a *App) CheckTx(txBytes []byte) abci.Result {
	tx, err := LoadTx(txBytes)
	if err != nil {
		return abci.NewError(err.Error())
	}
	db := a.IAVLApp.CheckDB()
	db.Set(s.Key, s.Value)
	return abci.Result{}
}
