package processing

import (
	"log"
	"net/http"
	"time"
	"fmt"
)

type HttpResponse struct {
	Url      string
	Response *http.Response
	Err      error
}

type BundledHttpGet struct {
	Urls          []string
	HttpResponses []*HttpResponse
	client        *http.Client
}

func (c *BundledHttpGet) add(url string) {
	c.Urls = append(c.Urls, url)
}

func (c *BundledHttpGet) Run() {
	ch := make(chan *HttpResponse)
	for _, url := range c.Urls {
		go func(url string) {
			resp, err := c.client.Get(url)
			ch <- &HttpResponse{url, resp, err}
			if err != nil && resp != nil && resp.StatusCode == http.StatusOK {
				resp.Body.Close()
			}
		}(url)
	}

	for {
		select {
		case r := <-ch:
			if r.Err != nil {
				log.Println("BundledHttpGet: ", r.Url, " ", r.Err)
			}
			c.HttpResponses = append(c.HttpResponses, r)
			if len(c.HttpResponses) == len(c.Urls) {
				break
			}
		case <-time.After(5 * time.Second):
		}
	}
}

func test() {
	pop := make(map[string]int)
	pop["test"]=1
	for k,v := range pop {
		fmt.Println(k,v)
	}
}