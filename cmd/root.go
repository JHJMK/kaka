package cmd

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:   "kaka",
	Short: "kaka is a kafka deploy tool",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func Execute() {
	rootCmd.Execute()
}
