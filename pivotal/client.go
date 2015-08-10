package pivotal

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"sync"
)

const BASE_URL string = "https://www.pivotaltracker.com/services/v5/"

var NO_PARAMS paramset = paramset{}

type paramset map[string]string

func NewClient(apiToken string) *Client {
	c := &Client{
		Token:      apiToken,
		httpClient: &http.Client{},
	}
	return c
}

type Me struct {
	ID       uint64
	Username string
	Projects []Project
}

type Project struct {
	ID   uint64 `json:"project_id"`
	Name string `json:"project_name"`
}

type Story struct {
	ID   uint64
	Name string
}

func MyStories(apiToken string) (stories []Story, err error) {
	c := NewClient(apiToken)
	me, err := c.Me()

	params := paramset{
		"filter": "mywork:" + me.Username + " -state:delivered",
	}

	results := make(chan Story)
	errors := make(chan error)

	var hadFailure bool
	var cwg = sync.WaitGroup{}

	cwg.Add(1)
	go func() {
		for err := range errors {
			log.Printf("Error: %s\n", err)
			hadFailure = true
		}

		cwg.Done()
	}()

	cwg.Add(1)
	go func() {
		for s := range results {
			stories = append(stories, s)
		}
		cwg.Done()
	}()

	var pwg = sync.WaitGroup{}
	for _, proj := range me.Projects {
		pwg.Add(1)
		go func(proj Project) {
			// log.Printf("Fetching your stories for %s\n", proj.Name)
			projStories, err := c.ProjectStories(proj.ID, params)
			if err != nil {
				errors <- err
			} else {
				for _, s := range projStories {
					results <- s
				}
			}
			pwg.Done()
		}(proj)
	}
	pwg.Wait()
	close(results)
	close(errors)
	cwg.Wait()

	if hadFailure {
		err = fmt.Errorf("ERROR: not successful getting all stories")
	}

	return
}

func (c *Client) ProjectStories(projectId uint64, params paramset) (stories []Story, err error) {
	resp, err := c.apiRequest(fmt.Sprintf("projects/%d/stories", projectId), params)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	err = parseResponse(resp, &stories)
	return
}

type Client struct {
	Token      string
	httpClient *http.Client
}

func (c *Client) Me() (m Me, err error) {
	resp, err := c.apiRequest("me", NO_PARAMS)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	err = parseResponse(resp, &m)
	return
}

func (c *Client) apiRequest(path string, params paramset) (resp *http.Response, err error) {
	req, err := http.NewRequest("GET", buildURL(path, params), nil)
	if err != nil {
		return
	}
	req.Header.Add("X-TrackerToken", c.Token)

	// log.Printf("[pivotal-tracker] GET %s\n", req.URL)
	resp, err = c.httpClient.Do(req)
	if err != nil {
		return
	}
	return
}

func parseResponse(resp *http.Response, data interface{}) (err error) {
	if resp.StatusCode != 200 {
		snip := make([]byte, 1000)
		resp.Body.Read(snip)
		err = fmt.Errorf("Error %d while making request %s\n%s\n", resp.StatusCode, resp.Request.URL, snip)
		return
	}

	d := json.NewDecoder(resp.Body)
	err = d.Decode(data)
	if err != nil {
		return
	}

	return
}

func buildURL(path string, params paramset) string {
	u := url.URL{}
	u.Path = BASE_URL + path

	query := url.Values{}
	for key, value := range params {
		query.Add(key, value)
	}

	u.RawQuery = query.Encode()

	return u.String()
}
