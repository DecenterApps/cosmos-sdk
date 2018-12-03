package wbank

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Mapping that keeps track of all signers
var isSigner map[string]bool

// Checks if address is a signer
func IsSigner(addr sdk.AccAddress) bool {
	return isSigner[addr.String()]
}

// Adds new signer
func AddSigner(addr sdk.AccAddress) {
	isSigner[addr.String()] = true
}
