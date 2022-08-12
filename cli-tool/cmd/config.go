/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"errors"
	"fmt"
	"github.com/manifoldco/promptui"
	"github.com/spf13/viper"
	"net/url"
	"os"
	"strconv"

	"github.com/spf13/cobra"
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
		fmt.Println("config called")
		if err := viper.ReadInConfig(); err == nil {
			fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
		}

		//if viper.Get("config") != nil {
		//	fmt.Println(viper.Get("config"))
		//} else {
		//	fmt.Println("config is nil")
		//}
		configuration()
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}

func promptGetInputInt(p Prompt) int {
	validate := func(input string) error {
		if len(input) <= 0 {
			return errors.New(p.errorMsg)
		}

		id, err := strconv.Atoi(input)
		if err != nil {
			return errors.New("has to be integer")
		}
		if id <= 0 {
			return errors.New("has to be greater than 0")
		}
		return nil
	}

	templates := &promptui.PromptTemplates{
		Prompt:  "{{ . }} ",
		Valid:   "{{ . | green }} ",
		Invalid: "{{ . | red }} ",
		Success: "{{ . | bold }} ",
	}

	prompt := promptui.Prompt{
		Label:     p.promptMsg,
		Templates: templates,
		Validate:  validate,
	}

	result, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		os.Exit(1)
	}

	contestID, _ := strconv.Atoi(result)

	return contestID
}

func promptGetInputStr(p Prompt) string {
	validate := func(input string) error {
		if len(input) <= 0 {
			return errors.New(p.errorMsg)
		}
		_, err := url.ParseRequestURI(input)
		if err != nil {
			return errors.New("invalid url")
		}
		return nil
	}

	templates := &promptui.PromptTemplates{
		Prompt:  "{{ . }} ",
		Valid:   "{{ . | green }} ",
		Invalid: "{{ . | red }} ",
		Success: "{{ . | bold }} ",
	}

	prompt := promptui.Prompt{
		Label:     p.promptMsg,
		Templates: templates,
		Validate:  validate,
	}

	result, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		os.Exit(1)
	}

	return result
}

func saveConfiguration() {
	viper.WriteConfig()
}

func configuration() {
	isuloggerAPIPrompt := Prompt{
		promptMsg: "Enter isulogger API URL: ",
		errorMsg:  "has to be valid URL",
	}
	isuloggerAPI := promptGetInputStr(isuloggerAPIPrompt)

	contestIDPrompt := Prompt{
		"Default Contest ID: ",
		"Contest ID must be integer and greater than 0",
	}
	contestID := promptGetInputInt(contestIDPrompt)

	viper.Set("isulogger_api", isuloggerAPI)
	viper.Set("contest_id", contestID)

	saveConfiguration()
}
