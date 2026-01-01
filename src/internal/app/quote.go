// Package app provides functionality for fetching quotes from external APIs.
package app

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"gti/src/internal/config"
	"gti/src/internal/session"
)

type QuoteResponse struct {
	Q string `json:"q"`
	A string `json:"a"`
}

func FetchQuote(cfg *config.Config) string {
	q := FetchQuoteWithAuthor(cfg)
	return q.Text
}

func FetchQuoteWithAuthor(cfg *config.Config) session.Quote {
	client := &http.Client{
		Timeout: time.Duration(cfg.Network.TimeoutMs) * time.Millisecond,
	}

	resp, err := client.Get("https://zenquotes.io/api/random")
	if err != nil {
		return session.Quote{Text: config.DefaultPracticeText, Author: "Unknown"}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return session.Quote{Text: config.DefaultPracticeText, Author: "Unknown"}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return session.Quote{Text: config.DefaultPracticeText, Author: "Unknown"}
	}

	var quotes []QuoteResponse
	err = json.Unmarshal(body, &quotes)
	if err != nil {
		return session.Quote{Text: config.DefaultPracticeText, Author: "Unknown"}
	}

	if len(quotes) == 0 {
		return session.Quote{Text: config.DefaultPracticeText, Author: "Unknown"}
	}

	quote := quotes[0]
	if quote.Q == "" {
		return session.Quote{Text: config.DefaultPracticeText, Author: "Unknown"}
	}

	return session.Quote{Text: quote.Q, Author: quote.A}
}

func FetchMultipleQuotes(cfg *config.Config, count int) []session.Quote {
	if count <= 0 {
		count = 1
	}
	if count > 10 {
		count = 10
	}

	var quotes []session.Quote
	client := &http.Client{
		Timeout: time.Duration(cfg.Network.TimeoutMs) * time.Millisecond,
	}

	for i := 0; i < count; i++ {
		resp, err := client.Get("https://zenquotes.io/api/random")
		if err != nil {
			quotes = append(quotes, session.Quote{Text: config.DefaultPracticeText, Author: "Unknown"})
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			quotes = append(quotes, session.Quote{Text: config.DefaultPracticeText, Author: "Unknown"})
			continue
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			quotes = append(quotes, session.Quote{Text: config.DefaultPracticeText, Author: "Unknown"})
			continue
		}

		var quoteResponses []QuoteResponse
		err = json.Unmarshal(body, &quoteResponses)
		if err != nil || len(quoteResponses) == 0 {
			quotes = append(quotes, session.Quote{Text: config.DefaultPracticeText, Author: "Unknown"})
			continue
		}

		qr := quoteResponses[0]
		if qr.Q != "" {
			quotes = append(quotes, session.Quote{Text: qr.Q, Author: qr.A})
		} else {
			quotes = append(quotes, session.Quote{Text: config.DefaultPracticeText, Author: "Unknown"})
		}
	}

	if len(quotes) == 0 {
		return []session.Quote{{Text: config.DefaultPracticeText, Author: "Unknown"}}
	}

	return quotes
}
