/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"github.com/spf13/viper"
	"os"

	"github.com/spf13/cobra"
)

// truncateCmd represents the truncate command
var truncateCmd = &cobra.Command{
	Use:   "truncate",
	Short: "Truncate isulogger logs",
	Long:  `Truncate isulogger logs.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Set the access log path
		if viper.Get("access_log_path") != nil {
			accessLogPath = viper.Get("access_log_path").(string)
		}

		// Set the slow log path
		if viper.Get("slow_log_path") != nil {
			slowLogPath = viper.Get("slow_log_path").(string)
		}

		if accessLogPath == "" && slowLogPath == "" {
			fmt.Println("Access log path and/or slow log path are not set. ")
			os.Exit(1)
		}

		// Set the skip flag
		if s, err := cmd.Flags().GetBool("skip"); err == nil {
			skip = s
		}

		if !skip {
			if !skipTruncateConfirmation() {
				return
			}
		}
		err := os.Truncate(accessLogPath, 0)
		if err != nil {
			printError("Access log truncate failed.")
		}
		printSuccess("Access log truncated.")
		err = os.Truncate(slowLogPath, 0)
		if err != nil {
			printError("Slow log truncate failed.")
		}
		printSuccess("Slow log truncated.")
	},
}

func init() {
	rootCmd.AddCommand(truncateCmd)

	truncateCmd.Flags().BoolP("skip", "s", false, "Skip confirmation")
}

func skipTruncateConfirmation() bool {
	confirmPrompt := Prompt{
		promptMsg: "Are you sure you want to truncate logs? (Y/n): ",
		errorMsg:  "Please enter y or n",
	}
	return PromptGetYN(confirmPrompt)
}
