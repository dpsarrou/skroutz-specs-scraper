package crawler

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type HttpClient interface {
	Get(url string) ([]byte, error)
}

type Client struct {
	UserAgent string
	Cookie    string
}

func NewClient(userAgent string, cookie string) HttpClient {
	if userAgent == "" {
		userAgent = "skroutz-specs/1.0"
	}

	c := &Client{
		UserAgent: userAgent,
		Cookie:    cookie,
	}
	ch := &CachingClient{c}
	return ch
}

func (c *Client) Get(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", c.UserAgent)
	if c.Cookie != "" {
		req.Header.Set("Cookie", c.Cookie)
	}

	// TODO remove rate limit
	time.Sleep(1 * time.Second)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, err
	}

	data, err := io.ReadAll(res.Body)
	return data, err
}

type CachingClient struct {
	Client *Client
}

func (c *CachingClient) Get(url string) ([]byte, error) {
	hash := fmt.Sprintf("%x", md5.Sum([]byte(url))) + ".html"
	f, err := os.Open("cache/" + hash)
	if err == nil {
		log.Printf("[Client]\t[Cache HIT]\tRetreiving page from cache for url %s, hash: %s", url, hash)
		data, err := io.ReadAll(f)
		return data, err
	}
	log.Printf("[Client]\t[Cache MISS]\tSending request to %s and saving to %s", url, hash)
	data, err := c.Client.Get(url)
	if err != nil {
		return nil, err
	}
	f, err = os.Create("cache/" + hash)
	if err != nil {
		return nil, err
	}
	f.Write(data)
	return data, err
}
