package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"leetsolv/command"
	"leetsolv/config"
	"leetsolv/core"
	"leetsolv/handler"
	"leetsolv/internal/clock"
	"leetsolv/storage"
	"leetsolv/usecase"
)

func main() {
	// Setup dependencies once
	env := config.Env()
	clock := clock.NewClock()
	storage := storage.NewFileStorage(env.QuestionsFile, env.DeltasFile)
	scheduler := core.NewSM2Scheduler(clock)
	questionUseCase := usecase.NewQuestionUseCase(storage, scheduler, clock)
	ioHandler := handler.NewIOHandler(clock)
	h := handler.NewHandler(ioHandler, questionUseCase)

	commandRegistry := command.NewCommandRegistry(h.HandleUnknown)

	listCommand := &command.ListCommand{Handler: h}
	commandRegistry.Register("list", listCommand)
	commandRegistry.Register("ls", listCommand)

	getCommand := &command.GetCommand{Handler: h}
	commandRegistry.Register("detail", getCommand)
	commandRegistry.Register("get", getCommand)

	statusCommand := &command.StatusCommand{Handler: h}
	commandRegistry.Register("status", statusCommand)
	commandRegistry.Register("stat", statusCommand)

	upsertCommand := &command.UpsertCommand{Handler: h}
	commandRegistry.Register("upsert", upsertCommand)
	commandRegistry.Register("add", upsertCommand)

	deleteCommand := &command.DeleteCommand{Handler: h}
	commandRegistry.Register("remove", deleteCommand)
	commandRegistry.Register("rm", deleteCommand)
	commandRegistry.Register("delete", deleteCommand)
	commandRegistry.Register("del", deleteCommand)

	undoCommand := &command.UndoCommand{Handler: h}
	commandRegistry.Register("undo", undoCommand)
	commandRegistry.Register("back", undoCommand)

	helpCommand := &command.HelpCommand{Handler: h}
	commandRegistry.Register("help", helpCommand)
	commandRegistry.Register("h", helpCommand)

	clearCommand := &command.ClearCommand{Handler: h}
	commandRegistry.Register("clear", clearCommand)
	commandRegistry.Register("cls", clearCommand)

	quitCommand := &command.QuitCommand{Handler: h}
	commandRegistry.Register("quit", quitCommand)
	commandRegistry.Register("q", quitCommand)
	commandRegistry.Register("exit", quitCommand)

	scanner := bufio.NewScanner(os.Stdin)

	// --- CLI argument mode ---
	if len(os.Args) > 1 {
		// Combine all arguments into a single command string
		commandRegistry.Execute(scanner, os.Args[1], os.Args[2:])
		os.Exit(0)
	}

	// --- Interactive mode ---

	// Set up graceful shutdown signal listener
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		<-signalChan
		fmt.Println("\nReceived shutdown signal. Please wait...")
		cancel() // Cancel the context

		// timeout
		time.Sleep(5 * time.Second)
		os.Exit(0)
	}()

	h.HandleHelp()

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Shutting down gracefully...")
			return
		default:
			fmt.Print(prompt())
			scanner.Scan()

			input := strings.TrimSpace(scanner.Text())
			if input == "" {
				continue
			}

			// Parse command and arguments
			parts := strings.Fields(input)
			cmd := parts[0]
			args := parts[1:]

			// Execute command
			if quit := commandRegistry.Execute(scanner, cmd, args); quit {
				return
			}
		}
	}
}

func prompt() string {
	return "\nleetsolv ❯ "
}
