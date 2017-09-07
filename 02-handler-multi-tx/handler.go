package two

import (
	wire "github.com/tendermint/go-wire"
	"github.com/tendermint/go-wire/data"

	sdk "github.com/cosmos/cosmos-sdk"
	"github.com/cosmos/cosmos-sdk/errors"
	"github.com/cosmos/cosmos-sdk/state"
)

// Handler has no state for now, a more complex app could store state here
type Handler struct{}

// NewHandler constructs a Handler and enforces the interface
func NewHandler() sdk.Handler {
	return Handler{}
}

// SetTx sets the data at the key
type SetTx struct {
	Key   data.Bytes
	Value data.Bytes
}

// use var _ to force the compiler to check that the implement the interface
var _ Tx = SetTx{}

// Bytes encodes this, must be symetric with LoadTx
func (s SetTx) Bytes() []byte {
	// note: we must also encode it wrapped in order to get the type byte properly
	return wire.BinaryBytes(TxWrapper{s})
}

// ValidateBasic checks internal validity of SetTx regardless of any state
// Reject any impossible Tx here
func (s SetTx) ValidateBasic() error {
	if len(s.Key) == 0 || len(s.Value) == 0 {
		return errors.ErrInternal("must have non-empty key and value")
	}
	return nil
}

// RemoveTx deletes the contents at a key
type RemoveTx struct {
	Key data.Bytes
}

var _ Tx = RemoveTx{}

// Bytes encodes this, must be symetric with LoadTx
func (r RemoveTx) Bytes() []byte {
	return wire.BinaryBytes(TxWrapper{r})
}

// ValidateBasic checks internal validity of RemoveTx regardless of any state
// Reject any impossible Tx here
func (r RemoveTx) ValidateBasic() error {
	if len(r.Key) == 0 {
		return errors.ErrInternal("must have non-empty key")
	}
	return nil
}

// LoadTx handles parsing the binary format, here is a simple intro to go-wire
func LoadTx(txBytes []byte) (Tx, error) {
	var tx TxWrapper
	err := data.FromWire(txBytes, &tx)
	return tx.Tx, err
}

// DeliverTx applies the tx
func (h Handler) DeliverTx(ctx sdk.Context, store state.SimpleDB,
	msg interface{}) (res sdk.DeliverResult, err error) {

	// here we switch on which implementation of tx we use,
	// and then take the appropriate action.
	switch tx := tx.(type) {
	case SendTx:
		err = tx.ValidateBasic()
		if err != nil {
			break
		}
		store.Set(tx.Key, tx.Value)
		res.Data = tx.Key
	case RemoveTx:
		err = tx.ValidateBasic()
		if err != nil {
			break
		}
		store.Remove(tx.Key)
		res.Data = tx.Key
	default:
		err = errors.ErrInvalidFormat(TxWrapper{}, msg)
	}

	return
}

// CheckTx verifies if it is legit and returns info on how
// to prioritize it in the mempool
func (h Handler) CheckTx(ctx sdk.Context, store state.SimpleDB,
	msg interface{}) (res sdk.CheckResult, err error) {

	// make sure it is something valid
	tx, ok := msg.(Tx)
	if !ok {
		return res, errors.ErrInvalidFormat(TxWrapper{}, msg)
	}
	err = tx.ValidateBasic()
	if err != nil {
		return
	}

	// now return the costs (these should have meaning in your app)
	return sdk.CheckResult{
		GasAllocated: 50,
		GasPayment:   10,
	}
}
