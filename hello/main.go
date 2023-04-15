package main

import (
	"log"

	"go.starlark.net/starlark"
)

func main() {
	_, err := starlark.ExecFile(&starlark.Thread{}, "hello.star", nil, nil)
	if err != nil {
		log.Fatalf("Starlark Exec: %s", err)
	}
}
