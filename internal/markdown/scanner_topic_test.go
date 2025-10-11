package markdown

import (
	"testing"
)

func TestDiscoverTopics(t *testing.T) {
	scanner := NewScanner("")
	topics, err := scanner.DiscoverTopics("../../ddia-quiz-bot/content/chapters")
	
	if err != nil {
		t.Fatalf("DiscoverTopics failed: %v", err)
	}
	
	if len(topics) == 0 {
		t.Fatal("Expected at least one topic, got none")
	}
	
	t.Logf("Discovered %d topics", len(topics))
	
	for _, topic := range topics {
		t.Logf("Topic: %s | Display: %s | Path: %s | Questions: %d", 
			topic.Name, topic.DisplayName, topic.Path, topic.TotalCount)
		
		if topic.Name == "" {
			t.Errorf("Topic name is empty")
		}
		if topic.DisplayName == "" {
			t.Errorf("Topic display name is empty")
		}
		if topic.Path == "" {
			t.Errorf("Topic path is empty")
		}
		if topic.TotalCount == 0 {
			t.Errorf("Topic %s has no questions", topic.Name)
		}
		
		// Check level counts
		totalFromLevels := 0
		for level, count := range topic.LevelCounts {
			t.Logf("  Level %s: %d questions", level, count)
			totalFromLevels += count
		}
		
		if totalFromLevels != topic.TotalCount {
			t.Errorf("Level counts (%d) don't match total count (%d)", totalFromLevels, topic.TotalCount)
		}
	}
}

func TestGetProgressiveQuestions(t *testing.T) {
	// Test with GFS topic
	scanner := NewScanner("../../ddia-quiz-bot/content/chapters/09-distributed-systems-gfs/subjective")
	index, err := scanner.ScanQuestions()
	
	if err != nil {
		t.Fatalf("ScanQuestions failed: %v", err)
	}
	
	questions := scanner.GetProgressiveQuestions(index)
	
	if len(questions) == 0 {
		t.Fatal("Expected questions, got none")
	}
	
	t.Logf("Got %d questions in progressive order", len(questions))
	
	// Verify progressive ordering: L3 -> L4 -> L5 -> L6 -> L7
	levelOrder := []string{"L3", "L4", "L5", "L6", "L7"}
	currentLevelIdx := 0
	
	for i, q := range questions {
		t.Logf("Question %d: ID=%s, Level=%s, Category=%s", i+1, q.ID, q.Level, q.Category)
		
		// Check if we're on the expected level or higher
		foundLevel := false
		for idx := currentLevelIdx; idx < len(levelOrder); idx++ {
			if q.Level == levelOrder[idx] {
				currentLevelIdx = idx
				foundLevel = true
				break
			}
		}
		
		if !foundLevel {
			t.Errorf("Question %d has level %s, expected level >= %s", i+1, q.Level, levelOrder[currentLevelIdx])
		}
		
		// Within same level, baseline should come before bar-raiser
		if i > 0 && questions[i-1].Level == q.Level {
			prevCat := questions[i-1].Category
			currCat := q.Category
			if prevCat == "bar-raiser" && currCat == "baseline" {
				t.Errorf("Question %d: baseline should come before bar-raiser within level %s", i+1, q.Level)
			}
		}
	}
}

func TestFormatTopicName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"03-storage-and-retrieval", "Storage And Retrieval"},
		{"09-distributed-systems-gfs", "Distributed Systems Gfs"},
		{"simple", "simple"},
	}
	
	for _, test := range tests {
		result := formatTopicName(test.input)
		if result != test.expected {
			t.Errorf("formatTopicName(%s) = %s, expected %s", test.input, result, test.expected)
		}
	}
}
