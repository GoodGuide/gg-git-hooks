package githooks

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"regexp"

	"github.com/rtlong/selection_ui"
)

var (
	storyTagRegexp   = regexp.MustCompile(`\[(?:(?:(?:complete[sd]?|(?:finish|fix)(?:e[sd])?)\s+)?\#\d{4,}|\#?no ?story)\]`)
	gitScissorMarker = []byte(`# ------------------------ >8 ------------------------`)
)

// Runs after supplying a commit message, is meant to check the contents of the
// message
func CommitMsg(msgFilepath string, config Config) {
	// fmt.Println(msgFilepath)
	file, err := os.OpenFile(msgFilepath, os.O_RDWR, 0664)
	if err != nil {
		log.Fatalf("Error opening %s: %s\n", msgFilepath, err)
	}
	defer file.Close()

	foundTag, err := fileContainsPivotalTag(file)
	if err != nil {
		switch err.(type) {
		case FileIsBlankError:
			// if all lines are 'blank', we want to exit 0 so git can abort due to empty
			// commit, if that setting is enabled
			os.Exit(0)

		default:
			panic(err)
		}
	}

	if !foundTag {
		if tags := promptForTag(&config); len(tags) > 0 {
			foundTag = addTagsToFile(file, tags)
		}
	}
	if !foundTag {
		fmt.Println("Missing Pivotal Tracker story tag in commit message")
		os.Exit(1)
	}
}

func promptForTag(config *Config) (tagsToAdd []string) {
	storiesCache, err := ioutil.ReadFile(config.StoriesCachePath)
	if err != nil {
		log.Printf("Unable to prompt you to tag the commit. Error while reading Tracker stories cache. Try the update-pivotal-stories command. Error was: %s\n", err)
		return
	}
	storiesBytes := bytes.Split(storiesCache, []byte{'\n'})

	var stories []string
	for _, storyBytes := range storiesBytes {
		s := bytes.Trim(storyBytes, "# ")
		if len(s) > 0 {
			stories = append(stories, string(s))
		}
	}
	stories = append(stories, "[no story]")
	selections := selection_ui.Prompt(stories, "Pivotal Tracker stories to tag this commit")
	for i, s := range selections {
		if s {
			tagsToAdd = append(tagsToAdd, stories[i])
		}
	}
	return
}

func addTagsToFile(file *os.File, tags []string) (success bool) {
	file.Seek(0, 0) // rewind to start with
	r, err := ioutil.ReadAll(io.Reader(file))
	if err != nil {
		log.Println("ERROR: %s", err)
		return
	}
	loc := bytes.Index(r, gitScissorMarker)
	if loc < 0 {
		if _, err := file.Seek(0, 2); err != nil { // seek to end of file
			log.Println("ERROR: %s", err)
			return
		}
	} else {
		if _, err := file.Seek(int64(loc), 0); err != nil { // seek to scissor line then truncate the file here
			log.Println("ERROR: %s", err)
			return
		}
		if err := file.Truncate(int64(loc + 1)); err != nil {
			log.Println("ERROR: %s", err)
			return
		}
	}
	fmt.Fprintln(file, "") // blank line first to ensure it's not stuck to the summary
	for _, tag := range tags {
		fmt.Fprintln(file, tag)
	}
	return true
}

type FileIsBlankError error

func fileContainsPivotalTag(file io.Reader) (bool, error) {
	var fileIsBlank bool = true

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		} else if line[0] == '#' {
			if bytes.Index(line, gitScissorMarker) >= 0 {
				// if the git-commit -v option is used, there is a diff block below the commit
				// message template, and we need to ignore that for the purpose of this
				// test
				break
			} else {
				continue
			}
		}
		fileIsBlank = false

		if storyTagRegexp.Match(line) {
			return true, nil
		}
	}

	if err := scanner.Err(); err != nil {
		return false, err
	}

	if fileIsBlank {
		return false, FileIsBlankError(fmt.Errorf("File didn't contain any non-blank lines"))
	}

	return false, nil
}
