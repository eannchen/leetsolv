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
	scheduler := core.SM2Scheduler{}
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Enter command (status/list/upsert/delete/quit): ")
		scanner.Scan()
		cmd := scanner.Text()
		switch cmd {
		case "list":
			const pageSize = 5
			page := 0

			for {
				questions, totalPages, err := commands.PaginatedListQuestions(storage, pageSize, page)
				if err != nil {
					fmt.Println("Error:", err)
					break
				}

				if len(questions) == 0 {
					fmt.Println("No questions available.")
					break
				}

				// Display the current page
				fmt.Printf("-- Page %d/%d --\n", page+1, totalPages)
				for _, q := range questions {
					fmt.Printf("[%d] %s (Next: %s)\n", q.ID, q.URL, q.NextReview.Format("2006-01-02"))
					fmt.Printf("   Note: %s\n", q.Note)
				}

				// Handle user input for pagination
				if page+1 == totalPages {
					fmt.Println("\nEnd of list.")
					break
				}

				fmt.Print("\nPress [Enter] for next page, [q] to quit: ")
				scanner.Scan()
				input := strings.TrimSpace(scanner.Text())
				if input == "q" {
					break
				}

				page++
			}
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
			url := readLine(scanner, "URL: ")
			note := readLine(scanner, "Note: ")

			fmt.Println("Familiarity:")
			fmt.Println("1. Struggled    - Solved, but barely. Needed heavy effort or help.")
			fmt.Println("2. Clumsy       - Solved with partial understanding, some errors.")
			fmt.Println("3. Decent       - Solved mostly right, but not smooth.")
			fmt.Println("4. Smooth       - Solved confidently and clearly.")
			fmt.Println("5. Fluent       - Solved perfectly and instantly.")
			famInput := readLine(scanner, "\nEnter a number (1-5): ")
			fam, err := strconv.Atoi(famInput)
			if err != nil || fam < 1 || fam > 5 {
				fmt.Println("Invalid familiarity level. Please enter a number between 1 and 5.")
				continue
			}

			// Adjust familiarity to match the `Familiarity` enum (0-based index)
			familiarity := core.Familiarity(fam - 1)

			// Call the updated UpsertQuestion function
			upsertedQuestion, err := commands.UpsertQuestion(storage, scheduler, url, note, familiarity)
			if err != nil {
				fmt.Println("Error:", err)
			} else {
				// Display the upserted question
				fmt.Println("Question upserted:")
				fmt.Printf("[%d] %s\n", upsertedQuestion.ID, upsertedQuestion.URL)
				fmt.Printf("   Note: %s\n", upsertedQuestion.Note)
				fmt.Printf("   Familiarity: %d\n", upsertedQuestion.Familiarity)
				fmt.Printf("   Last Reviewed: %s\n", upsertedQuestion.LastReviewed.Format("2006-01-02 15:04:05"))
				fmt.Printf("   Next Review: %s\n", upsertedQuestion.NextReview.Format("2006-01-02 15:04:05"))
				fmt.Printf("   Review Count: %d\n", upsertedQuestion.ReviewCount)
				fmt.Printf("   Ease Factor: %.2f\n", upsertedQuestion.EaseFactor)
				fmt.Printf("   Created At: %s\n", upsertedQuestion.CreatedAt.Format("2006-01-02 15:04:05"))
			}
			fmt.Printf("\n")
		case "delete":
			input := readLine(scanner, "Enter ID, URL or type '--last' to delete the most recently added: ")

			// Confirm before deleting
			confirm := strings.ToLower(readLine(scanner, "Do you want to delete the question? [y/N]: "))
			if confirm != "y" && confirm != "yes" {
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

func readLine(scanner *bufio.Scanner, prompt string) string {
	fmt.Print(prompt)
	scanner.Scan()
	return strings.TrimSpace(scanner.Text())
}
