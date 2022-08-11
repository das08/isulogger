/*
Copyright © 2022 快適PandA. All rights reserved.

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "See version of this tool",
	Long:  `See version of this tool.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("version 0.0.1")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
