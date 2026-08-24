package model

import (
	"errors"
)

var InvalidCommandError = errors.New("Invalid command")

type Command struct {
	name   string
	params []string
}

func NewCommand(n string, p []string) *Command {
	return &Command{
		name:   n,
		params: p,
	}
}

func (c *Command) Name() string {
	return c.name
}

func (c *Command) Params() []string {
	return c.params
}
