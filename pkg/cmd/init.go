package cmd

import (
	"gopkg.in/alecthomas/kingpin.v2"
	"strings"
)

func (c *commands) Init(context *kingpin.ParseContext) error {
	gitlabEnabled := askBoolQuestion("Do you want to setup Gitlab configuration? y/n")
	if gitlabEnabled {
		gitlabToken := ask("Gitlab token: ")
		gitlabUserId := askIntQuestion("Gitlab user ID: ")
		gitlabProjectIds := ask("Gitlab project IDs (comma separated): ")

		c.c.Data.Gitlab.Token = gitlabToken
		c.c.Data.Gitlab.UserId = gitlabUserId
		c.c.Data.Gitlab.Projects = strings.Split(gitlabProjectIds, ",")

		err := c.c.Write()
		if err != nil {
			return err
		}
	}

	return nil
}
