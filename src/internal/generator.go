package internal

import (
	"bufio"
	"math/rand"
	"strings"
	"sync"
	"time"

	"gti/src/assets"
)

var defaultWords = []string{
	"the", "quick", "brown", "fox", "jumps", "over", "lazy", "dog",
	"hello", "world", "typing", "speed", "test", "practice", "accuracy",
	"keyboard", "computer", "software", "development", "programming",
	"go", "language", "bubble", "tea", "terminal", "user", "interface",
}

var languageFiles = map[string]string{
	"english":    "eng",
	"spanish":    "spa",
	"french":     "fre",
	"german":     "ger",
	"japanese":   "jap",
	"russian":    "ru",
	"italian":    "ita",
	"portuguese": "por",
	"chinese":    "chi",
	"arabic":     "ara",
	"hindi":      "hin",
	"korean":     "kor",
	"dutch":      "dut",
	"swedish":    "swe",
	"czech":      "cze",
	"danish":     "dan",
	"finnish":    "fin",
	"greek":      "gre",
	"hebrew":     "heb",
	"hungarian":  "hun",
	"norwegian":  "nor",
	"polish":     "pol",
	"thai":       "tha",
	"turkish":    "tur",
	"random":     "ran",
}

var loadedWords = make(map[string][]string)
var loadMutex sync.Mutex

func loadWords(language string) []string {
	loadMutex.Lock()
	defer loadMutex.Unlock()

	if words, exists := loadedWords[language]; exists {
		return words
	}

	fileName, exists := languageFiles[language]
	if !exists {
		fileName = languageFiles["random"]
	}

	filePath := "words/" + fileName

	data, err := assets.Words.ReadFile(filePath)
	if err != nil {
		return defaultWords
	}

	var words []string
	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			words = append(words, line)
		}
	}

	if len(words) == 0 {
		words = defaultWords
	}

	loadedWords[language] = words
	return words
}

func GenerateWord(language string) string {
	rand.Seed(time.Now().UnixNano())
	words := loadWords(language)
	return words[rand.Intn(len(words))]
}

func GenerateWordsDynamic(count int, language string) string {
	rand.Seed(time.Now().UnixNano())
	var selected []string
	for i := 0; i < count; i++ {
		selected = append(selected, GenerateWord(language))
	}
	return strings.Join(selected, " ")
}

func IsLanguageSupported(language string) bool {
	_, exists := languageFiles[language]
	return exists
}
