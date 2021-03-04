package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func ask(question, defaultValue string) string {
	reader := bufio.NewReader(os.Stdin)

	askQuestion := question
	if defaultValue != "" {
		askQuestion += " [default: " + defaultValue + "]"
	}
	askQuestion += ": "
	fmt.Print(askQuestion)

	text, _ := reader.ReadString('\n')
	text = strings.Trim(text, "\n")
	text = strings.Trim(text, " ")
	fmt.Println()

	if text == "" {
		text = defaultValue
	}

	return text
}

func askBoolQuestion(question string) bool {
	answer := strings.ToLower(ask(question, ""))
	yeses := []string{"y", "ye", "yes"}
	for _, yes := range yeses {
		if answer == yes {
			return true
		}
	}

	return false
}

func askIntQuestion(question string) int {
	answer := ask(question, "")
	i, err := strconv.Atoi(answer)
	if err != nil {
		return -1
	}

	return i
}
