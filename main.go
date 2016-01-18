package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/goodguide/goodguide-git-hooks/git"
	"github.com/goodguide/goodguide-git-hooks/githooks"
	kingpin "gopkg.in/alecthomas/kingpin.v1"
)

var (
	config              githooks.Config
	messageFilepath     *string = new(string)
	messageSourceCommit *string = new(string)
	messageSource       *string = new(string)
	clobber             *bool   = new(bool)
	noclobber           *bool   = new(bool)
	HOOKS                       = [8]string{
		"commit-msg",
		"prepare-commit-msg",
		"pre-commit",
		"applypatch-msg",
		"post-update",
		"pre-applypatch",
		"pre-push",
		"pre-rebase",
	}
)

// These are set at build-time via goxc
var (
	VERSION    string
	BUILD_DATE string
)

func Version() string {
	return fmt.Sprintf("%s - built %s", VERSION, BUILD_DATE)
}

// Installs small shell scripts for all the git hooks in the .git/hooks directory for the git repo in which this command is run.
// TODO: Make this smarter about offering to overwrite an existing file when it already has the exact contents we want to write
func InstallHookShims(hooksDir string, hooks []string) {
	for _, hook := range hooks {
		var confirmed bool = *clobber

		hookPath := filepath.Join(hooksDir, hook)

		stat, err := os.Stat(hookPath)
		if err != nil {
			if os.IsNotExist(err) {
				confirmed = true
			} else {
				log.Printf("[%s] Error while installing: %s\n", hook, err)
				continue
			}
		}

		if !confirmed {
			if !stat.Mode().IsRegular() {
				log.Printf("[%s] ERROR: File already exists but is not a regular file!\n", hook)
				continue
			}
			if *noclobber {
				confirmed = false
			} else if *clobber {
				confirmed = true
			} else {
				confirmed, err = confirm(fmt.Sprintf("[%s] File already exists. Overwrite?", hook))
				if err != nil {
					panic(err)
				}
			}
		}

		if confirmed {
			log.Printf("[%s] installing shim\n", hook)
			file, err := os.Create(hookPath)
			if err != nil {
				log.Printf("[%s] Error while opening file: %s\n", hook, err)
				continue
			}
			defer file.Close()

			if err := writeHookShim(file, hook); err != nil {
				log.Printf("[%s] Error while writing file: %s\n", hook, err)
			}
			if err := file.Chmod(0755); err != nil {
				log.Printf("[%s] Error while setting file permissions: %s\n", hook, err)
			}
		}
	}
}

func SelfUpdate() {
	fmt.Println("The self-update feature doesn't exist yet. Please check github for latest binary release, or try")
	fmt.Println("  go get -u github.com/goodguide/goodguide-git-hooks")
}

func initKingpin() {
	var cmd *kingpin.CmdClause

	kingpin.Version(Version())

	cmd = kingpin.Command("install", "Install scripts at .git/hooks/* for each git-hook provided by this tool")
	cmd.Flag("force", "Don't ask to overwrite target file, always clobber").BoolVar(clobber)
	cmd.Flag("noclobber", "Don't ask to overwrite target file, always skip").BoolVar(noclobber)

	kingpin.Command("self-update", "Check for updates of goodguide-git-hooks and download the newer version if available")

	kingpin.Command("update-pivotal-stories", "Update cache of pivotal stories manually")

	// git hooks commands:
	kingpin.Command("pre-commit", "Verifies the files about to be committed follow certain guidelines regarding e.g. whitespace, syntax, etc.")

	cmd = kingpin.Command("prepare-commit-msg", "Augment the default commit message template with commented-out PivotalTracker Story IDs to make it easy to tag commits")
	cmd.Arg("message_path", "Path to the file which will be sent to the editor and ultimately become the commit log message").
		Required().
		ExistingFileVar(messageFilepath)
	cmd.Arg("source", "Source of the commit message going into this hook").
		EnumVar(&messageSource, "message", "merge", "commit", "squash", "template")
	cmd.Arg("commit_sha", "If source is 'commit', this is the SHA1 of the source commit").
		StringVar(messageSourceCommit)

	cmd = kingpin.Command("commit-msg", "Checks the commit message for PivotalTracker story ID, bad whitespace, syntax, etc.")
	cmd.Arg("message_path", "Path to the file that holds the proposed commit log message").
		Required().
		ExistingFileVar(messageFilepath)

	// no-ops:
	cmd = kingpin.Command("applypatch-msg", "no-op")
	cmd.Arg("args", "").Strings()
	cmd = kingpin.Command("post-update", "no-op")
	cmd.Arg("args", "").Strings()
	cmd = kingpin.Command("pre-applypatch", "no-op")
	cmd.Arg("args", "").Strings()
	cmd = kingpin.Command("pre-push", "no-op")
	cmd.Arg("args", "").Strings()
	cmd = kingpin.Command("pre-rebase", "no-op")
	cmd.Arg("args", "").Strings()
}

func main() {
	initKingpin()

	config.APIToken = GetAPIToken()
	config.StoriesCachePath = PivotalStoriesCacheFilePath()

	switch kingpin.Parse() {
	case "install":
		gitDir, err := git.GitDir()
		if err != nil {
			log.Fatal("Error while installing:\n", err)
		}
		hooksDir := filepath.Join(gitDir, "hooks")
		if err := os.MkdirAll(hooksDir, 0775); err != nil {
			log.Fatalf("Error while creating hooks directory: %s\n", err)
		}

		InstallHookShims(hooksDir, HOOKS[:])

	case "self-update":
		SelfUpdate()

	case "update-pivotal-stories":
		githooks.UpdatePivotalStories(config)

	case "pre-commit":
		githooks.PreCommit()

	case "prepare-commit-msg":
		config.StoriesCachePath = PivotalStoriesCacheFilePath()

		githooks.PrepareCommitMsg(*messageFilepath, *messageSource, *messageSourceCommit, config)

	case "commit-msg":
		config.StoriesCachePath = PivotalStoriesCacheFilePath()

		githooks.CommitMsg(*messageFilepath, config)
	}
}
