package main

import (
	"log"
	"os"
	"os/user"
	"path"

	"github.com/goodguide/goodguide-git-hooks/git"
)

func PivotalStoriesCacheFilePath() string {
	if custom := os.Getenv("GIT_HOOKS_CACHE"); custom != "" {
		return custom
	}
	usr, _ := user.Current()
	dir := usr.HomeDir
	return path.Join(dir, ".gg-git-hooks-cache.json")
}

func GetAPIToken() string {
	if str := os.Getenv("PIVOTAL_API_TOKEN"); str != "" {
		return str
	}
	str, err := git.ConfigGetString("pivotal.api-token")
	if err != nil {
		log.Fatal("Can't find Pivotal API Token. Set it in git-config at the pivotal.api-token key.", err)
	}
	return str
}
