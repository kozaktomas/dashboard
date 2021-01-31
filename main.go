package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/kozaktomas/dashboard/pkg/gitlab"
	gui2 "github.com/kozaktomas/dashboard/pkg/gui"
	"os"
	"strconv"
	"strings"
)

func main() {
	_ = godotenv.Load(".env")
	config, _ := loadConfig()

	gitlabUserId, err := strconv.Atoi(config["GITLAB_USER_ID"])
	if err != nil {
		panic("Gitlab user ID need to be integer")
	}
	gl, _ := gitlab.New(
		config["GITLAB_TOKEN"],
		gitlabUserId,
		strings.Split(config["GITLAB_PROJECTS"], ","),
	)

	gui := gui2.New(gl)
	gui.Run()

}


func loadConfig() (map[string]string, error) {
	keys := []string{
		"GITLAB_USER_ID",
		"GITLAB_TOKEN",
		"GITLAB_PROJECTS",
	}

	config := make(map[string]string)
	for _, key := range keys {
		v := os.Getenv(key)
		if v == "" {
			return nil, fmt.Errorf("could not find environment variable %s", key)
		}
		config[key] = v
	}

	return config, nil
}
