package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
)

func main() {
	origin := flag.String("origin", "", "A URL of the server to which the requests will be forwarded.")
	port := flag.Int("port", 0, "A port on which the caching proxy server will run.")

	args := os.Args
	if len(args) < 2 {
		fmt.Println("please enter at least a command or run -h to see help")
		return
	}
	flag.Parse()

	if *origin == "" {
		fmt.Println("please enter the origin")
		return
	}

	if *port == 0 {
		fmt.Println("please enter the port")
		return
	}
	cps := NewCachingProxyServer(*origin)

	http.HandleFunc("/", cps.Start)

	go func() {
		http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
	}()

	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("Server is running on port : %d \n", *port)
	fmt.Println("Enter 'clear-cache' to clear the cache, or 'quit' to exit: ")

	for {
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		switch input {
		case "clear-cache":
			cps.ClearCache()
			fmt.Println("cache is cleared.")
		case "quit":
			return
		default:
			fmt.Println("Unknown command. Please try again.")
		}
	}
}
