package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
)

func confirm(prompt string) (bool, error) {
	fmt.Printf("%s (y/N): ", prompt)
	bio := bufio.NewReader(os.Stdin)
	line, _, err := bio.ReadLine()
	if err != nil {
		return false, err
	}
	if len(line) > 0 && bytes.ToUpper(line)[0] == 'Y' {
		return true, nil
	}
	return false, nil
}

func writeHookShim(w io.Writer, hookName string) (err error) {
	_, err = fmt.Fprintln(w, "#!/bin/bash")
	if err != nil {
		return
	}
	_, err = fmt.Fprint(w, "set -e\n\n")
	if err != nil {
		return
	}
	_, err = fmt.Fprintf(w, "exec goodguide-git-hooks %s $@\n", hookName)
	return
}
