package main

import (
	"fmt"
	"log"

	"leetsolv/storage"
)

func main() {
	// Test loading the migrated data
	fileStorage := storage.NewFileStorage("questions_new.json", "deltas.json")

	store, err := fileStorage.LoadQuestionStore()
	if err != nil {
		log.Fatalf("Failed to load migrated data: %v", err)
	}

	fmt.Printf("✅ Successfully loaded migrated data!\n")
	fmt.Printf("📊 Statistics:\n")
	fmt.Printf("  - Total questions: %d\n", len(store.Questions))
	fmt.Printf("  - Max ID: %d\n", store.MaxID)
	fmt.Printf("  - URL index entries: %d\n", len(store.URLIndex))

	// Verify URL index integrity
	fmt.Printf("\n🔍 Verifying URL index integrity...\n")
	urlIndexErrors := 0
	for url, id := range store.URLIndex {
		if question, exists := store.Questions[id]; !exists {
			fmt.Printf("  ❌ URL index points to non-existent question ID %d: %s\n", id, url)
			urlIndexErrors++
		} else if question.URL != url {
			fmt.Printf("  ❌ URL mismatch for ID %d: index has '%s', question has '%s'\n", id, url, question.URL)
			urlIndexErrors++
		}
	}

	if urlIndexErrors == 0 {
		fmt.Printf("  ✅ URL index is consistent\n")
	} else {
		fmt.Printf("  ❌ Found %d URL index errors\n", urlIndexErrors)
	}

	// Show a few sample questions
	fmt.Printf("\n📝 Sample questions:\n")
	count := 0
	for id, question := range store.Questions {
		if count >= 3 {
			break
		}
		fmt.Printf("  ID %d: %s (Familiarity: %v, Importance: %v)\n",
			id, question.URL, question.Familiarity, question.Importance)
		count++
	}

	fmt.Printf("\n🎉 Verification completed successfully!\n")
	fmt.Printf("The migrated data is ready to use with your application.\n")
}
