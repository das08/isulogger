/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Prompt struct {
	promptMsg string
	errorMsg  string
}

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configuration for isulogger.",
	Long:  `Configuration for isulogger.`,
	Run: func(cmd *cobra.Command, args []string) {
		configuration()
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}

func saveConfiguration() {
	viper.WriteConfig()
}

func configuration() {
	isuloggerAPIPrompt := Prompt{
		promptMsg: "Enter isulogger API URL: ",
		errorMsg:  "has to be valid URL",
	}
	isuloggerAPI := PromptGetURL(isuloggerAPIPrompt)

	secretKeyPrompt := Prompt{
		promptMsg: "Enter secret key: ",
		errorMsg:  "has to be valid secret key",
	}
	secretKey := PromptGetString(secretKeyPrompt)

	contestIDPrompt := Prompt{
		"Default Contest ID: ",
		"Contest ID must be integer and greater than 0",
	}
	contestID := PromptGetContestID(contestIDPrompt)

	accessLogPrompt := Prompt{
		"Access Log Path: ",
		"Access Log must be valid path",
	}
	accessLog := PromptGetString(accessLogPrompt)

	slowLogPrompt := Prompt{
		"Slow Log Path: ",
		"Slow Log must be valid path",
	}
	slowLog := PromptGetString(slowLogPrompt)

	viper.Set("isulogger_api", isuloggerAPI)
	viper.Set("secret_key", secretKey)
	viper.Set("contest_id", contestID)
	viper.Set("access_log_path", accessLog)
	viper.Set("slow_log_path", slowLog)

	saveConfiguration()
	printSuccess(fmt.Sprintf("Configuration saved: %s\n", viper.ConfigFileUsed()))
}
