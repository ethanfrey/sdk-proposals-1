package main

import "fmt"

func main() {
	fmt.Println("vim-go")
}

type MasterMerkleStore interface {
	Commit() []byte // returns a Simple merkle hash

	Set(name string, store MerkleStore)
	Get(name string) MerkleStore
	GetScratch(name string) MerkleStore

	// Convenience getters
	GetKVStore(name string) KVStore
	GetKVStoreScratch(name string) KVStore
	GetHeapStore(name string) Heapstore
	GetHeapStoreScratch(name string) HeapStore
}

// NOTE: not necessarily a KV store.
// .ScratchCopy() returns a copy of the store that will be discarded.
type MerkleStore interface {
	ScratchCopy() MerkleStore // calling Commit on this will panic.
	Commit() (hash []byte)    // may panic if scratch copy

	// Convenience functions like .KVStore() and .HeapStore()?
}

type KVStore interface {
	MerkleStore
	Get(key string) (value string)
	Set(key string, value string)
	Domain(prefix string) KVStore
}

// TODO: just for example
type HeapStore interface {
	MerkleStore
	Add(key string, value string)
	Pop() (key string, value string)
}

//--------------------------------------------------------------------------------

type AccountStore struct {
	KVStore
}

// NOTE: Account is declared in app/types/account.go, and is application-dependent
// This means AccountStore needs to be declared at the application level as well.
// We could generate code into app/types/*_store.go, but it's also easy enough to hand-code.
func (as AccountStore) GetAccount(addr []byte) *Account {
	return as.Get(string(addr)).(*Account) // or whatever
}

func (as AccountStore) UpdateAccount(account *Account) {
	// whatever
}

//--------------------------------------------------------------------------------

func main1() {

	// ABCI init
	var kv1 KVStore
	var kv2 KVStore
	var hs1 HeapStore

	store := NewMasterMerkleStore()
	store.Set("kv1", kv1)
	store.Set("kv2", kv2)
	store.Set("hs1", hs3)
	app.store = store

	// ABCI deliver tx
	db := app.store.GetKVStore("kv1")
	db.Get()
	db.Set()
	// ...

	// ABCI check tx
	db := app.store.GetKVStoreScratch("kv1")
	db.Get()
	db.Set()
	// ...

	// ABCI commit
	hash := app.store.Commit()
	// ...

}

//--------------------------------------------------------------------------------

func main2() {

	// ABCI init
	var kv1 KVStore
	var kv2 KVStore
	var hs1 HeapStore

	store := NewMasterMerkleStore()
	store.Set("kv1", kv1)
	store.Set("kv2", kv2)
	store.Set("hs1", hs3)
	app.store = store

	// ABCI begin block
	app.kv1 = app.store.GetKVStore("kv1")
	app.kv1CheckTx = app.store.GetKVStoreScratch("kv1")
	app.kv2 = app.store.GetKVStore("kv2")
	app.kv2CheckTx = app.store.GetKVStoreScratch("kv2")
	app.hs1 = app.store.GetKVStore("hs1")
	app.hs1CheckTx = app.store.GetKVStoreScratch("hs1")
	app.accountStore = AccountStore{app.kv1}
	app.accountStoreCheckTx = AccountStore{app.kv1}

	// ABCI deliver tx
	accountStore.GetAccount(addr)
	// ...

	// ABCI check tx
	accountStoreCheckTx.GetAccount(addr)
	// ...

	// ABCI commit
	hash := app.store.Commit()
	// ...

}
