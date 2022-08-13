/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

var (
	isuloggerAPI  string
	secretKey     string
	contestID     int
	accessLogPath string
	slowLogPath   string
	score         int
	message       string
	skip          bool
	noScore       bool
)

// uploadCmd represents the up command
var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Upload score, message, and logs to isulogger",
	Long:  `Upload score, message, and logs to isulogger.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := viper.ReadInConfig(); err != nil {
			fmt.Println("Config file not found. Run ./isulogger config to create one.")
			os.Exit(1)
		}

		// Set the API endpoint
		if viper.Get("isulogger_api") != nil {
			isuloggerAPI = viper.Get("isulogger_api").(string)
		} else {
			fmt.Println("isulogger_api not found in config file. Run ./isulogger config .")
			os.Exit(1)
		}

		// Set the secret key
		if viper.Get("secret_key") != nil {
			secretKey = viper.Get("secret_key").(string)
		} else {
			fmt.Println("secret_key not found in config file. Run ./isulogger config .")
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

		// Set the skip flag
		if s, err := cmd.Flags().GetBool("skip"); err == nil {
			skip = s
		}

		// Set the no-score flag
		if x, err := cmd.Flags().GetBool("no-score"); err == nil {
			noScore = x
		}

		//fmt.Println("contestID", contestID)
		//fmt.Println("accessLogPath", accessLogPath)
		//fmt.Println("slowLogPath", slowLogPath)

		contestName := getContestName()

		fmt.Printf("[Contest Name] ")
		color.Set(color.Bold, color.Underline)
		fmt.Printf("%s (ID: %d)\n", contestName, contestID)
		color.Unset()

		if !noScore {
			getScoreMessage()
		}

		// Check if score and message are set
		if score != 0 {
			postScoreMessage()
		} else {
			printWarning("Score is not set. Skipping score and message upload.")
		}

		if !skip {
			if !confirmMessage() {
				return
			}
		}

		// Upload logs
		if accessLogPath != "" {
			postLog("access")
		} else {
			printWarning("Access log path is not set. Skipping access log upload.")
		}
		if slowLogPath != "" {
			postLog("slow")
		} else {
			printWarning("Slow log path is not set. Skipping slow log upload.")
		}
	},
}

func init() {
	rootCmd.AddCommand(uploadCmd)

	uploadCmd.Flags().IntP("contestid", "c", 0, "Contest ID")
	uploadCmd.Flags().BoolP("skip", "s", false, "Skip confirmation prompt")
	uploadCmd.Flags().BoolP("no-score", "x", false, "Skip score prompt")
}

func getScoreMessage() {
	fmt.Println("Enter score and message. Leave score blank to skip creating new log entry.")
	scorePrompt := Prompt{
		promptMsg: "Enter score: ",
		errorMsg:  "Score has to be greater than 0",
	}
	_score := PromptGetScore(scorePrompt)

	messagePrompt := Prompt{
		promptMsg: "Enter message: ",
		errorMsg:  "Message can't be empty",
	}
	_message := PromptGetString(messagePrompt)

	score = _score
	message = _message
}

func getContestName() string {
	contestName := ""
	type Contest struct {
		ContestID   int    `json:"contest_id"`
		ContestName string `json:"contest_name"`
	}
	var contests []Contest

	endpoint := isuloggerAPI + "/contest"
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		panic("Error")
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Secret-Key", secretKey)

	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		panic("Error")
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = json.Unmarshal(body, &contests)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for _, contest := range contests {
		if contest.ContestID == contestID {
			contestName = contest.ContestName
		}
	}

	return contestName
}

func getBranchName() string {
	output, err := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD").Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(output))
}

func postScoreMessage() {
	type postJson struct {
		ContestID  int    `json:"contest_id"`
		BranchName string `json:"branch_name"`
		Score      int    `json:"score"`
		Message    string `json:"message"`
	}
	postData := postJson{
		ContestID:  contestID,
		BranchName: getBranchName(),
		Score:      score,
		Message:    message,
	}
	postDataJson, _ := json.Marshal(postData)

	endpoint := isuloggerAPI + "/entry"
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(postDataJson))
	if err != nil {
		panic("Error")
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Secret-Key", secretKey)

	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		panic("Error")
	}
	defer resp.Body.Close()

	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		panic("Error")
	}

	if resp.StatusCode == 200 {
		printSuccess(fmt.Sprintf("Score and message successfully posted to IsuLogger. %d", resp.StatusCode))
	} else {
		printError(fmt.Sprintf("Score and message posting failed. %d", resp.StatusCode))
	}
}

func confirmMessage() bool {
	confirmPrompt := Prompt{
		promptMsg: "Are you sure you want to upload logs? (Y/n): ",
		errorMsg:  "Please enter y or n",
	}
	return PromptGetYN(confirmPrompt)
}

func postLog(logType string) {
	if logType != "access" && logType != "slow" {
		fmt.Println("Log type is invalid. It has to be access or slow.")
		os.Exit(1)
	}

	var logPath, fileName string
	now := time.Now().Format("2006-01-02-15-04-05")
	if logType == "access" {
		logPath = accessLogPath
		fileName = "access-" + now + ".log"
	} else {
		logPath = slowLogPath
		fileName = "slow-" + now + ".log"
	}

	file, err := os.Open(logPath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("log", fileName)
	if err != nil {
		panic("Error")
	}
	if _, err := io.Copy(part, file); err != nil {
		panic("Error")
	}
	err = writer.Close()
	if err != nil {
		panic("Error")
	}

	endpoint := fmt.Sprintf("%s/entry/%d/%s", isuloggerAPI, contestID, logType)
	req, err := http.NewRequest("POST", endpoint, body)

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("X-Secret-Key", secretKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic("Error")
	}
	defer resp.Body.Close()
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		panic("Error")
	}

	if resp.StatusCode == 200 {
		printSuccess(fmt.Sprintf("%s log posted successfully. %d", logType, resp.StatusCode))
	} else {
		printError(fmt.Sprintf("%s log posting failed. %d", logType, resp.StatusCode))
	}
}
