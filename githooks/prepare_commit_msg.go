package githooks

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
)

func WritePivotalStories(w io.Writer, config *Config) {
	w.Write([]byte("\n# Uncomment one of your active stories, below:\n"))

	stories, err := ioutil.ReadFile(config.StoriesCachePath)
	if err == nil {
		w.Write(stories)
	} else {
		w.Write([]byte("# There was a problem getting your Tracker Stories from ~/.gg-git-hooks-cache\n# To (re)create/update the file:\n#\n#   goodguide-git-hooks update-stories\n#\n"))
	}

	w.Write([]byte("#[no story]\n\n\n"))
}

// Runs just before opening the editor to get a message from the user. In this
// case, it fetches pivotal tracker stories and modifies the message template to
// include the story ids as commented-out lines
func PrepareCommitMsg(msgFilepath string, source string, commitSha string, config Config) {
	fmt.Println("prepare-commit-msg", msgFilepath, source, commitSha)

	if source == "merge" {
		return
	}

	originalFile, err := os.OpenFile(msgFilepath, os.O_RDWR, 0664)
	if err != nil {
		log.Fatal(err)
	}
	defer originalFile.Close()
	originalMsg := bufio.NewReader(originalFile)

	newMsg, err := ioutil.TempFile("", "goodguide-git-hooks")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		newMsg.Close()
		os.Remove(newMsg.Name())
	}()

	var insertedStories bool

	for {
		line, err := originalMsg.ReadBytes('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("ERR", err)
		}

		if !insertedStories && line[0] == '#' {
			WritePivotalStories(newMsg, &config)
			insertedStories = true
		}

		newMsg.Write(line)
	}

	originalFile.Truncate(0)
	originalFile.Seek(0, 0)

	if o, err := newMsg.Seek(0, 0); err != nil {
		log.Fatal(err, o)
	}

	if _, err := io.Copy(originalFile, newMsg); err != nil {
		log.Fatal(err)
	}
}
