package data

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

// Client is a client to maestrod
type Client struct {
	DaemonAddr     string
	PollIntervalMS int
	httpClient     *http.Client
}

// NewClient returns a pointer to an instance of Client
func NewClient(daemonAddr string, pollIntrvl int) *Client {
	return &Client{
		DaemonAddr:     daemonAddr,
		PollIntervalMS: pollIntrvl,
		httpClient:     &http.Client{},
	}
}

func (c *Client) stateEndPoint() string {
	return fmt.Sprintf("%s/state", c.DaemonAddr)
}

func (c *Client) projectEndPoint(project string) string {
	projectEsc := url.QueryEscape(project)
	return fmt.Sprintf("%s?project=%s", c.stateEndPoint(), projectEsc)
}

func (c *Client) singleEndPoint(project, branch string) string {
	branchEsc := url.QueryEscape(branch)
	return fmt.Sprintf("%s&branch=%s", c.projectEndPoint(project), branchEsc)
}

func (c *Client) get(endpoint string) (chan []byte, chan error) {
	resChan := make(chan []byte)
	errChan := make(chan error)
	go func() {
		resp, getErr := c.httpClient.Get(endpoint)
		defer resp.Body.Close()
		if getErr != nil {
			errChan <- getErr
			return
		}
		res, readErr := ioutil.ReadAll(resp.Body)
		resChan <- res
		errChan <- readErr
	}()
	return resChan, errChan
}

func (c *Client) watch(endpoint string, interrupt chan bool) (chan []byte, chan error) {
	resChan := make(chan []byte)
	errChan := make(chan error)
	shouldRun := true
	go func() {
		for shouldRun {
			res, err := c.get(endpoint)
			resChan <- <-res
			errChan <- <-err
			time.Sleep(time.Millisecond * (time.Duration)(c.PollIntervalMS))
		}
	}()
	go func() {
		shouldRun = <-interrupt
	}()
	return resChan, errChan
}

// GetAll gets all states of maestrod
func (c *Client) GetAll() (chan []byte, chan error) {
	endpoint := c.stateEndPoint()
	return c.get(endpoint)
}

// GetOne gets the state of a particular branch of a project
func (c *Client) GetOne(project, branch string) (chan []byte, chan error) {
	endpoint := c.singleEndPoint(project, branch)
	return c.get(endpoint)
}

// WatchOne will watch a current build in maestrod
func (c *Client) WatchOne(project, branch string, interrupt chan bool) (chan []byte, chan error) {
	endpoint := c.singleEndPoint(project, branch)
	return c.watch(endpoint, interrupt)
}

// GetProject will get the state of a project in maestrod
func (c *Client) GetProject(project string) (chan []byte, chan error) {
	endpoint := c.projectEndPoint(project)
	return c.get(endpoint)
}

// WatchProject will watch for builds within a project
func (c *Client) WatchProject(project string, interrupt chan bool) (chan []byte, chan error) {
	endpoint := c.projectEndPoint(project)
	return c.watch(endpoint, interrupt)
}
