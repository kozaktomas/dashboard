package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func ask(question string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println(question)
	text, _ := reader.ReadString('\n')
	text = strings.Trim(text, "\n")
	text = strings.Trim(text, " ")
	return text
}

func askBoolQuestion(question string) bool {
	answer := strings.ToLower(ask(question))
	yeses := []string{"y", "ye", "yes"}
	for _, yes := range yeses {
		if answer == yes {
			return true
		}
	}

	return false
}

func askIntQuestion(question string) int {
	answer := ask(question)
	i, err := strconv.Atoi(answer)
	if err != nil {
		return -1
	}

	return i
}
