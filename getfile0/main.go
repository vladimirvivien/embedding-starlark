package main

import (
	"io"
	"log"
	"net/http"
	"os"
)

// Download/save text file

var (
	sourceUrl string
	destFile  string
)

func main() {
	sourceUrl = "https://www.gutenberg.org/files/408/408-0.txt"
	destFile = "./soul-black-folks.txt"

	// 1. download resource
	rsp, err := http.Get(sourceUrl)
	if err != nil {
		log.Fatal(err)
	}

	// 2. write resource to destination
	var content []byte
	if content, err = io.ReadAll(rsp.Body); err != nil {
		log.Fatal(err)
	}
	defer rsp.Body.Close()

	if err := os.WriteFile(destFile, content, 0644); err != nil {
		log.Fatal(err)
	}
}
