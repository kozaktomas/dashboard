package cmd

import "github.com/kozaktomas/dashboard/pkg/config"

type commands struct {
	c *config.Config
}

func New(c *config.Config) *commands {
	return &commands{
		c: c,
	}
}