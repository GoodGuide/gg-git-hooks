package githooks

import (
	"bytes"
	"fmt"
	"os"
	"sync"

	"github.com/GoodGuide/goodguide-git-hooks/git"
)

const (
	RESULT_OKAY int = 0
	RESULT_SKIP int = 1
	RESULT_FAIL int = 2
)

// Assemble list of files to check. Ask git which files have changes and weren't deleted
func changedFiles(results chan<- *fileToCheck) {
	defer close(results)

	out, err := git.Command("diff-index", "HEAD", "--cached", "--name-only", "-z", "--diff-filter=ACRMT")
	if err != nil {
		fmt.Printf("ERROR: git-diff-index: %s\n%s\n", err, out)
		os.Exit(2)
	}

	for _, chunk := range bytes.Split(out, []byte{'\x00'}) {
		if len(chunk) == 0 {
			continue
		}

		results <- &fileToCheck{Path: string(chunk)}
	}
}

// Runs while building the commit snapshot, can check the files that will be
// changed for syntax, whitespace, etc.
func PreCommit() {
	filesToCheck := make(chan *fileToCheck)
	go changedFiles(filesToCheck)

	var wg sync.WaitGroup
	results := make(chan *fileCheckResult)

	for file := range filesToCheck {
		wg.Add(1)
		go checkFile(file, results, &wg)
	}
	go func() {
		wg.Wait()
		close(results)
	}()

	var fail bool

	for result := range results {
		switch result.Result {
		case RESULT_OKAY:
			// fmt.Println("OKAY:", result.File.Path)
		case RESULT_SKIP:
			// fmt.Println("SKIP:", result.File.Path, "-", result.Error)
		case RESULT_FAIL:
			fmt.Println("FAIL:", result.File.Path, "-", result.Error)
			fail = true
		}
	}

	if fail {
		os.Exit(1)
	}
}

type fileCheckResult struct {
	File   *fileToCheck
	Result int
	Error  error
}

func checkFile(file *fileToCheck, results chan<- *fileCheckResult, wg *sync.WaitGroup) {
	defer wg.Done()

	result := fileCheckResult{File: file}
	defer func() { results <- &result }()

	if err := file.PreCheck(); err != nil {
		result.Result = RESULT_SKIP
		result.Error = err
		return
	}

	if err := file.Check(); err != nil {
		result.Result = RESULT_FAIL
		result.Error = err
		return
	}
}
