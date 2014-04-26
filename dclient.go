package doracle

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type Client struct {
	shutdown    bool
	req         chan chan int64
	servers     []string
	cacheServer string
}

func NewClient(addresses []string) (*Client, error) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	if len(addresses) <= 0 {
		return nil, errors.New("no server")
	}

	cl := &Client{
		shutdown: false,
		req:      make(chan chan int64, 100000),
		servers:  addresses,
	}
	cl.cacheServer = cl.servers[0]
	go cl.start()
	return cl, nil
}

// Close the client after all TS responses are returned
func (c *Client) Close() {
	c.shutdown = true
}

func (c *Client) TS() (int64, error) {
	if c.shutdown {
		return -1, errors.New("close")
	}
	ch := make(chan int64)
	c.req <- ch
	if ts := <-ch; ts >= 0 {
		return ts, nil
	} else {
		return -1, errors.New("invalid ts")
	}
}

func (c *Client) GetTS(num int32) (int64, error) {
	if c.shutdown {
		return -1, errors.New("already close")
	}

	numStr := strconv.FormatInt(int64(num), 10)
	var b bytes.Buffer
	b.WriteString(numStr)

	resp, err := http.Post(fmt.Sprintf("http://%s/doracle", c.cacheServer), "text/plain", &b)
	for err != nil || resp == nil || resp.StatusCode != http.StatusOK {
		for _, server := range c.servers {
			var b bytes.Buffer
			b.WriteString(numStr)
			resp, err = http.Post(fmt.Sprintf("http://%s/doracle", server), "text/plain", &b)
			log.Println("tried: ", server)
			if err == nil && resp.StatusCode == http.StatusOK {
				c.cacheServer = server
				break
			}
			if err != nil {
				log.Println("error: ", err)
			}
			if resp != nil {
				// Must close the response body
				resp.Body.Close()
				log.Println("code: ", resp.StatusCode, server)
			}
		}
	}

	defer resp.Body.Close()
	if r, err := ioutil.ReadAll(resp.Body); err != nil {
		return -1, err
	} else {
		return strconv.ParseInt(string(r), 10, 64)
	}
}

func (c *Client) start() {
	for !c.shutdown {
		ch := <-c.req
		l := len(c.req)
		// batch count
		// log.Println(l)
		ts, err := c.GetTS(int32(l + 1))
		if err != nil {
			log.Println("get ts error", err)
			c.shutdown = true
			break
		}
		ch <- ts
		for i := 1; i <= l; i++ {
			ch = <-c.req
			ch <- ts - int64(i)
		}
	}

	close(c.req)
}
