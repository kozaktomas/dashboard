package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/kozaktomas/dashboard/pkg/gui"
	"github.com/kozaktomas/dashboard/pkg/integrations"
	"github.com/kozaktomas/dashboard/pkg/integrations/gitlab"
	"os"
)

func main() {
	_ = godotenv.Load(".env")
	config, _ := loadConfig()

	ins := []integrations.Integration{
		gitlab.New(config),
	}

	g := gui.New(ins)
	g.Run()

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
