package main

import (
	"fmt"
	"os"

	"github.com/GoodGuide/goodguide-git-hooks/githooks"
)

// Prints a help message instructing with the available commands
func printHelp() {
	fmt.Println("Don't know what you mean by:", os.Args)
	fmt.Println("Usage:", os.Args[0], "COMMAND [ARG...]")
	fmt.Println("  Commands: commit-msg, prepare-commit-msg, pre-commit")
	os.Exit(0)
}

func main() {
	if len(os.Args) < 2 {
		printHelp()
	}
	switch os.Args[1] {
	case "commit-msg":
		if len(os.Args) < 3 {
			printHelp()
		}
		githooks.CommitMsg(os.Args[2])

	case "prepare-commit-msg":
		var commitSha string
		var source string
		if len(os.Args) < 3 {
			printHelp()
		}
		if len(os.Args) >= 4 {
			source = os.Args[3]
		}
		if len(os.Args) >= 5 {
			commitSha = os.Args[4]
		}
		githooks.PrepareCommitMsg(os.Args[2], source, commitSha)

	case "pre-commit":
		githooks.PreCommit()

	default:
		printHelp()
	}
}
