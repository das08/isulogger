package cmd

import (
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
	"net/url"
	"os"
	"strconv"
)

func PromptGetScore(p Prompt) int {
	validate := func(input string) error {
		if input == "0" {
			return nil
		}
		// check if it is a number
		_, err := strconv.Atoi(input)
		if err != nil && input != "" {
			return errors.New(p.errorMsg)
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
		Default:   "",
		Templates: templates,
		Validate:  validate,
		Stdout:    &BellSkipper{},
	}

	result, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		os.Exit(1)
	}

	score, _ := strconv.Atoi(result)
	if result == "" {
		score = -1
	}

	return score
}

func PromptGetString(p Prompt) string {
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
		Stdout:    &BellSkipper{},
	}

	result, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		os.Exit(1)
	}

	return result
}

func PromptGetContestID(p Prompt) int {
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
		Stdout:    &BellSkipper{},
	}

	result, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		os.Exit(1)
	}

	contestID, _ := strconv.Atoi(result)

	return contestID
}

func PromptGetURL(p Prompt) string {
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
		Stdout:    &BellSkipper{},
	}

	result, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		os.Exit(1)
	}

	return result
}

func PromptGetYN(p Prompt) bool {
	validate := func(input string) error {
		if input != "y" && input != "n" && input != "Y" && input != "N" && input != "" {
			return errors.New(p.errorMsg)
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
		Default:   "y",
		Templates: templates,
		Validate:  validate,
		Stdout:    &BellSkipper{},
	}

	result, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		os.Exit(1)
	}
	if result == "y" || result == "Y" || result == "" {
		return true
	}
	return false
}

func printSuccess(successMsg string) {
	color.Set(color.FgGreen, color.Bold)
	fmt.Printf("[OK] ")
	color.Unset()
	fmt.Println(successMsg)
}

func printWarning(warningMsg string) {
	color.Set(color.FgHiYellow, color.Bold)
	fmt.Printf("[--] ")
	color.Unset()
	fmt.Println(warningMsg)
}

func printError(errorMsg string) {
	color.Set(color.FgRed, color.Bold)
	fmt.Printf("[Error] ")
	color.Unset()
	fmt.Println(errorMsg)
}
