package githooks

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
)

var (
	storyTagRegexp      = regexp.MustCompile(`\[(?:(?:(?:complete[sd]?|(?:finish|fix)(?:e[sd])?)\s+)?\#\d{4,}|\#?no ?story)\]`)
	gitDiffMarkerRegexp = regexp.MustCompile(`-{10,}\s*\>8\s*-{10,}`)
)

// Runs after supplying a commit message, is meant to check the contents of the
// message
func CommitMsg(msgFilepath string) {
	file, err := os.Open(msgFilepath)
	if err != nil {
		log.Fatalf("Error opening %s: %s\n", msgFilepath, err)
	}
	defer file.Close()

	var foundTag, atLeastOneLineWasntBlank bool

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		} else if line[0] == '#' {
			if gitDiffMarkerRegexp.Match(line) {
				// if the git-commit -v option is used, there is a diff block below the commit
				// message template, and we need to ignore that for the purpose of this
				// test
				break
			} else {
				continue
			}
		}
		// if all lines are 'blank', we want to exit 0 so git can abort due to empty
		// commit, if that setting is enabled
		atLeastOneLineWasntBlank = true

		if storyTagRegexp.Match(line) {
			foundTag = true
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal("Error while reading file:", err)
	}

	if atLeastOneLineWasntBlank && !foundTag {
		fmt.Println("Missing Pivotal Tracker story tag in commit message")
		os.Exit(1)
	}
}
