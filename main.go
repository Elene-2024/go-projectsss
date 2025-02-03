package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

type CacheConfig struct {
	TTL                  time.Duration
	CacheableMethods     []string
	CacheableStatusCodes []int
}

type CacheMiddleware struct {
	Transport http.RoundTripper
	Config    CacheConfig
	cache     sync.Map
}

type CacheEntry struct {
	Response  *http.Response
	Timestamp time.Time
}

func NewCacheMiddleware(transport http.RoundTripper, config CacheConfig) *CacheMiddleware {
	return &CacheMiddleware{
		Transport: transport,
		Config:    config,
	}
}

func (m *CacheMiddleware) RoundTrip(req *http.Request) (*http.Response, error) {

	cacheKey := req.URL.String()

	if cachedEntry, found := m.cache.Load(cacheKey); found {
		entry := cachedEntry.(CacheEntry)

		if time.Since(entry.Timestamp) < m.Config.TTL {

			fmt.Println("Returning cached response")
			return entry.Response, nil
		}

		m.cache.Delete(cacheKey)
	}

	resp, err := m.Transport.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	if isCacheable(req, resp, m.Config) {

		m.cache.Store(cacheKey, CacheEntry{
			Response:  resp,
			Timestamp: time.Now(),
		})
	}

	return resp, nil
}

func isCacheable(req *http.Request, resp *http.Response, config CacheConfig) bool {

	for _, method := range config.CacheableMethods {
		if req.Method == method {

			for _, statusCode := range config.CacheableStatusCodes {
				if resp.StatusCode == statusCode {
					return true
				}
			}
		}
	}
	return false
}

func main() {

	client := &http.Client{
		Transport: NewCacheMiddleware(http.DefaultTransport, CacheConfig{
			TTL:                  5 * time.Second, // Short TTL for testing
			CacheableMethods:     []string{"GET"},
			CacheableStatusCodes: []int{200},
		}),
	}

	url := "http://example.com/another-resource"

	resp, err := client.Get(url)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}
	fmt.Println("First Response (from server):")
	fmt.Println(string(body))

	time.Sleep(2 * time.Second)

	resp, err = client.Get(url)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}
	fmt.Println("\nSecond Response (from cache):")
	fmt.Println(string(body))

	time.Sleep(6 * time.Second)

	resp, err = client.Get(url)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}
	fmt.Println("\nThird Response (from server):")
	fmt.Println(string(body))
}
