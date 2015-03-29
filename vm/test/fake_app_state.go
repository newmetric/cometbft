package vm

import (
	"fmt"

	. "github.com/tendermint/tendermint/common"
	. "github.com/tendermint/tendermint/vm"
	"github.com/tendermint/tendermint/vm/sha3"
)

type FakeAppState struct {
	accounts map[string]*Account
	storage  map[string]Word256
	logs     []*Log
}

func (fas *FakeAppState) GetAccount(addr Word256) *Account {
	account := fas.accounts[addr.String()]
	if account != nil {
		return account
	} else {
		panic(Fmt("Invalid account addr: %X", addr))
	}
}

func (fas *FakeAppState) UpdateAccount(account *Account) {
	_, ok := fas.accounts[account.Address.String()]
	if !ok {
		panic(Fmt("Invalid account addr: %X", account.Address))
	} else {
		// Nothing to do
	}
}

func (fas *FakeAppState) RemoveAccount(account *Account) {
	_, ok := fas.accounts[account.Address.String()]
	if !ok {
		panic(Fmt("Invalid account addr: %X", account.Address))
	} else {
		// Remove account
		delete(fas.accounts, account.Address.String())
	}
}

func (fas *FakeAppState) CreateAccount(creator *Account) *Account {
	addr := createAddress(creator)
	account := fas.accounts[addr.String()]
	if account == nil {
		return &Account{
			Address:     addr,
			Balance:     0,
			Code:        nil,
			Nonce:       0,
			StorageRoot: Zero256,
		}
	} else {
		panic(Fmt("Invalid account addr: %X", addr))
	}
}

func (fas *FakeAppState) GetStorage(addr Word256, key Word256) Word256 {
	_, ok := fas.accounts[addr.String()]
	if !ok {
		panic(Fmt("Invalid account addr: %X", addr))
	}

	value, ok := fas.storage[addr.String()+key.String()]
	if ok {
		return value
	} else {
		return Zero256
	}
}

func (fas *FakeAppState) SetStorage(addr Word256, key Word256, value Word256) {
	_, ok := fas.accounts[addr.String()]
	if !ok {
		panic(Fmt("Invalid account addr: %X", addr))
	}

	fas.storage[addr.String()+key.String()] = value
}

func (fas *FakeAppState) AddLog(log *Log) {
	fas.logs = append(fas.logs, log)
}

func main() {
	appState := &FakeAppState{
		accounts: make(map[string]*Account),
		storage:  make(map[string]Word256),
		logs:     nil,
	}
	params := Params{
		BlockHeight: 0,
		BlockHash:   Zero256,
		BlockTime:   0,
		GasLimit:    0,
	}
	ourVm := NewVM(appState, params, Zero256)

	// Create accounts
	account1 := &Account{
		Address: Uint64ToWord256(100),
	}
	account2 := &Account{
		Address: Uint64ToWord256(101),
	}

	var gas uint64 = 1000
	output, err := ourVm.Call(account1, account2, []byte{0x5B, 0x60, 0x00, 0x56}, []byte{}, 0, &gas)
	fmt.Printf("Output: %v Error: %v\n", output, err)
}

// Creates a 20 byte address and bumps the nonce.
func createAddress(creator *Account) Word256 {
	nonce := creator.Nonce
	creator.Nonce += 1
	temp := make([]byte, 32+8)
	copy(temp, creator.Address[:])
	PutUint64(temp[32:], nonce)
	return RightPadWord256(sha3.Sha3(temp)[:20])
}
