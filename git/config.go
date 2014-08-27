package git

import "fmt"

func ConfigGetString(key string) (string, error) {
	out, err := Command("config", "--get", key)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s", out), nil
}
