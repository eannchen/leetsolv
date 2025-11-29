// Package command links between the command line and the handler.
package command

import (
	"bufio"
	"strings"

	"github.com/eannchen/leetsolv/handler"
)

type Command interface {
	Execute(scanner *bufio.Scanner, args []string) bool
}

type CommandRegistry struct {
	commands              map[string]Command
	unknownCommandHandler func(command string)
}

func NewCommandRegistry(unknownCommandHandler func(command string)) *CommandRegistry {
	return &CommandRegistry{
		commands:              make(map[string]Command),
		unknownCommandHandler: unknownCommandHandler,
	}
}

func (r *CommandRegistry) Register(name string, cmd Command) {
	// Convert command name to lowercase for case-insensitive registration
	r.commands[strings.ToLower(name)] = cmd
}

func (r *CommandRegistry) Execute(scanner *bufio.Scanner, name string, args []string) bool {
	// Convert command name to lowercase for case-insensitive lookup
	lowerName := strings.ToLower(name)
	if cmd, exists := r.commands[lowerName]; exists {
		if quit := cmd.Execute(scanner, args); quit {
			return true
		}
	} else {
		r.unknownCommandHandler(name)
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

type SearchCommand struct {
	Handler handler.Handler
}

func (c *SearchCommand) Execute(scanner *bufio.Scanner, args []string) bool {
	c.Handler.HandleSearch(scanner, args)
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
	var rawURL string
	if len(args) > 0 {
		rawURL = args[0]
	}
	c.Handler.HandleUpsert(scanner, rawURL)
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

type HelpCommand struct {
	Handler handler.Handler
}

func (c *HelpCommand) Execute(scanner *bufio.Scanner, args []string) bool {
	c.Handler.HandleHelp()
	return false
}

type ClearCommand struct {
	Handler handler.Handler
}

func (c *ClearCommand) Execute(scanner *bufio.Scanner, args []string) bool {
	c.Handler.HandleClear()
	return false
}

type QuitCommand struct {
	Handler handler.Handler
}

func (c *QuitCommand) Execute(scanner *bufio.Scanner, args []string) bool {
	c.Handler.HandleQuit()
	return true
}

type HistoryCommand struct {
	Handler handler.Handler
}

func (c *HistoryCommand) Execute(scanner *bufio.Scanner, args []string) bool {
	c.Handler.HandleHistory()
	return false
}

type SettingCommand struct {
	Handler handler.Handler
}

func (c *SettingCommand) Execute(scanner *bufio.Scanner, args []string) bool {
	c.Handler.HandleSetting(scanner, args)
	return false
}

type VersionCommand struct {
	Handler handler.Handler
}

func (c *VersionCommand) Execute(scanner *bufio.Scanner, args []string) bool {
	c.Handler.HandleVersion()
	return false
}
