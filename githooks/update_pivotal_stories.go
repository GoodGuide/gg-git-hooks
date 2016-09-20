package githooks

import (
	"github.com/goodguide/goodguide-git-hooks/pivotal"
	"log"
	"sort"
)

func UpdatePivotalStories(config Config) (stories []pivotal.Story) {
	stories, err := updatePivotalStoriesCache(config)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Wrote %d stories to cache file at %s\n", len(stories), config.StoriesCachePath)
	return stories
}

type ByDate []pivotal.Story

func (a ByDate) Len() int           { return len(a) }
func (a ByDate) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByDate) Less(i, j int) bool { return a[i].UpdatedAt.After(a[j].UpdatedAt) }

func updatePivotalStoriesCache(config Config) (stories []pivotal.Story, err error) {
	stories, err = pivotal.MyStories(config.APIToken)
	if err != nil {
		return
	}
	sort.Sort(ByDate(stories))
	// fmt.Printf("%#v\n", stories)
	if err = writeStoriesToCache(config.StoriesCachePath, stories); err != nil {
		return
	}
	return stories, nil
}
