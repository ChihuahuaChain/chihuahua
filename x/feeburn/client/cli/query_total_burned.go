package cli

import (
	"context"

	"github.com/ChihuahuaChain/chihuahua/x/feeburn/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
)

func CmdQueryTotalBurned() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "total-burned",
		Aliases: []string{"burned", "burn"},
		Short:   "shows the transaction fees burned so far",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.TotalBurned(context.Background(), &types.QueryTotalBurnedRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
