package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	cyoa "github.com/zrnorth/gopher/gopher-3"
)

func main() {
	port := flag.Int("port", 3000, "port to start the webapp on")
	filename := flag.String("file", "story.json", "the json with the choose-your-own-adventure story")
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

	h := cyoa.NewHandler(story)
	fmt.Printf("Starting the server on port %d", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), h))
}
