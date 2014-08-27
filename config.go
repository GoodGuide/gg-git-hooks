package main

import (
	"log"
	"os/user"
	"path"

	"github.com/GoodGuide/goodguide-git-hooks/git"
)

func PivotalStoriesCacheFilePath() string {
	usr, _ := user.Current()
	dir := usr.HomeDir
	return path.Join(dir, ".gg-git-hooks-cache")
}

func GetAPIToken() string {
	str, err := git.ConfigGetString("pivotal.api-token")
	if err != nil {
		log.Fatal("Can't find Pivotal API Token. Set it in git-config at the pivotal.api-token key.", err)
	}
	return str
}
