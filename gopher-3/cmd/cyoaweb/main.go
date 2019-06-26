package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	cyoa "github.com/zrnorth/gopher/gopher-3"
)

func main() {
	filename := flag.String("file", "story.json", "the json with the choose-your-own-adventure story")
	flag.Parse()
	fmt.Printf("Using the story %s.\n", *filename)

	f, err := os.Open(*filename)
	if err != nil {
		panic(err) // Shouldn't do this really but whatever
	}

	d := json.NewDecoder(f)
	var story cyoa.Story
	if err := d.Decode(&story); err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", story)
}
