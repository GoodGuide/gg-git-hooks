package githooks

import "fmt"

// Runs after supplying a commit message, is meant to check the contents of the
// message
func CommitMsg(msgFilepath string) {
	fmt.Println("commit-msg", msgFilepath)
}
