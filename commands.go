package main

import (
	"errors"

	"github.com/kobayashicomputing/gator/internal/database"
)

type command struct {
	name string
	args []string
}

type commands struct {
	handler	map[string]func(*state, command) error
}

type middleware struct {
	handler	map[string]func(*state, command, database.User) error
}



func (c *commands) register(name string, f func(*state, command) error) {
	c.handler[name] = f
}

func (c *commands) run(s *state, cmd command) error {
	if op, ok := c.handler[cmd.name]; ok {
		return op(s, cmd)
	} else {
		return errors.New("Command '" + cmd.name + "' does not exist...")
	}
}