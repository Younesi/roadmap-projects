package main

import (
	"bufio"
	"fmt"
	"math/rand/v2"
	"os"
	"strconv"
	"strings"
	"time"
)

type DifficultyLevel int

const (
	Easy DifficultyLevel = iota + 1
	Medium
	Hard
)

func (d DifficultyLevel) String() string {
	return [...]string{"Easy", "Medium", "Hard"}[d-1]
}

func (d DifficultyLevel) Chances() int {
	return [...]int{10, 5, 3}[d-1]
}

func IsValidDifficulty(level int) bool {
	return level >= int(Easy) && level <= int(Hard)
}

func main() {
	fmt.Print(`
		Welcome to the Number Guessing Game!
		I'm thinking of a number between 1 and 100.
		You have 5 chances to guess the correct number.

		Please select the difficulty level:
		1. Easy (10 chances)
		2. Medium (5 chances)
		3. Hard (3 chances)
		
	`)

	reader := bufio.NewReader(os.Stdin)
	var level int
	var err error
	for {
		fmt.Print("Enter difficulty level (1, 2, or 3): ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		level, err = strconv.Atoi(input)
		if err == nil && IsValidDifficulty(level) {
			break
		}
		fmt.Println("Invalid input. Please enter 1, 2, or 3.")
	}

	difficultyLevel := DifficultyLevel(level)
	maxAttempts := difficultyLevel.Chances()

	for {
		playGame(reader, difficultyLevel, maxAttempts)

		fmt.Print("Do you want to continue? (yes/no): ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(strings.ToLower(input))
		if input != "yes" && input != "y" {
			break
		}
	}
}

func playGame(reader *bufio.Reader, difficultyLevel DifficultyLevel, maxAttempts int) {
	targetNumber := rand.IntN(100) + 1
	attempts := 0
	startTime := time.Now()

	fmt.Printf("Great! You have selected the %s difficulty level.\n", difficultyLevel.String())
	fmt.Println("Let's start the game!")

	for attempts < maxAttempts {
		fmt.Print("Enter your guess: ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		guess, err := strconv.Atoi(input)
		if err != nil {
			fmt.Println("Invalid input. Please enter a number.")
			continue
		}

		attempts++

		if targetNumber < guess {
			fmt.Printf("Incorrect! The number is less than %d.\n", guess)
		} else if targetNumber > guess {
			fmt.Printf("Incorrect! The number is greater than %d.\n", guess)
		} else {
			elapsedTime := time.Since(startTime)
			fmt.Printf("Congratulations! You guessed the correct number in %d attempts and %s.\n", attempts, elapsedTime.Round(time.Second))
			return
		}
	}

	fmt.Printf("Sorry, you've run out of attempts. The number was %d.\n", targetNumber)
}
