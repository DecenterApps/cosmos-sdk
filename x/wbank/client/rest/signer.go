package rest

import (
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/utils"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/wbank/client"

	"github.com/gorilla/mux"
)

type addSignerReq struct {
	BaseReq utils.BaseReq `json:"base_req"`
}

// AddSignerRequestHandlerFn - http request handler for adding new signers.
func AddSignerRequestHandlerFn(cdc *codec.Codec, kb keys.Keybase, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		bech32Addr := vars["address"]

		signer, err := sdk.AccAddressFromBech32(bech32Addr)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		var req addSignerReq
		err = utils.ReadRESTReq(w, r, cdc, &req)
		if err != nil {
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		info, err := kb.Get(baseReq.Name)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := client.CreateAddSignerMsg(sdk.AccAddress(info.GetPubKey().Address()), signer)
		utils.CompleteAndBroadcastTxREST(w, r, cliCtx, baseReq, []sdk.Msg{msg}, cdc)
	}
}
