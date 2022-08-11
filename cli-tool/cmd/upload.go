/*
Copyright © 2022 快適PandA. All rights reserved.

*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

type Upload struct {
	ContestID     int    `json:"contest_id"`
	Score         int    `json:"score"`
	Message       string `json:"message"`
	AccessLogPath string `json:"access_log_path"`
	SlowLogPath   string `json:"slow_log_path"`
}

// uploadCmd represents the upload command
var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Command for uploading logs to ISULOGGER",
	Long:  `Command for uploading logs to ISULOGGER.`,
	Run: func(cmd *cobra.Command, args []string) {
		contestID, _ := cmd.Flags().GetInt("contest_id")
		score, _ := cmd.Flags().GetInt("score")
		message, _ := cmd.Flags().GetString("message")
		accessLogPath, _ := cmd.Flags().GetString("access_log_path")
		slowLogPath, _ := cmd.Flags().GetString("slow_log_path")
		fmt.Println("contestID:", contestID)
		fmt.Println("score:", score)
		fmt.Println("message:", message)
		fmt.Println("accessLogPath:", accessLogPath)
		fmt.Println("slowLogPath:", slowLogPath)
	},
}

func init() {
	rootCmd.AddCommand(uploadCmd)

	uploadCmd.Flags().IntP("contest_id", "c", 0, "Contest ID")
	uploadCmd.Flags().IntP("score", "s", 0, "Score")
	uploadCmd.Flags().StringP("message", "m", "", "Message")
	uploadCmd.Flags().StringP("access_log_path", "a", "", "Access Log Path")
	uploadCmd.Flags().StringP("slow_log_path", "l", "", "Slow Log Path")

}
