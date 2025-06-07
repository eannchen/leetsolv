package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"leetsolv/commands"
	"leetsolv/core"
	"leetsolv/storage"
)

func main() {
	storage := &storage.FileStorage{File: "questions.json"}
	scheduler := core.SimpleScheduler{}
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Enter command (status/upsert/delete/quit): ")
		scanner.Scan()
		cmd := scanner.Text()
		switch cmd {
		case "status":
			due, upcoming, total, err := commands.ListQuestionsSummary(storage)
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}
			fmt.Printf("Total Questions: %d\n\n", total)

			fmt.Printf("Due Questions: %d\n", len(due))
			for _, q := range due {
				fmt.Printf("[%d] %s\n   Note: %s\n", q.ID, q.URL, q.Note)
			}

			fmt.Printf("\nUpcoming Questions (within 14 days): %d\n", len(upcoming))
			for _, q := range upcoming {
				fmt.Printf("[%d] %s (Next: %s)\n   Note: %s\n", q.ID, q.URL, q.NextReview.Format("2006-01-02"), q.Note)
			}
			fmt.Printf("\n")
		case "upsert":
			fmt.Print("URL: ")
			scanner.Scan()
			url := scanner.Text()
			fmt.Print("Note: ")
			scanner.Scan()
			note := scanner.Text()

			fmt.Println("Familiarity:")
			fmt.Println("1. Struggled    - Solved, but barely. Needed heavy effort or help.")
			fmt.Println("2. Clumsy       - Solved with partial understanding, some errors.")
			fmt.Println("3. Decent       - Solved mostly right, but not smooth.")
			fmt.Println("4. Smooth       - Solved confidently and clearly.")
			fmt.Println("5. Fluent       - Solved perfectly and instantly.")
			fmt.Print("\nEnter a number (1-5): ")
			scanner.Scan()
			famInput := scanner.Text()
			fam, err := strconv.Atoi(famInput)
			if err != nil || fam < 1 || fam > 5 {
				fmt.Println("Invalid familiarity level. Please enter a number between 1 and 5.")
				continue
			}

			// Adjust familiarity to match the `Familiarity` enum (0-based index)
			familiarity := core.Familiarity(fam - 1)

			if err := commands.UpsertQuestion(storage, scheduler, url, note, familiarity); err != nil {
				fmt.Println("Error:", err)
			} else {
				fmt.Println("Question upserted.")
			}
			fmt.Printf("\n")
		case "delete":
			fmt.Print("Enter URL or type '--last' to delete the most recently added: ")
			scanner.Scan()
			input := scanner.Text()

			// Confirm before deleting
			fmt.Print("Do you want to delete the question? [y/N]: ")
			scanner.Scan()
			confirm := strings.ToLower(strings.TrimSpace(scanner.Text()))
			if confirm != "y" && confirm != "Y" && confirm != "Yes" && confirm != "yes" {
				fmt.Println("Cancelled.")
				fmt.Printf("\n")
				continue
			}

			if err := commands.DeleteQuestion(storage, input); err != nil {
				fmt.Println("Error:", err)
			} else {
				fmt.Println("Question deleted.")
			}
			fmt.Printf("\n")
		case "quit":
			return
		default:
			fmt.Println("Unknown command.")
		}
	}
}
