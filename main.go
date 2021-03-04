package main

import (
	"fmt"
	"github.com/kozaktomas/dashboard/pkg/cmd"
	"github.com/kozaktomas/dashboard/pkg/config"
	"github.com/kozaktomas/dashboard/pkg/gui"
	"github.com/kozaktomas/dashboard/pkg/integrations"
	"github.com/kozaktomas/dashboard/pkg/integrations/gitlab"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
	"os/user"
)

var app = kingpin.New("ddboard", "Developer dashboard.")
var c *config.Config // global configuration

func main() {
	usr, err := user.Current()
	if err != nil {
		fmt.Printf("could not get current user: %w", err)
		return
	}

	c = config.New(usr.HomeDir)
	if len(os.Args) == 1 && c.IsReady() {
		// run UI app
		ins := []integrations.Integration{
			gitlab.New(c.Data.Gitlab.Url, c.Data.Gitlab.Token, c.Data.Gitlab.UserId, c.Data.Gitlab.Projects),
		}

		g := gui.New(ins)
		g.Run()
	} else {
		// run kingpin app
		commands := cmd.New(c)
		app.Command("init", "Init configuration").Action(commands.Init)
		command, err := app.Parse(os.Args[1:])
		kingpin.MustParse(command, err)
	}
}
