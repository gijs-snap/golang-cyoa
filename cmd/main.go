package main

import (
	"flag"
	"fmt"
	"os"
	"github.com/gijs-snap/golang-cyoa"
	"net/http"
	"log"
)

func main() {
	port := flag.Int("port", 3000, "The port the web app will run on")
	// Retrieve file containing our story
	file := flag.String("file", "story.json", "JSON File")
	flag.Parse()
	fmt.Printf("Using the story in %s", *file)

	f, err := os.Open(*file)
	if err != nil {
		panic(err)
	}

	story, err := cyoa.JsonStory(f)
	if err != nil {
		panic(err)
	}

	h := cyoa.NewHandler(story)
	fmt.Printf("Starting server on port: %d", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), h))
}
