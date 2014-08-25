package main

import (
	"github.com/GoodGuide/goodguide-git-hooks/githooks"
	kingpin "gopkg.in/alecthomas/kingpin.v1"
)

// Installs small shell scripts for all the git hooks in the .git/hooks
// directory for the git repo in which this command is run
func install() {
}

var (
	debug *bool
)

func main() {
	debug = kingpin.Flag("debug", "Show extra info").Bool()

	kingpin.Command("install", "Install scripts at .git/hooks/* for each git-hook provided by this tool")

	cmdCommitMsg := kingpin.Command("commit-msg", "Checks the commit message for PivotalTracker story ID, bad whitespace, syntax, etc.")
	messageFilepath := cmdCommitMsg.Arg("message_path", "Path to the file that holds the proposed commit log message").
		Required().
		ExistingFile()

	cmdPrepareCommitMsg := kingpin.Command("prepare-commit-msg", "Augment the default commit message template with commented-out PivotalTracker Story IDs to make it easy to tag commits")
	cmdPrepareCommitMsg.Arg("message_path", "Path to the file which will be sent to the editor and ultimately become the commit log message").
		Required().
		ExistingFileVar(messageFilepath)

	messageSource := cmdPrepareCommitMsg.Arg("source", "Source of the commit message going into this hook").
		Enum("message", "merge", "commit", "squash", "template")

	messageSourceCommit := cmdPrepareCommitMsg.Arg("commit_sha", "If source is 'commit', this is the SHA1 of the source commit").
		String()

	kingpin.Command("pre-commit", "Verifies the files about to be committed follow certain guidelines regarding e.g. whitespace, syntax, etc.")

	// no-ops:
	kingpin.Command("applypatch-msg", "no-op")
	kingpin.Command("post-update", "no-op")
	kingpin.Command("pre-applypatch", "no-op")
	kingpin.Command("pre-push", "no-op")
	kingpin.Command("pre-rebase", "no-op")
	kingpin.Command("update", "no-op")

	switch kingpin.Parse() {
	case "install":
		install()

	case "commit-msg":
		githooks.CommitMsg(*messageFilepath)

	case "prepare-commit-msg":
		githooks.PrepareCommitMsg(*messageFilepath, *messageSource, *messageSourceCommit)

	case "pre-commit":
		githooks.PreCommit()
	}
}
