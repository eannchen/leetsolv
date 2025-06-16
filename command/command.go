package command

import (
	"bufio"
	"fmt"

	"leetsolv/handler"
)

type Command interface {
	Execute(scanner *bufio.Scanner, args []string) bool
}

type CommandRegistry struct {
	commands map[string]Command
}

func NewCommandRegistry() *CommandRegistry {
	return &CommandRegistry{
		commands: make(map[string]Command),
	}
}

func (r *CommandRegistry) Register(name string, cmd Command) {
	r.commands[name] = cmd
}

func (r *CommandRegistry) Execute(scanner *bufio.Scanner, name string, args []string) bool {
	if cmd, exists := r.commands[name]; exists {
		if quit := cmd.Execute(scanner, args); quit {
			return true
		}
	} else {
		fmt.Println("Unknown command.")
	}
	return false
}

// command implementations
// ListCommand, StatusCommand, UpsertCommand, DeleteCommand, UndoCommand, QuitCommand
// are defined below, each implementing the Command interface.

type ListCommand struct {
	Handler handler.Handler
}

func (c *ListCommand) Execute(scanner *bufio.Scanner, args []string) bool {
	c.Handler.HandleList(scanner)
	return false
}

type GetCommand struct {
	Handler handler.Handler
}

func (c *GetCommand) Execute(scanner *bufio.Scanner, args []string) bool {
	var target string
	if len(args) > 0 {
		target = args[0]
	}
	c.Handler.HandleGet(scanner, target)
	return false
}

type StatusCommand struct {
	Handler handler.Handler
}

func (c *StatusCommand) Execute(scanner *bufio.Scanner, args []string) bool {
	c.Handler.HandleStatus()
	return false
}

type UpsertCommand struct {
	Handler handler.Handler
}

func (c *UpsertCommand) Execute(scanner *bufio.Scanner, args []string) bool {
	c.Handler.HandleUpsert(scanner)
	return false
}

type DeleteCommand struct {
	Handler handler.Handler
}

func (c *DeleteCommand) Execute(scanner *bufio.Scanner, args []string) bool {
	var target string
	if len(args) > 0 {
		target = args[0]
	}
	c.Handler.HandleDelete(scanner, target)
	return false
}

type UndoCommand struct {
	Handler handler.Handler
}

func (c *UndoCommand) Execute(scanner *bufio.Scanner, args []string) bool {
	c.Handler.HandleUndo(scanner)
	return false
}

type QuitCommand struct{}

func (c *QuitCommand) Execute(scanner *bufio.Scanner, args []string) bool {
	return true
}
