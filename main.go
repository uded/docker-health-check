package main

import (
	"net/http"
	"strings"
	"time"

	"github.com/urfave/cli/v2"
)

// globals
var (
	app        *cli.App
	url        string
	httpVerb   string
	statusCode int
	timeOut    int
)

func init() {
	app = &cli.App{
		Name:    "healtcheck",
		Version: "0.0.4",
		Usage:   "Hits an endpoint for you.  healthcheck -url=http://localhost/ping",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "url, U",
				Usage:       "the url to hit (required)",
				EnvVars:     []string{"HEALTHCHECK_URL"},
				Destination: &url,
			},
			&cli.StringSliceFlag{
				Name:  "headers, H",
				Usage: "specify a header and value for the request (-H=key:value)",
			},
			&cli.StringFlag{
				Name:        "verb, V",
				Usage:       "the HTTP verb to use",
				Value:       "GET",
				EnvVars:     []string{"HEALTHCHECK_VERB"},
				Destination: &httpVerb,
			},
			&cli.IntFlag{
				Name:        "code, C",
				Usage:       "expected response code",
				Value:       http.StatusOK,
				EnvVars:     []string{"HEALTHCHECK_RESPONSECODE"},
				Destination: &statusCode,
			},
			&cli.IntFlag{
				Name:        "timeout, T",
				Usage:       "timeout for HTTP connection",
				Value:       0,
				EnvVars:     []string{"HEALTHCHECK_TIMEOUT"},
				Destination: &timeOut,
			},
		},
	}

}

func httpResp(c *cli.Context) (*http.Response, error) {
	if len(url) < 0 {
		return nil, cli.Exit("url length must be > 0 ", 1)
	}

	req, err := http.NewRequest(httpVerb, url, nil)
	if err != nil {
		return nil, cli.Exit(err.Error(), 1)
	}
	for _, str := range c.StringSlice("headers") {
		kv := strings.Split(str, ":")
		if len(kv) == 2 {
			req.Header.Add(kv[0], kv[1])
		} else {
			return nil, cli.Exit("header field must be in the format \"key:value\"", 1)
		}
	}
	req.Close = true

	var client *http.Client
	if timeOut > 0 {
		timeout := time.Duration(timeOut) * time.Second
		client = &http.Client{Timeout: timeout}
	} else {
		client = &http.Client{}
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, cli.Exit(err.Error(), 1)
	}
	if resp != nil {
		defer func(r *http.Response) {
			if r.Body != nil {
				r.Body.Close()
			}
		}(resp)
		return resp, nil
	} else {
		return nil, cli.Exit("no erro but also no response was returned", 1)
	}
}
