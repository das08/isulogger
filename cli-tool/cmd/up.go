/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
)

var (
	isuloggerAPI  string
	contestID     int
	accessLogPath string
	slowLogPath   string
	score         int
	message       string
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
		getScoreMessage()

		// Check if score and message are set
		if score == 0 || message == "" {
			//fmt.Println("Score and message are not set. Run ./isulogger config to set them.")
			//os.Exit(1)
			postAccessLog()
		} else {
			postScoreMessage()
		}
	},
}

func init() {
	rootCmd.AddCommand(upCmd)

	upCmd.Flags().IntP("contestid", "c", 0, "Contest ID")
}

func promptGetScore(p Prompt) int {
	validate := func(input string) error {
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

func promptGetMessage(p Prompt) string {
	validate := func(input string) error {
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

func getScoreMessage() {
	scorePrompt := Prompt{
		promptMsg: "Enter score: ",
		errorMsg:  "Score has to be greater than 0",
	}
	_score := promptGetScore(scorePrompt)

	messagePrompt := Prompt{
		promptMsg: "Enter message: ",
		errorMsg:  "Message can't be empty",
	}
	_message := promptGetMessage(messagePrompt)

	score = _score
	message = _message
}

func postScoreMessage() {
	type postJson struct {
		ContestID int    `json:"contest_id"`
		Score     int    `json:"score"`
		Message   string `json:"message"`
	}
	postData := postJson{
		ContestID: contestID,
		Score:     score,
		Message:   message,
	}
	postDataJson, _ := json.Marshal(postData)

	endpoint := isuloggerAPI + "/entry"
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(postDataJson))
	if err != nil {
		panic("Error")
	}

	req.Header.Set("Content-Type", "application/json")

	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		panic("Error")
	}
	defer resp.Body.Close()

	byteArray, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic("Error")
	}

	if resp.StatusCode == 200 {
		fmt.Println("[OK] Score and message posted successfully.")
	} else {
		fmt.Println("[ERROR] Score and message posting failed.")
		fmt.Println(string(byteArray), resp.Status)
		os.Exit(1)
	}
}

func postAccessLog() {
	file, err := os.Open(accessLogPath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("log", "name.txt")
	if err != nil {
		panic("Error")
	}
	if _, err := io.Copy(part, file); err != nil {
		panic("Error")
	}
	err = writer.WriteField("token", "token")
	if err != nil {
		panic("Error")
	}
	err = writer.Close()
	if err != nil {
		panic("Error")
	}
	endpoint := isuloggerAPI + "/entry/4/access"
	req, err := http.NewRequest("POST", endpoint, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
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
	fmt.Println("[OK] Access log posted successfully.", resp.Status)
}
