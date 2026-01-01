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

	"github.com/eannchen/leetsolv/command"
	"github.com/eannchen/leetsolv/config"
	"github.com/eannchen/leetsolv/core"
	"github.com/eannchen/leetsolv/handler"
	"github.com/eannchen/leetsolv/internal/clock"
	"github.com/eannchen/leetsolv/internal/fileutil"
	"github.com/eannchen/leetsolv/internal/logger"
	"github.com/eannchen/leetsolv/storage"
	"github.com/eannchen/leetsolv/usecase"
)

// Version information - will be set during build
var (
	Version   = "dev"
	BuildTime = "unknown"
	GitCommit = "unknown"
)

func main() {
	// Setup dependencies once
	clock := clock.NewClock()
	fileutil := fileutil.NewJSONFileUtil()
	cfg, err := config.NewConfig(fileutil)
	if err != nil {
		fmt.Println("Failed to load configuration:", err)
		os.Exit(1)
	}
	if err := logger.Init(cfg.InfoLogFile, cfg.ErrorLogFile); err != nil {
		fmt.Println("Failed to initialize logger:", err)
		os.Exit(1)
	}
	storage := storage.NewFileStorage(cfg.QuestionsFile, cfg.DeltasFile, fileutil)
	scheduler := core.NewSM2Scheduler(cfg, clock)
	questionUseCase := usecase.NewQuestionUseCase(cfg, storage, scheduler, clock)
	ioHandler := handler.NewIOHandler(clock)
	h := handler.NewHandler(cfg, questionUseCase, ioHandler, Version)

	commandRegistry := command.NewCommandRegistry(h.HandleUnknown)

	listCommand := &command.ListCommand{Handler: h}
	commandRegistry.Register("list", listCommand)
	commandRegistry.Register("ls", listCommand)

	searchCommand := &command.SearchCommand{Handler: h}
	commandRegistry.Register("search", searchCommand)
	commandRegistry.Register("s", searchCommand)

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

	historyCommand := &command.HistoryCommand{Handler: h}
	commandRegistry.Register("history", historyCommand)
	commandRegistry.Register("hist", historyCommand)
	commandRegistry.Register("log", historyCommand)

	settingCommand := &command.SettingCommand{Handler: h}
	commandRegistry.Register("setting", settingCommand)
	commandRegistry.Register("config", settingCommand)
	commandRegistry.Register("cfg", settingCommand)

	helpCommand := &command.HelpCommand{Handler: h}
	commandRegistry.Register("help", helpCommand)
	commandRegistry.Register("h", helpCommand)

	versionCommand := &command.VersionCommand{Handler: h}
	commandRegistry.Register("version", versionCommand)
	commandRegistry.Register("ver", versionCommand)
	commandRegistry.Register("v", versionCommand)

	migrateCommand := &command.MigrateCommand{Handler: h}
	commandRegistry.Register("migrate", migrateCommand)

	resetCommand := &command.ResetCommand{Handler: h}
	commandRegistry.Register("reset", resetCommand)

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
