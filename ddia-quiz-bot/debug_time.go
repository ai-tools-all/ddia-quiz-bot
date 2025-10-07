package main

import (
	"fmt"
	"time"
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
		return
	}
	
	chap := s.Chapters[0]
	startDate, _ := chap.GetStartDate()
	fmt.Printf("Start date: %s\n", startDate)
	
	qSched := chap.Questions[1] // day 2
	fmt.Printf("Question day: %d\n", qSched.Day)
	
	postDate := startDate.AddDate(0, 0, qSched.Day-1)
	postTime := time.Date(postDate.Year(), postDate.Month(), postDate.Day(), 10, 0, 0, 0, time.UTC)
	fmt.Printf("Post time: %s\n", postTime)
	
	now := time.Now().UTC()
	fmt.Printf("Now: %s\n", now)
	fmt.Printf("Should post: %v\n", now.After(postTime))
}
