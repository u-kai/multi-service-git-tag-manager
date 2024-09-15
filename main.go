package main

import "github.com/spf13/cobra"

func main() {
	rootCmd := &cobra.Command{
		Use:   "msgtm",
		Short: "msgtm is a tool for multi service git tag manager",
		Run: func(cmd *cobra.Command, args []string) {
			println("msgtm")
		},
	}
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
