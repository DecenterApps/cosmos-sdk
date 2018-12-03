package cli

import (
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/utils"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	authtxb "github.com/cosmos/cosmos-sdk/x/auth/client/txbuilder"
	"github.com/cosmos/cosmos-sdk/x/wbank/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	flagSigner = "signer"
)

func AddSigner(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-signer",
		Short: "Adds a new signer",
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := authtxb.NewTxBuilderFromCLI().WithCodec(cdc)
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithAccountDecoder(authcmd.GetAccountDecoder(cdc))

			if err := cliCtx.EnsureAccountExists(); err != nil {
				return err
			}

			signerStr := viper.GetString(flagSigner)

			signer, err := sdk.AccAddressFromBech32(signerStr)
			if err != nil {
				return err
			}

			from, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}

			msg := client.CreateAddSignerMsg(from, signer)
			if cliCtx.GenerateOnly {
				return utils.PrintUnsignedStdTx(txBldr, cliCtx, []sdk.Msg{msg}, false)
			}

			return utils.CompleteAndBroadcastTxCli(txBldr, cliCtx, []sdk.Msg{msg})
		},
	}

	cmd.Flags().String(flagSigner, "", "Address to add to signers list")
	cmd.MarkFlagRequired(flagSigner)

	return cmd
}
