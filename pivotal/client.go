package pivotal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
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
	params := make(paramset)
	// params["filter"] = "state:started mywork:" + me.Username
	params["filter"] = "mywork:" + me.Username

	// TODO: make this concurrent
	for _, proj := range me.Projects {
		projStories, err := c.ProjectStories(proj.ID, params)
		if err != nil {
			return nil, err
		}
		for _, s := range projStories {
			stories = append(stories, s)
		}
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
