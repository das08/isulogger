/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"errors"
	"fmt"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"strconv"
)

var (
	isuloggerAPI  string
	contestID     int
	accessLogPath string
	slowLogPath   string
)

// upCmd represents the up command
var upCmd = &cobra.Command{
	Use:   "up",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := viper.ReadInConfig(); err != nil {
			fmt.Println("Config file not found. Run ./isulogger config to create one.")
			os.Exit(1)
		}

		// Set the API endpoint
		if viper.Get("isulogger_api") != nil {
			isuloggerAPI = viper.Get("isulogger_api").(string)
		} else {
			fmt.Println("isulogger_api not found in config file. Run ./isulogger config to create one.")
			os.Exit(1)
		}

		// Set the contest ID
		if viper.Get("contest_id") != nil {
			contestID = viper.Get("contest_id").(int)
		}
		if id, err := cmd.Flags().GetInt("contestid"); err == nil && id > 0 {
			contestID = id
		}

		if contestID == 0 {
			fmt.Println("Contest id is not set or is invalid. Run ./isulogger config to set Contest ID.")
			os.Exit(1)
		}

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

		fmt.Println("contestID", contestID)
		fmt.Println("accessLogPath", accessLogPath)
		fmt.Println("slowLogPath", slowLogPath)
	},
}

func init() {
	rootCmd.AddCommand(upCmd)

	upCmd.Flags().IntP("contestid", "c", 0, "Contest ID")
}

func promptGetScore(p Prompt) int {
	validate := func(input string) error {
		id, err := strconv.Atoi(input)
		if err != nil {
			return errors.New("has to be integer")
		}
		if id < 0 {
			return errors.New("has to be positive integer")
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

func newEntry() {
	scorePrompt := Prompt{
		promptMsg: "Enter score: ",
		errorMsg:  "Score has to be greater than 0",
	}
	score := promptGetScore(scorePrompt)
	fmt.Println("score", score)
}
