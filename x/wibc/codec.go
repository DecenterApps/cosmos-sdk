package wibc

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// Register concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(ReceiveMsg{}, "cosmos-sdk/WithdrawMsg", nil)
	cdc.RegisterConcrete(DepositMsg{}, "cosmos-sdk/DepositMsg", nil)
	cdc.RegisterConcrete(WithdrawMsg{}, "cosmos-sdk/Send", nil)
}

var msgCdc = codec.New()

func init() {
	RegisterCodec(msgCdc)
}
