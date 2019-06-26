package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	cyoa "github.com/zrnorth/gopher/gopher-3"
)

func main() {
	port := flag.Int("port", 3000, "port to start the webapp on")
	filename := flag.String("file", "story.json", "the json with the choose-your-own-adventure story")
	commandLineMode := flag.Bool("cmd", false, "run in command line mode instead of as a web server")
	flag.Parse()
	fmt.Printf("Using the story %s.\n", *filename)

	f, err := os.Open(*filename)
	if err != nil {
		panic(err) // Shouldn't do this really but whatever
	}
	story, err := cyoa.JSONStory(f)
	if err != nil {
		panic(err)
	}

	if *commandLineMode {
		h := cyoa.NewHandler(story, cyoa.WithCommandLineMode())
		// Serve the initial text
		path := "intro" // this is bad obviously
		h.ServeTextToConsole(path)

		// When path comes back empty we have finished the game.
		for path != "" {
			// Setup a listener for user input on cmd line
			responseCh := make(chan string)
			go func() {
				var response string
				fmt.Scanf("%s\n", &response)
				responseCh <- response
			}()

			select {
			case response := <-responseCh:
				responseAsInt, _ := strconv.Atoi(response)
				path = h.GetNext(path, responseAsInt)

				h.ServeTextToConsole(path)
			}
		}
	} else {
		h := cyoa.NewHandler(story)

		fmt.Printf("Starting the server on port %d", *port)
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), h))
	}
}
