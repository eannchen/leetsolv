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
	ioHandler := handler.NewIOHandler()
	h := handler.NewHandler(ioHandler, questionUseCase)

	commandRegistry := command.NewCommandRegistry()
	commandRegistry.Register("list", &command.ListCommand{Handler: h})
	commandRegistry.Register("get", &command.GetCommand{Handler: h})
	commandRegistry.Register("status", &command.StatusCommand{Handler: h})
	commandRegistry.Register("upsert", &command.UpsertCommand{Handler: h})
	commandRegistry.Register("delete", &command.DeleteCommand{Handler: h})
	commandRegistry.Register("undo", &command.UndoCommand{Handler: h})
	commandRegistry.Register("quit", &command.QuitCommand{})

	scanner := bufio.NewScanner(os.Stdin)

	// --- CLI argument mode ---
	if len(os.Args) > 1 {
		// Handle "help" command
		if os.Args[1] == "help" {
			printHelp()
			os.Exit(0)
		}

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

	printWelcome()
	printInteractiveHelp()

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Shutting down gracefully...")
			return
		default:
			fmt.Print("\n> ")
			scanner.Scan()
			input := strings.TrimSpace(scanner.Text())

			if input == "" {
				continue
			}

			// Parse command and arguments
			parts := strings.Fields(input)
			cmd := parts[0]
			args := parts[1:]

			// Handle special commands
			switch cmd {
			case "help", "h":
				printInteractiveHelp()
				continue
			case "quit", "q", "exit":
				fmt.Println("Goodbye!")
				return
			case "clear", "cls":
				fmt.Print("\033[H\033[2J") // Clear screen
				printWelcome()
				printInteractiveHelp()
				continue
			}

			// Execute command
			if quit := commandRegistry.Execute(scanner, cmd, args); quit {
				return
			}
		}
	}
}

func printWelcome() {
	fmt.Println("╭────────────────────────────────────────╮")
	fmt.Println("│                                        │")
	fmt.Println("│    LeetSolv — CLI SRS for LeetCode     │")
	fmt.Println("│                                        │")
	fmt.Println("╰────────────────────────────────────────╯")
}

func printInteractiveHelp() {
	fmt.Println("\n📚 Available Commands:")
	fmt.Println("  status      - Show question status (due, upcoming, total)")
	fmt.Println("  list        - List all questions with pagination")
	fmt.Println("  get <id>    - Get details of a question by ID or URL")
	fmt.Println("  upsert      - Add or update a question")
	fmt.Println("  delete <id> - Delete a question by ID or URL")
	fmt.Println("  undo        - Undo the last action")
	fmt.Println("  help        - Show this help message")
	fmt.Println("  quit        - Exit the application")
	fmt.Println("  clear       - Clear the screen")
	fmt.Println("\n💡 Tips:")
	fmt.Println("  • Use 'h' for help, 'q' for quit")
	fmt.Println("  • Commands are case-insensitive")
	fmt.Println("  • Press Enter to continue pagination")
}

func printHelp() {
	fmt.Println("Usage: leetsolv [command]")
	fmt.Println("\nAvailable commands:")
	fmt.Println("  list       - List all questions with pagination.")
	fmt.Println("  get        - Get details of a question by ID or URL.")
	fmt.Println("  status     - Show the status of questions (due, upcoming, total).")
	fmt.Println("  upsert     - Add or update a question.")
	fmt.Println("  delete     - Delete a question by ID or URL.")
	fmt.Println("  undo       - Undo the last action.")
	fmt.Println("  help       - Show this help message.")
}
