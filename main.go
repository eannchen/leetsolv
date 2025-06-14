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

	"leetsolv/core"
	"leetsolv/handler"
	"leetsolv/storage"
	"leetsolv/usecase"
)

func main() {
	// Setup dependencies once
	storage := storage.NewFileStorage("questions.json", "snapshots.json")
	scheduler := core.NewSM2Scheduler()
	questionUseCase := usecase.NewQuestionUseCase(storage, scheduler)
	h := handler.NewHandler(questionUseCase)

	scanner := bufio.NewScanner(os.Stdin)

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

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Shutting down gracefully...")
			return
		default:
			fmt.Print("Enter command (status/list/upsert/delete/undo/quit): ")
			scanner.Scan()
			cmd := strings.TrimSpace(scanner.Text())
			switch cmd {
			case "list":
				h.HandleList(scanner)
			case "status":
				h.HandleStatus()
			case "upsert":
				h.HandleUpsert(scanner)
			case "delete":
				h.HandleDelete(scanner)
			case "undo":
				h.HandleUndo(scanner)
			case "quit":
				return
			default:
				fmt.Println("Unknown command.")
			}
		}
	}
}
