package main

import (
	"encoding/json"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
)

// This data management is from Paul Thiel, from GoLab Hackathon Dinner 2023
// Permission is being sought for its inclusion in this repository
// https://github.com/thielepaul/golab-schedule-2023/blob/main/main.go

type pageData struct {
	Props props `json:"props"`
}

type props struct {
	PageProps pageProps `json:"pageProps"`
}

type pageProps struct {
	Edition edition `json:"edition"`
}

type edition struct {
	Days []day `json:"days"`
}

type day struct {
	Title    string   `json:"title"`
	Schedule []Record `json:"schedule"`
}

type Record struct {
	Id                string    `json:"id"`
	Title             string    `json:"title"`
	Time              time.Time `json:"time"`
	DurationInMinutes int       `json:"durationInMinutes"`
	Text              string    `json:"text"`
}

func getData(scheduleURL string) ([]day, error) {
	resp, err := http.Get(scheduleURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	scheduleHtml, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	scheduleJson, err := extractJSON(string(scheduleHtml))
	if err != nil {
		return nil, err
	}

	var data pageData
	if err := json.Unmarshal([]byte(scheduleJson), &data); err != nil {
		return nil, err
	}

	return data.Props.PageProps.Edition.Days, nil
}

func extractJSON(htmlString string) (string, error) {
	re := regexp.MustCompile(`<script id="__NEXT_DATA__" type="application/json">(.*?)</script>`)
	match := re.FindStringSubmatch(htmlString)
	if len(match) < 2 {
		return "", nil // or return an error
	}
	return strings.TrimSpace(match[1]), nil
}
