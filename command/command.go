package command

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"leetsolv/handler"
)

type Command interface {
	Execute(scanner *bufio.Scanner) bool
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

func (r *CommandRegistry) Execute(name string, scanner *bufio.Scanner) bool {
	if cmd, exists := r.commands[name]; exists {
		if quit := cmd.Execute(scanner); quit {
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

func (c *ListCommand) Execute(scanner *bufio.Scanner) bool {
	c.Handler.HandleList(scanner)
	return false
}

type GetCommand struct {
	Handler handler.Handler
}

func (c *GetCommand) Execute(scanner *bufio.Scanner) bool {
	var input string
	if len(os.Args) > 2 { // Check if input is provided in CLI argument mode
		input = strings.Join(os.Args[2:], " ")
	}
	c.Handler.HandleGet(scanner, input)
	return false
}

type StatusCommand struct {
	Handler handler.Handler
}

func (c *StatusCommand) Execute(scanner *bufio.Scanner) bool {
	c.Handler.HandleStatus()
	return false
}

type UpsertCommand struct {
	Handler handler.Handler
}

func (c *UpsertCommand) Execute(scanner *bufio.Scanner) bool {
	c.Handler.HandleUpsert(scanner)
	return false
}

type DeleteCommand struct {
	Handler handler.Handler
}

func (c *DeleteCommand) Execute(scanner *bufio.Scanner) bool {
	c.Handler.HandleDelete(scanner)
	return false
}

type UndoCommand struct {
	Handler handler.Handler
}

func (c *UndoCommand) Execute(scanner *bufio.Scanner) bool {
	c.Handler.HandleUndo(scanner)
	return false
}

type QuitCommand struct{}

func (c *QuitCommand) Execute(scanner *bufio.Scanner) bool {
	return true
}
