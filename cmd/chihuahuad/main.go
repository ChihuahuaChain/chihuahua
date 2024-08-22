package main

import (
	"fmt"
	"os"

	"github.com/ChihuahuaChain/chihuahua/app"
	"github.com/ChihuahuaChain/chihuahua/cmd/chihuahuad/cmd"
	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
)

func main() {
	rootCmd, _ := cmd.NewRootCmd()

	if err := svrcmd.Execute(rootCmd, "", app.DefaultNodeHome); err != nil {
		fmt.Fprintln(rootCmd.OutOrStderr(), err)
		os.Exit(1)
	}
}
