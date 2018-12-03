package wibc

import (
	"github.com/cosmos/cosmos-sdk/types"
	"time"
)

var deposit map[string]int64

func SetNewDepositLock(addr types.AccAddress, time int64) error {
	if CheckDepositLockTime(addr, time) {
		return ErrInvalidDepositLockTime(DefaultCodespace)
	}

	deposit[addr.String()] = time
	return nil
}

func CheckDepositLockTime(addr types.AccAddress, time int64) bool {
	return deposit[addr.String()] < time
}

func HasDepositLockPassed(addr types.AccAddress) bool {
	return deposit[addr.String()] > time.Now().Unix()
}
