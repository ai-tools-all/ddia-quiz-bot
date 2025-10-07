package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/your-username/ddia-quiz-bot/internal/config"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Printf("Error loading .env: %v\n", err)
	}
	
	s, err := config.LoadSchedule("./content/schedule.yml")
	if err != nil {
		fmt.Printf("Error loading schedule: %v\n", err)
	} else {
		fmt.Printf("Start date: %s\n", s.Chapters[0].StartDate)
	}
}
