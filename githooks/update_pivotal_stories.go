package githooks

import (
	"log"

	"github.com/GoodGuide/goodguide-git-hooks/pivotal"
)

func UpdatePivotalStories(config Config) (stories []pivotal.Story) {
	stories, err := updatePivotalStoriesCache(config)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Wrote %d stories to cache file at %s\n", len(stories), config.StoriesCachePath)
	return stories
}

func updatePivotalStoriesCache(config Config) (stories []pivotal.Story, err error) {
	stories, err = pivotal.MyStories(config.APIToken)
	if err != nil {
		return
	}
	if err = writeStoriesToCache(config.StoriesCachePath, stories); err != nil {
		return
	}
	return stories, nil
}
