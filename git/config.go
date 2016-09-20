package git

import (
	"fmt"
	"strings"
)

func ConfigGetString(key string) (string, error) {
	out, err := Command("config", "--get", key)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(fmt.Sprintf("%s", out)), nil
}
