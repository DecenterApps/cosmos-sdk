package wibc

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	DefaultCodespace sdk.CodespaceType = 2

	CodeInvalidDepositLockTime sdk.CodeType = 101
	CodeDepositLockStillActive sdk.CodeType = 102
	CodeInvalidSequence        sdk.CodeType = 103
	CodeIdenticalChains        sdk.CodeType = 104
)

// NOTE: Don't stringer this, we'll put better messages in later.
func codeToDefaultMsg(code sdk.CodeType) string {
	switch code {
	case CodeInvalidDepositLockTime:
		return "invalid deposit lock time"
	case CodeDepositLockStillActive:
		return "deposit lock time is still active"
	case CodeInvalidSequence:
		return "invalid IBC packet sequence"
	case CodeIdenticalChains:
		return "source and destination chain cannot be identical"
	default:
		return sdk.CodeToDefaultMsg(code)
	}
}

//----------------------------------------
// Error constructors

func ErrInvalidDepositLockTime(codespace sdk.CodespaceType) sdk.Error {
	return newError(codespace, CodeInvalidDepositLockTime, "")
}
func ErrDepositLockStillActive(codespace sdk.CodespaceType) sdk.Error {
	return newError(codespace, CodeDepositLockStillActive, "")
}
func ErrInvalidSequence(codespace sdk.CodespaceType) sdk.Error {
	return newError(codespace, CodeInvalidSequence, "")
}
func ErrIdenticalChains(codespace sdk.CodespaceType) sdk.Error {
	return newError(codespace, CodeIdenticalChains, "")
}

//----------------------------------------

func msgOrDefaultMsg(msg string, code sdk.CodeType) string {
	if msg != "" {
		return msg
	}
	return codeToDefaultMsg(code)
}

func newError(codespace sdk.CodespaceType, code sdk.CodeType, msg string) sdk.Error {
	msg = msgOrDefaultMsg(msg, code)
	return sdk.NewError(codespace, code, msg)
}
