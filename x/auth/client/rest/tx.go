package rest

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"

	"github.com/certikfoundation/shentu/v2/x/auth/types"
)

// RegisterRoutes registers custom REST routes.
func RegisterRoutes(cliCtx client.Context, r *mux.Router) {
	r.HandleFunc("/txs/{address}/unlock", UnlockRequestHandlerFn(cliCtx)).Methods("POST")
}

// UnlockReq defines the properties of a unlock request's body.
type UnlockReq struct {
	BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`
	Amount  sdk.Coins    `json:"amount" yaml:"amount"`
}

// UnlockRequestHandlerFn handles a request sent by an unlocker
// to unlock coins from a manual vesting account
func UnlockRequestHandlerFn(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		bech32Addr := vars["address"]
		accAddr, err := sdk.AccAddressFromBech32(bech32Addr)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		var req UnlockReq
		if !rest.ReadRESTReq(w, r, cliCtx.LegacyAmino, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		fromAddr, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		msg := types.NewMsgUnlock(fromAddr, accAddr, req.Amount)
		tx.WriteGeneratedTxResponse(cliCtx, w, req.BaseReq, msg)
	}
}
