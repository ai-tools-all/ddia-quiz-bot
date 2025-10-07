package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/your-username/ddia-quiz-bot/internal/config"
	"github.com/your-username/ddia-quiz-bot/internal/store"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Printf("Error loading .env: %v\n", err)
	}
	
	s, err := config.LoadSchedule("./content/schedule.yml")
	if err != nil {
		fmt.Printf("Error loading schedule: %v\n", err)
		return
	}
	
	store, err := store.NewContentStore("./content")
	if err != nil {
		fmt.Printf("Error loading store: %v\n", err)
		return
	}
	
	fmt.Printf("Schedule chapters: %d\n", len(s.Chapters))
	fmt.Printf("Questions in store: %d\n", len(store.QuestionsByID))
	for id, q := range store.QuestionsByID {
		fmt.Printf("Question ID: %s, File: %s, Day: %d\n", id, q.Path, q.Day)
	}
	
	fmt.Printf("Schedule expects file: %s\n", s.Chapters[0].Questions[1].File)
}
