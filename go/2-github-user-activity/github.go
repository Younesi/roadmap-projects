package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const BaseUri = "https://api.github.com"
const countOfOutputEvents = 7
const cacheEventTimeInMin = 30 * time.Minute

type EventsData struct {
	Events []string
}

type GithubActivityFetcher struct {
	username string
	Data     map[string]EventsData
}

func NewGithubActivityFetcher() *GithubActivityFetcher {
	return &GithubActivityFetcher{
		Data: make(map[string]EventsData),
	}
}

func (gt *GithubActivityFetcher) FetchEvents(username string) error {
	client := http.Client{}
	url := fmt.Sprintf("%s/users/%s/events", BaseUri, username)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("error: %s", err)
	}

	req.Header = http.Header{
		"Accept": {"application/vnd.github+json"},
	}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error: %s", err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("Error reading the response body: %w", err)

	}

	var events []map[string]interface{}
	err = json.Unmarshal(body, &events)
	if err != nil {
		return fmt.Errorf("error unmarshaling JSON: %w", err)
	}

	var eventsData EventsData
	for _, event := range events {
		var desc = getEventDescription(event)

		if len(desc) != 0 {
			eventsData.Events = append(eventsData.Events, desc)
		}
	}

	if len(eventsData.Events) != 0 {
		gt.Data[username] = eventsData
	}
	return nil
}

func (gt *GithubActivityFetcher) PrintEvents(username string) {
	eventsData, ok := gt.Data[username]
	if !ok || len(eventsData.Events) == 0 {
		fmt.Println("No evets found")
		return
	}

	count := 0
	for _, data := range eventsData.Events {
		if count > countOfOutputEvents {
			return
		}
		fmt.Println(data)
		count++
	}
}

func getEventDescription(event map[string]interface{}) string {
	eventType, _ := event["type"].(string)
	repo, _ := event["repo"].(map[string]interface{})
	repoName, _ := repo["name"].(string)
	payload, _ := event["payload"].(map[string]interface{})

	caser := cases.Title(language.English)

	var result string
	switch eventType {
	case "PushEvent":
		commits, _ := payload["commits"].([]interface{})
		result = fmt.Sprintf("Pushed %d commits to %s", len(commits), repoName)

	case "IssuesEvent", "PullRequestEvent":
		action, _ := payload["action"].(string)
		itemType := "issue"
		if eventType == "PullRequestEvent" {
			itemType = "pull request"
		}
		switch action {
		case "opened":
			result = fmt.Sprintf("Opened a new %s in %s", itemType, repoName)
		case "closed", "reopened":
			result = fmt.Sprintf("%s a %s in %s", caser.String(action), itemType, repoName)
		case "merged":
			if eventType == "PullRequestEvent" {
				result = fmt.Sprintf("Merged a pull request in %s", repoName)
			}
		default:
			result = fmt.Sprintf("%s a %s in %s", caser.String(action), itemType, repoName)
		}

	case "WatchEvent":
		action, _ := payload["action"].(string)
		if action == "started" {
			result = fmt.Sprintf("Starred %s", repoName)
		}

	case "ForkEvent":
		result = fmt.Sprintf("Forked %s", repoName)

	case "CreateEvent", "DeleteEvent":
		refType, _ := payload["ref_type"].(string)
		if eventType == "CreateEvent" {
			result = fmt.Sprintf("Created a new %s in %s", refType, repoName)
		} else {
			result = fmt.Sprintf("Deleted a %s in %s", refType, repoName)
		}

	case "CommitCommentEvent":
		result = fmt.Sprintf("Commented on a commit in %s", repoName)

	case "IssueCommentEvent":
		result = fmt.Sprintf("Commented on an issue in %s", repoName)

	case "ReleaseEvent":
		action, _ := payload["action"].(string)
		if action == "published" {
			result = fmt.Sprintf("Published a new release in %s", repoName)
		}

	default:
		eventName := strings.TrimSuffix(eventType, "Event")
		result = fmt.Sprintf("Performed %s on %s", eventName, repoName)
	}

	return result
}
