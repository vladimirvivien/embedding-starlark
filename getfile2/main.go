package main

import (
	"bufio"
	"fmt"
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
	script, err := starlark.ExecFile(&starlark.Thread{}, "getfile.star", nil, registrar)
	if err != nil {
		log.Fatalf("Starlark Exec: %s", err)
	}

	// 2. retrieve a text line processor from script
	var procLineFn starlark.Value
	if procLine := script["proc_line"]; procLine != nil && procLine.Type() == "function" {
		procLineFn = procLine
	}

	// 3. download resource
	rsp, err := http.Get(sourceUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer rsp.Body.Close()

	// 4. process and save content
	file, err := os.Create(destFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scnr := bufio.NewScanner(bufio.NewReader(rsp.Body))
	scnr.Split(bufio.ScanLines)
	for scnr.Scan() {
		if scnr.Err() != nil {
			log.Printf("text scanning error: %s", err)
			continue
		}
		text := scnr.Text()
		if procLineFn != nil {
			// call starlark function proc_line, expect text back
			result, err := starlark.Call(&starlark.Thread{}, procLineFn, starlark.Tuple{starlark.String(text)}, nil)
			if err != nil {
				log.Printf("starlark text processing failed: %s", err)
			}
			if str, ok := starlark.AsString(result); ok {
				text = str
			}
		}
		file.WriteString(text)
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
