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
	"leetsolv/storage"
	"leetsolv/usecase"
)

func main() {
	// Setup dependencies once
	env := config.Env()
	storage := storage.NewFileStorage(env.QuestionsFile, env.SnapshotsFile)
	scheduler := core.NewSM2Scheduler()
	questionUseCase := usecase.NewQuestionUseCase(storage, scheduler)
	ioHandler := handler.NewIOHandler()
	h := handler.NewHandler(ioHandler, questionUseCase)

	commandRegistry := command.NewCommandRegistry()
	commandRegistry.Register("list", &command.ListCommand{Handler: h})
	commandRegistry.Register("status", &command.StatusCommand{Handler: h})
	commandRegistry.Register("upsert", &command.UpsertCommand{Handler: h})
	commandRegistry.Register("delete", &command.DeleteCommand{Handler: h})
	commandRegistry.Register("undo", &command.UndoCommand{Handler: h})
	commandRegistry.Register("quit", &command.QuitCommand{})

	scanner := bufio.NewScanner(os.Stdin)

	// Check for command-line arguments
	if len(os.Args) > 1 {
		// Handle "help" command
		if os.Args[1] == "help" {
			printHelp()
			os.Exit(0)
		}

		// Combine all arguments into a single command string
		cmd := strings.Join(os.Args[1:], " ")
		commandRegistry.Execute(cmd, scanner)
		os.Exit(0)
	}

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

	// Interactive mode
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Shutting down gracefully...")
			return
		default:
			fmt.Print("Enter command (status/list/upsert/delete/undo/quit): ")
			scanner.Scan()
			cmd := strings.TrimSpace(scanner.Text())
			if quit := commandRegistry.Execute(cmd, scanner); quit {
				return
			}
		}
	}
}

func printHelp() {
	fmt.Println("Usage: leetsolv [command]")
	fmt.Println("\nAvailable commands:")
	fmt.Println("  list       - List all questions with pagination.")
	fmt.Println("  status     - Show the status of questions (due, upcoming, total).")
	fmt.Println("  upsert     - Add or update a question.")
	fmt.Println("  delete     - Delete a question by ID or URL.")
	fmt.Println("  undo       - Undo the last action.")
	fmt.Println("  help       - Show this help message.")
}
