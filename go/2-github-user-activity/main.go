package main

import (
	"fmt"
	"os"
)

func main() {
	args := os.Args

	if len(args) < 2 {
		fmt.Println("Please provide username")
		return
	}

	username := args[1]
	gt := NewGithubActivityFetcher()
	err := gt.FetchEvents(username)
	if err != nil {
		fmt.Println("We could not fetch data, ", err)
	}

	gt.PrintEvents(username)
}
