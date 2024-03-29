package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

// Exit and display an error message
func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}

type problem struct {
	q string
	a string
}

// Takes in the lines from the csv as a slice and outputs them
// in problem-struct form
func parseLines(lines [][]string) []problem {
	ret := make([]problem, len(lines))
	for i, line := range lines {
		ret[i] = problem{
			q: line[0],
			a: strings.TrimSpace(line[1]),
		}
	}
	return ret
}

// If the shuffle flag is set, this is called to shuffle the problems before displaying them
func shuffleProblems(probs []problem) {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	for len(probs) > 0 {
		n := len(probs)
		randIndex := r.Intn(n)
		probs[n-1], probs[randIndex] = probs[randIndex], probs[n-1]
		probs = probs[:n-1]
	}
}

func main() {
	csvFilename := flag.String("csv", "problems.csv", "a csv file in the format of 'question,answer'")
	timeLimit := flag.Int("limit", 30, "the time limit for the quiz in seconds")
	shuffle := flag.Bool("shuffle", false, "if true, shuffle the problems")
	flag.Parse()

	file, err := os.Open(*csvFilename)
	if err != nil {
		exit(fmt.Sprintf("Failed to open the csv file: %s\n", *csvFilename))
		os.Exit(1)
	}
	// Open the csv file with a reader
	r := csv.NewReader(file)

	// Parse the entire file into a slice in memory
	lines, err := r.ReadAll()
	if err != nil {
		exit("Failed to parse the provided csv file.")
	}
	problems := parseLines(lines)
	if *shuffle {
		shuffleProblems(problems)
	}

	// Start the timer for the quiz
	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)

	// Print the problems and get user response
	correct := 0

	for i, p := range problems {
		fmt.Printf("Problem #%d: %s = ", i+1, p.q) // i+1 so starts at problem 1
		// Setup the channel to listen for answers from the user
		answerCh := make(chan string)
		go func() {
			var answer string
			fmt.Scanf("%s\n", &answer)
			answerCh <- answer
		}()

		select {
		// If the timer ends the game, stop now
		case <-timer.C:
			fmt.Printf("\nYou scored %d out of %d.\n", correct, len(problems))
			return
		// if you get an inputted answer, check for correctness
		case answer := <-answerCh:
			if answer == p.a {
				correct++
			}
		}
	}
	fmt.Printf("You scored %d out of %d.\n", correct, len(problems))
}
