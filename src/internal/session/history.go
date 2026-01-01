package session

import (
	"bufio"
	"encoding/json"
	"os"
	"sort"
	"time"

	"gti/src/internal/config"
)

type SessionRecord struct {
	Timestamp   time.Time `json:"timestamp"`
	Mode        string    `json:"mode"`
	TextLength  int       `json:"text_length"`
	DurationMs  int64     `json:"duration_ms"`
	WPM         float64   `json:"wpm"`
	CPM         float64   `json:"cpm"`
	Accuracy    float64   `json:"accuracy"`
	Mistakes    int       `json:"mistakes"`
	Tier        string    `json:"tier,omitempty"`
	QuoteAuthor string    `json:"quote_author,omitempty"`

	NetWPM            float64 `json:"net_wpm,omitempty"`
	AdjustedWPM       float64 `json:"adjusted_wpm,omitempty"`
	CorrectedErrors   int     `json:"corrected_errors,omitempty"`
	UncorrectedErrors int     `json:"uncorrected_errors,omitempty"`
	BackspaceCount    int     `json:"backspace_count,omitempty"`
	AvgWordLength     float64 `json:"avg_word_length,omitempty"`
}

func SaveSessionRecord(cfg *config.Config, record *SessionRecord) error {
	if !cfg.History.Enabled {
		return nil
	}

	filePath := config.ExpandPath(cfg.History.File)
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	record.Timestamp = time.Now()
	data, err := json.Marshal(record)
	if err != nil {
		return err
	}

	_, err = file.WriteString(string(data) + "\n")
	return err
}

func LoadSessionRecords(cfg *config.Config) ([]*SessionRecord, error) {
	if !cfg.History.Enabled {
		return nil, nil
	}

	filePath := config.ExpandPath(cfg.History.File)
	file, err := os.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return []*SessionRecord{}, nil
		}
		return nil, err
	}
	defer file.Close()

	var records []*SessionRecord
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var record SessionRecord
		if err := json.Unmarshal(scanner.Bytes(), &record); err != nil {
			continue
		}
		records = append(records, &record)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	sort.Slice(records, func(i, j int) bool {
		return records[i].Timestamp.After(records[j].Timestamp)
	})

	return records, nil
}

func CalculateStreaks(validSessions []*SessionRecord) (int, int) {
	if len(validSessions) == 0 {
		return 0, 0
	}

	dates := extractUniqueDates(validSessions)
	if len(dates) == 0 {
		return 0, 0
	}

	currentStreak := calculateCurrentStreak(dates)
	longestStreak := calculateLongestStreak(dates)

	return currentStreak, longestStreak
}

func extractUniqueDates(sessions []*SessionRecord) []string {
	dateMap := make(map[string]bool)
	for _, session := range sessions {
		date := session.Timestamp.Format("2006-01-02")
		dateMap[date] = true
	}

	var dates []string
	for date := range dateMap {
		dates = append(dates, date)
	}
	sort.Strings(dates)
	return dates
}

func calculateCurrentStreak(dates []string) int {
	now := time.Now()
	yesterday := now.AddDate(0, 0, -1)
	todayStr := now.Format("2006-01-02")
	yesterdayStr := yesterday.Format("2006-01-02")

	hasRecentActivity := false
	for _, date := range dates {
		if date == todayStr || date == yesterdayStr {
			hasRecentActivity = true
			break
		}
	}

	if !hasRecentActivity {
		return 0
	}

	currentStreak := 0
	for i := len(dates) - 1; i >= 0; i-- {
		expectedDate := now.AddDate(0, 0, -(len(dates) - 1 - i)).Format("2006-01-02")
		if dates[i] == expectedDate {
			currentStreak++
		} else {
			break
		}
	}

	return currentStreak
}

func calculateLongestStreak(dates []string) int {
	if len(dates) == 0 {
		return 0
	}

	longestStreak := 1
	current := 1

	for i := 1; i < len(dates); i++ {
		prevDate, _ := time.Parse("2006-01-02", dates[i-1])
		currDate, _ := time.Parse("2006-01-02", dates[i])

		if currDate.Sub(prevDate).Hours() == 24 {
			current++
		} else {
			if current > longestStreak {
				longestStreak = current
			}
			current = 1
		}
	}

	if current > longestStreak {
		longestStreak = current
	}

	return longestStreak
}
