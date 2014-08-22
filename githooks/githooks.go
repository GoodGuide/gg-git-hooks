package githooks

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/GoodGuide/goodguide-git-hooks/git"
)

// Runs after supplying a commit message, is meant to check the contents of the
// message
func CommitMsg(msgFilepath string) {
	// fmt.Println("commit-msg", msgFilepath)
}

// Runs just before opening the editor to get a message from the user. In this
// case, it fetches pivotal tracker stories and modifies the message template to
// include the story ids as commented-out lines
func PrepareCommitMsg(msgFilepath string, source string, commitSha string) {
	// fmt.Println("prepare-commit-msg", msgFilepath, source, commitSha)
}

type FileToCheck struct {
	Path string
	Stat os.FileInfo
}

// Runs while building the commit snapshot, can check the files that will be
// changed for syntax, whitespace, etc.
func PreCommit() {
	// fmt.Println("pre-commit")

	// Assemble list of files to check. Ask git which files have changes and weren't deleted
	out, err := git.Command("diff-index", "HEAD", "--cached", "--name-only", "-z", "--diff-filter=ACRMT")
	if err != nil {
		fmt.Println("ERROR: git-diff-index:", err)
		fmt.Printf("%s\n", out)
		os.Exit(1)
	}

	results := make(chan FileToCheck)
	go func() {
		var wg sync.WaitGroup
		for _, chunk := range bytes.Split(out, []byte{'\x00'}) {
			if len(chunk) == 0 {
				continue
			}
			filepath := string(chunk)
			wg.Add(1)
			go lookupFileStat(filepath, results, &wg)
		}
		wg.Wait()
		close(results)
	}()

	problematicFiles := make(chan FileToCheck)
	go func() {
		var wg sync.WaitGroup
		for fileToCheck := range results {
			wg.Add(1)
			go checkFile(fileToCheck, problematicFiles, &wg)
		}
		wg.Wait()
		close(problematicFiles)
	}()

	var problem bool
	for file := range problematicFiles {
		log.Println("Problem with file:", file.Path)
		problem = true
	}

	if problem {
		// os.Exit(1)
	}
}

func lookupFileStat(filepath string, result chan<- FileToCheck, wg *sync.WaitGroup) {
	defer wg.Done()

	file, err := os.Lstat(filepath)
	if err != nil {
		log.Printf("ERR: %s", err)
	}

	if file.Mode().IsRegular() {
		result <- FileToCheck{filepath, file}
	}
}

func checkFile(file FileToCheck, problematicFiles chan<- FileToCheck, wg *sync.WaitGroup) {
	defer wg.Done()

	if file.Stat.Size() == 0 {
		log.Println("WARN:", file.Path, "is an empty file")
		return
	}

	// Check mime type before we go further
	cmd := exec.Command("file", "--mime-type", "--brief", file.Path)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Println("WARN: Couldn't determine file type for", file.Path, "; skipping: ", err)
		return
	}

	mimetype := string(out[:len(out)-1])
	// fmt.Println(file.Path, mimetype)
	if strings.Split(mimetype, "/")[0] != "text" {
		log.Println("WARN: Skipping", file.Path, "for file contents check")
		return
	}

	problematicFiles <- file
}
