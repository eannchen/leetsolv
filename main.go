package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"leetsolv/core"
	"leetsolv/handler"
	"leetsolv/storage"
)

func main() {
	storage := &storage.FileStorage{QuestionsFile: "questions.json", SnapshotsFile: "snapshots.json"}
	scheduler := core.NewSM2Scheduler()
	scanner := bufio.NewScanner(os.Stdin)

	h := handler.NewHandler(storage, scheduler)

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
