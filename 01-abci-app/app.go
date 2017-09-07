package one

import (
	wire "github.com/tendermint/go-wire"
	"github.com/tendermint/go-wire/data"
	"github.com/tendermint/tmlibs/log"

	"github.com/cosmos/cosmos-sdk/abci"
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

// SetTx is one tx, the only one we support
// here we can introduce data.Bytes, which is like []byte
// but serializes to hex in encoding/json
type SetTx struct {
	Key   data.Bytes
	Value data.Bytes
}

// LoadTx handles parsing the binary format, here is a simple intro to go-wire
func LoadTx(txBytes []byte) (SetTx, error) {
	var tx SetTx
	err := data.FromWire(txBytes, &tx)
	return tx, err
}

// Bytes encodes this, must be symetric with LoadTx
// Shows how to encode for the client side
func (s SetTx) Bytes() []byte {
	// note that this type and the recevier must differ by exactly one points
	// eg. if you serialize SetTx, you pass *SetTx to the loader
	// if you serialize *SetTx, you must pass **SetTx to the loader
	return wire.BinaryBytes(s)
}

// Show the minimal way to handle delivertx
func (a *App) DeliverTx(txBytes []byte) abci.Result {
	tx, err := LoadTx(txBytes)
	if err != nil {
		return abci.NewError(err.Error())
	}
	db := a.IAVLApp.DeliverState()
	db.Set(s.Key, s.Value)
	return abci.Result{}
}

// Show the minimal way to handle checktx
func (a *App) CheckTx(txBytes []byte) abci.Result {
	tx, err := LoadTx(txBytes)
	if err != nil {
		return abci.NewError(err.Error())
	}
	db := a.IAVLApp.CheckState()
	db.Set(s.Key, s.Value)
	return abci.Result{}
}
