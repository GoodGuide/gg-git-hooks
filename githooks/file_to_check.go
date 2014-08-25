package githooks

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"regexp"

	"github.com/GoodGuide/goodguide-git-hooks/git"
)

type fileToCheck struct {
	Path string
	stat os.FileInfo
}

var (
	uncommitablesRegexp = regexp.MustCompile(
		`\A\x1B\[32m\+.+(binding.pry|debugger|(?i)rbc)`,
	)
)

func (f *fileToCheck) getStat() error {
	if f.stat == nil {
		fstat, err := os.Lstat(f.Path)
		if err != nil {
			return err
		}
		f.stat = fstat
	}
	return nil
}

func (f *fileToCheck) VerifyType() error {
	if err := f.getStat(); err != nil {
		return err
	}

	if !f.stat.Mode().IsRegular() {
		return fmt.Errorf("is not a regular file")
	}

	return nil
}

func (f *fileToCheck) VerifyIsNotEmpty() error {
	if err := f.getStat(); err != nil {
		return err
	}

	if f.stat.Size() == 0 {
		return fmt.Errorf("is empty")
	}

	return nil
}

func (f *fileToCheck) VerifyIsText() error {
	cmd := exec.Command("file", "--mime-type", "--brief", f.Path)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("Couldn't determine file type: %s\n%s", err, out)
	}

	if t := bytes.Split(out[:len(out)-1], []byte("/"))[0]; string(t) != "text" {
		return fmt.Errorf("doesn't appear to be textual")
	}

	return nil
}

func (f *fileToCheck) VerifyEOFNewline() error {
	if err := f.getStat(); err != nil {
		return err
	}

	file, err := os.Open(f.Path)
	if err != nil {
		return err
	}
	defer file.Close()

	var b = make([]byte, 1)
	if _, err := file.ReadAt(b, f.stat.Size()-1); err != nil {
		return err
	}

	if b[0] != byte('\n') {
		return fmt.Errorf("missing newline character at EOF")
	}

	return nil
}

// uses the --check option in git-diff to check for invalid whitespace
func (f *fileToCheck) GitDiffCheck() error {
	out, err := git.Command("diff", "--check", "--cached", "--color=always", "--", f.Path)
	if err != nil {
		return fmt.Errorf("git-diff --check failed:\n%s", indentLines(out, 1))
	}

	return nil
}

// verifies certain keywords aren't being added
func (f *fileToCheck) VerifyNoUncommittables() error {
	out, err := git.Command("diff", "--exit-code", "--cached", "--color=always", "--", f.Path)
	if err == nil {
		return nil
	}

	for _, line := range bytes.Split(out, []byte("\n")) {
		if uncommitablesRegexp.Match(line) {
			return fmt.Errorf("forbidden string found:\n%s", indentLines(out, 1))
		}
	}

	return nil
}

func (file *fileToCheck) PreCheck() error {
	if err := file.VerifyType(); err != nil {
		return err
	}

	if err := file.VerifyIsNotEmpty(); err != nil {
		return err
	}

	if err := file.VerifyIsText(); err != nil {
		return err
	}

	return nil
}

func (file *fileToCheck) Check() error {
	if err := file.VerifyEOFNewline(); err != nil {
		return err
	}

	if err := file.GitDiffCheck(); err != nil {
		return err
	}

	if err := file.VerifyNoUncommittables(); err != nil {
		return err
	}

	return nil
}
