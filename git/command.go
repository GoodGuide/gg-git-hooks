package git

import "os/exec"

func Command(args ...string) ([]byte, error) {
	cmd := exec.Command("git", args...)
	out, err := cmd.CombinedOutput()
	return out, err
}
