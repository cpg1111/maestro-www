package data

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	DaemonAddr     string
	PollIntervalMS int
	httpClient     *http.Client
}

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
	return fmt.Sprintf("%s&branch=%s", c.projectEndPoint(), branchEsc)
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
			time.Sleep(c.PollIntervalMS * time.Millisecond)
		}
	}()
	go func() {
		shouldRun = <-interrupt
	}()
	return resChan, errChan
}

func (c *Client) GetAll() (chan []byte, chan error) {
	endpoint := c.stateEndPoint()
	return c.get(endpoint)
}

func (c *Client) GetOne(project, branch string) (chan []byte, chan error) {
	endpoint := c.singleEndPoint(project, branch)
	return c.get(endpoint)
}

func (c *Client) WatchOne(project, branch string, interrupt chan bool) (chan []byte, chan error) {
	endpoint := c.singleEndPoint(project, branch)
	return c.watch(endpoint)
}

func (c *Client) GetProject(project string) (chan []byte, chan error) {
	endpoint := c.projectEndPoint(project)
	return c.get(endpoint)
}

func (c *Client) WatchProject(project string) (chan []byte, chan error) {
	endpoint := c.projectEndPoint(project)
	return c.watch(endpoint)
}
