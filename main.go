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

	for {
		fmt.Print("Enter command (status/list/upsert/delete/undo/quit): ")
		scanner.Scan()
		cmd := strings.TrimSpace(scanner.Text())
		switch cmd {
		case "list":
			handler.HandleList(scanner, storage)
		case "status":
			handler.HandleStatus(storage)
		case "upsert":
			handler.HandleUpsert(scanner, storage, scheduler)
		case "delete":
			handler.HandleDelete(scanner, storage)
		case "undo":
			handler.HandleUndo(scanner, storage)
		case "quit":
			return
		default:
			fmt.Println("Unknown command.")
		}
	}
}
