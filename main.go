package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"leetsolv/core"
	"leetsolv/handler"
	"leetsolv/storage"
	"leetsolv/usecase"
)

func main() {
	storage := storage.NewFileStorage("questions.json", "snapshots.json")
	scheduler := core.NewSM2Scheduler()
	questionUseCase := usecase.NewQuestionUseCase(storage, scheduler)
	h := handler.NewHandler(questionUseCase)

	scanner := bufio.NewScanner(os.Stdin)
	for {
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
