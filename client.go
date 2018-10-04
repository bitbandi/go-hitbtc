package hitbtc

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
)

type client struct {
	apiKey      string
	apiSecret   string
	httpClient  *http.Client
	httpTimeout time.Duration
	debug       bool
}

// NewClient return a new HitBtc HTTP client
func NewClient(apiKey, apiSecret string) (c *client) {
	return &client{apiKey, apiSecret, &http.Client{}, 30 * time.Second, false}
}

// NewClientWithCustomHttpConfig returns a new HitBtc HTTP client using the predefined http client
func NewClientWithCustomHttpConfig(apiKey, apiSecret string, httpClient *http.Client) (c *client) {
	timeout := httpClient.Timeout
	if timeout <= 0 {
		timeout = 30 * time.Second
	}
	return &client{apiKey, apiSecret, httpClient, timeout, false}
}

// NewClient returns a new HitBtc HTTP client with custom timeout
func NewClientWithCustomTimeout(apiKey, apiSecret string, timeout time.Duration) (c *client) {
	return &client{apiKey, apiSecret, &http.Client{}, timeout, false}
}

func (c client) dumpRequest(r *http.Request) {
	if r == nil {
		log.Print("dumpReq ok: <nil>")
		return
	}
	dump, err := httputil.DumpRequest(r, true)
	if err != nil {
		log.Print("dumpReq err:", err)
	} else {
		log.Print("dumpReq ok:", string(dump))
	}
}

func (c client) dumpResponse(r *http.Response) {
	if r == nil {
		log.Print("dumpResponse ok: <nil>")
		return
	}
	dump, err := httputil.DumpResponse(r, true)
	if err != nil {
		log.Print("dumpResponse err:", err)
	} else {
		log.Print("dumpResponse ok:", string(dump))
	}
}

// doTimeoutRequest do a HTTP request with timeout
func (c *client) doTimeoutRequest(timer *time.Timer, req *http.Request) (*http.Response, error) {
	// Do the request in the background so we can check the timeout
	type result struct {
		resp *http.Response
		err  error
	}
	done := make(chan result, 1)
	go func() {
		if c.debug {
			c.dumpRequest(req)
		}
		resp, err := c.httpClient.Do(req)
		if c.debug {
			c.dumpResponse(resp)
		}
		done <- result{resp, err}
	}()
	// Wait for the read or the timeout
	select {
	case r := <-done:
		return r.resp, r.err
	case <-timer.C:
		return nil, errors.New("timeout on reading data from HitBtc API")
	}
}

// do prepare and process HTTP request to HitBtc API
func (c *client) do(method string, ressource string, payload map[string]string, authNeeded bool) (response []byte, err error) {
	connectTimer := time.NewTimer(c.httpTimeout)

	var rawurl string
	if strings.HasPrefix(ressource, "http") {
		rawurl = ressource
	} else {
		rawurl = fmt.Sprintf("%s/%s", API_BASE, ressource)
	}
	var formData string
	if method == "GET" {
		var URL *url.URL
		URL, err = url.Parse(rawurl)
		if err != nil {
			return
		}
		q := URL.Query()
		for key, value := range payload {
			q.Set(key, value)
		}
		formData = q.Encode()
		URL.RawQuery = formData
		rawurl = URL.String()
	} else {
		formValues := url.Values{}
		for key, value := range payload {
			formValues.Set(key, value)
		}
		formData = formValues.Encode()
	}
	req, err := http.NewRequest(method, rawurl, strings.NewReader(formData))
	if err != nil {
		return
	}
	req.Header.Add("Accept", "application/json")

	// Auth
	if authNeeded {
		if len(c.apiKey) == 0 || len(c.apiSecret) == 0 {
			err = errors.New("You need to set API Key and API Secret to call this method")
			return
		}
		req.SetBasicAuth(c.apiKey, c.apiSecret)
	}

	resp, err := c.doTimeoutRequest(connectTimer, req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	response, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return response, err
	}
	if resp.StatusCode != 200 && resp.StatusCode != 401 {
		return response, errors.New(resp.Status)
	}
	return response, err
}
