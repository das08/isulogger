/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
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
		if id, err := cmd.Flags().GetInt("contestid"); err == nil {
			contestID = id
		}

		if contestID == 0 {
			fmt.Println("Contest id is not set or is invalid. Run ./isulogger config to set Contest ID.")
			os.Exit(1)
		}

		// Set the access log path
		if path, err := cmd.Flags().GetString("accesslog"); err == nil {
			accessLogPath = path
		}

		// Set the slow log path
		if path, err := cmd.Flags().GetString("slowlog"); err == nil {
			slowLogPath = path
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
	upCmd.Flags().StringP("accesslog", "a", "", "accesslog")
	upCmd.Flags().StringP("slowlog", "s", "", "slowlog")
}
