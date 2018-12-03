package rest

import (
	"github.com/cosmos/cosmos-sdk/x/wibc"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/utils"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type depositReq struct {
	BaseReq utils.BaseReq `json:"base_req"`
	Time    int64         `json:"time"`
}

// TransferRequestHandler - http request handler to transfer coins to a address
// on a different chain via IBC.
func DepostiRequestHandlerFn(cdc *codec.Codec, kb keys.Keybase, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req depositReq
		err := utils.ReadRESTReq(w, r, cdc, &req)
		if err != nil {
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		info, err := kb.Get(baseReq.Name)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}

		msg := wibc.NewMsgDeposit(sdk.AccAddress(info.GetPubKey().Address()), req.Time)
		utils.CompleteAndBroadcastTxREST(w, r, cliCtx, baseReq, []sdk.Msg{msg}, cdc)
	}
}
