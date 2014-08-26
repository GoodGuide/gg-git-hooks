package git

import (
	"fmt"
	"path/filepath"
	"strings"
)

func GitDir() (string, error) {
	out, err := Command("rev-parse", "--git-dir")
	if err != nil {
		err = fmt.Errorf("%s\n", out)
		return "", err
	}
	gitDir := strings.TrimSpace(string(out))

	gitDir, err = filepath.Abs(gitDir)
	return gitDir, err
}
