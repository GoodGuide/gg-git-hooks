package githooks

import (
	"fmt"
	"log"
	"os"

	"github.com/GoodGuide/goodguide-git-hooks/pivotal"
)

func UpdatePivotalStories(config Config) {
	file, err := os.Create(config.StoriesCachePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	stories, err := pivotal.MyStories(config.APIToken)
	if err != nil {
		log.Fatal(err)
	}

	for _, story := range stories {
		fmt.Fprintf(file, "#[#%d] %s\n", story.ID, story.Name)
	}

	log.Printf("Wrote stories cache to %s\n", file.Name())
}
