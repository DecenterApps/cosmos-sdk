package cli

import (
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/utils"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	authtxb "github.com/cosmos/cosmos-sdk/x/auth/client/txbuilder"
	"github.com/cosmos/cosmos-sdk/x/wibc"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"strconv"
)

const (
	flagTime = "time"
)

// Set deposit lock time
func DepositLockCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use: "transfer",
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := authtxb.NewTxBuilderFromCLI().WithCodec(cdc)
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithAccountDecoder(authcmd.GetAccountDecoder(cdc))

			from, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}

			timeStr := viper.GetString(flagTime)
			time, err := strconv.ParseInt(timeStr, 10, 64)

			msg := wibc.NewMsgDeposit(from, time)

			return utils.CompleteAndBroadcastTxCli(txBldr, cliCtx, []sdk.Msg{msg})
		},
	}

	return cmd
}
