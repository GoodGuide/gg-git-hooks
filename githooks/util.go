package githooks

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	"github.com/GoodGuide/goodguide-git-hooks/pivotal"
)

func indentLines(b []byte, level int) (out []byte) {
	lineSep := []byte("\n")
	indent := bytes.Repeat([]byte("\t"), level)

	buf := bytes.NewBuffer(out)

	for _, line := range bytes.Split(b, lineSep) {
		if len(line) == 0 {
			continue
		}
		buf.Write(indent)
		buf.Write(line)
		buf.Write(lineSep)
	}

	return buf.Bytes()
}

func loadStoriesFromCache(filepath string) (stories []pivotal.Story, err error) {
	var f *os.File
	if f, err = os.Open(filepath); err != nil {
		return
	}
	defer f.Close()

	d := json.NewDecoder(f)
	err = d.Decode(&stories)
	return
}
func writeStoriesToCache(filepath string, stories []pivotal.Story) (err error) {
	var f *os.File
	if f, err = os.Create(filepath); err != nil {
		return
	}
	defer f.Close()

	d := json.NewEncoder(f)
	err = d.Encode(stories)
	return
}
func formatStoriesAsStrings(stories []pivotal.Story) (strings []string) {
	strings = make([]string, len(stories))
	for i, s := range stories {
		strings[i] = fmt.Sprintf("[#%d] %s", s.ID, s.Name)
	}
	return
}
