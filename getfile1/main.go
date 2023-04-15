package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"go.starlark.net/starlark"
)

var (
	sourceUrl string
	destFile  string
)

func main() {
	// 0. register Starlark builtin
	registrar := starlark.StringDict{"config": starlark.NewBuiltin("config", configFn)}

	// 1. Exec script file
	_, err := starlark.ExecFile(&starlark.Thread{}, "getfile.star", nil, registrar)
	if err != nil {
		log.Fatalf("Starlark Exec: %s", err)
	}

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

// configFn implements the Starlark builtin for script function config:
// config(source_url="<source>", dest_file="<file>")
func configFn(_ *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	if err := starlark.UnpackArgs(
		"config", args, kwargs,
		"source_url", &sourceUrl,
		"dest_file", &destFile,
	); err != nil {
		return starlark.None, fmt.Errorf("config: %s", err)
	}
	return starlark.None, nil
}
