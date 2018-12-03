package wibc

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/wbank"
)

func NewHandler(ibcm Mapper, ck wbank.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case WithdrawMsg:
			return handleWithdrawMsg(ctx, ibcm, ck, msg)
		case ReceiveMsg:
			return handleReceiveMsg(ctx, ibcm, ck, msg)
		default:
			errMsg := "Unrecognized IBC Msg type: " + msg.Type()
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// WithdrawMsg deducts coins from the account and creates an egress IBC packet.
func handleWithdrawMsg(ctx sdk.Context, ibcm Mapper, ck wbank.Keeper, msg WithdrawMsg) sdk.Result {
	packet := msg.Packet

	_, _, err := ck.SubtractCoins(ctx, packet.SrcAddr, packet.Coins)
	if err != nil {
		return err.Result()
	}

	err = ibcm.PostPacket(ctx, packet)
	if err != nil {
		return err.Result()
	}

	return sdk.Result{}
}

// ReceiveMsg adds coins to the destination address and creates an ingress IBC packet.
func handleReceiveMsg(ctx sdk.Context, ibcm Mapper, ck wbank.Keeper, msg ReceiveMsg) sdk.Result {
	packet := msg.Packet

	seq := ibcm.GetIngressSequence(ctx, packet.SrcChain)
	if msg.Sequence != seq {
		return ErrInvalidSequence(ibcm.codespace).Result()
	}

	_, _, err := ck.AddCoins(ctx, packet.DestAddr, packet.Coins)
	if err != nil {
		return err.Result()
	}

	ibcm.SetIngressSequence(ctx, packet.SrcChain, seq+1)

	return sdk.Result{}
}
