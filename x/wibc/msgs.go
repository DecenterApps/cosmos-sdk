package wibc

import (
	"encoding/json"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Packet struct {
	SrcAddr   sdk.AccAddress `json:"src_addr"`
	DestAddr  sdk.AccAddress `json:"dest_addr"`
	Coins     sdk.Coins      `json:"coins"`
	SrcChain  string         `json:"src_chain"`
	DestChain string         `json:"dest_chain"`
}

func NewIBCPacket(srcAddr sdk.AccAddress, destAddr sdk.AccAddress, coins sdk.Coins,
	srcChain string, destChain string) Packet {

	return Packet{
		SrcAddr:   srcAddr,
		DestAddr:  destAddr,
		Coins:     coins,
		SrcChain:  srcChain,
		DestChain: destChain,
	}
}

//nolint
func (p Packet) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(p)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// validator the ibc packey
func (p Packet) ValidateBasic() sdk.Error {
	if p.SrcChain == p.DestChain {
		return ErrIdenticalChains(DefaultCodespace).TraceSDK("")
	}
	if !p.Coins.IsValid() {
		return sdk.ErrInvalidCoins("")
	}
	return nil
}

type DepositMsg struct {
	From        sdk.AccAddress
	DepositTime int64
}

func NewMsgDeposit(addr sdk.AccAddress, depositTime int64) DepositMsg {
	return DepositMsg{
		From:        addr,
		DepositTime: depositTime,
	}
}

// Implements Msg.
func (msg DepositMsg) Route() string { return "wrapper" }

// Implements Msg.
func (msg DepositMsg) Type() string { return "deposit" }

// Implements Msg.
func (msg DepositMsg) ValidateBasic() sdk.Error {
	if CheckDepositLockTime(msg.From, msg.DepositTime) {
		return ErrInvalidDepositLockTime(DefaultCodespace)
	}

	return nil
}

// Implements Msg.
func (msg DepositMsg) GetSignBytes() []byte {
	bin, err := msgCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(bin)
}

// Implements Msg.
func (msg DepositMsg) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.From}
}

type WithdrawMsg struct {
	Packet Packet `json:"package"` // Deposit Packet
}

var _ sdk.Msg = WithdrawMsg{}

func NewMsgWithdraw(srcAddr sdk.AccAddress, destAddr sdk.AccAddress, coins sdk.Coins,
	srcChain string, destChain string, forTime uint) WithdrawMsg {
	return WithdrawMsg{
		Packet: Packet{
			SrcAddr:   srcAddr,
			DestAddr:  destAddr,
			Coins:     coins,
			SrcChain:  srcChain,
			DestChain: destChain,
		},
	}
}

// Implements Msg.
func (msg WithdrawMsg) Route() string { return "wrapper" }
func (msg WithdrawMsg) Type() string  { return "deposit" }

// Implements Msg.
func (msg WithdrawMsg) ValidateBasic() sdk.Error {
	if len(msg.Packet.SrcAddr) == 0 {
		return sdk.ErrInvalidAddress(msg.Packet.SrcAddr.String())
	}
	if !msg.Packet.Coins.IsValid() {
		return sdk.ErrInvalidCoins(msg.Packet.Coins.String())
	}
	if !msg.Packet.Coins.IsPositive() {
		return sdk.ErrInvalidCoins(msg.Packet.Coins.String())
	}

	return nil
}

// Implements Msg.
func (msg WithdrawMsg) GetSignBytes() []byte {
	bin, err := msgCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(bin)
}

// Implements Msg.
func (msg WithdrawMsg) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Packet.SrcAddr}
}

// ----------------------------------
// ReceiveMsg

// ReceiveMsg defines the message that a relayer uses to post an Packet
// to the destination chain.
type ReceiveMsg struct {
	Packet
	Relayer  sdk.AccAddress
	Sequence int64
}

// nolint
func (msg ReceiveMsg) Route() string            { return "ibc" }
func (msg ReceiveMsg) Type() string             { return "receive" }
func (msg ReceiveMsg) ValidateBasic() sdk.Error { return msg.Packet.ValidateBasic() }

// x/bank/tx.go MsgSend.GetSigners()
func (msg ReceiveMsg) GetSigners() []sdk.AccAddress { return []sdk.AccAddress{msg.Relayer} }

// get the sign bytes for ibc receive message
func (msg ReceiveMsg) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(struct {
		IBCPacket json.RawMessage
		Relayer   sdk.AccAddress
		Sequence  int64
	}{
		IBCPacket: json.RawMessage(msg.Packet.GetSignBytes()),
		Relayer:   msg.Relayer,
		Sequence:  msg.Sequence,
	})
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}
