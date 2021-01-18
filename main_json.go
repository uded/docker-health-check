// +build json

package main

import (
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/sinhashubham95/jsonic"
	"github.com/urfave/cli/v2"
)

var (
	json      bool
	jsonCheck cli.StringSlice
)

func main() {
	newFlags := []cli.Flag{
		&cli.BoolFlag{
			Name:        "jsonResponse",
			Aliases:     []string{"json", "j"},
			Usage:       "the response should be treated as JSON",
			Value:       false,
			EnvVars:     []string{"HEALTHCHECK_JSON"},
			Destination: &json,
		},
		&cli.StringSliceFlag{
			Name:        "jsonCheck",
			Aliases:     []string{"json-check", "jc"},
			Usage:       "JSON check",
			Destination: &jsonCheck,
		},
	}
	app.Flags = append(app.Flags, newFlags...)

	app.Action = jsonActionFunc

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func jsonActionFunc(c *cli.Context) error {
	if err := httpCheck(c); err != nil {
		return err
	}

	if json && len(jsonCheck.Value()) > 0 {
		if data, err := ioutil.ReadAll(resp.Body); err != nil {
			cli.Exit(err.Error(), 1)
		} else {
			result, err := jsonCheckValues(data)
			if err != nil {
				cli.Exit(err.Error(), 1)
			}
			if !result {
				cli.Exit("given checks were not matching response", 1)
			}
		}
	}

	return nil
}

func jsonCheckValues(data []byte) (bool, error) {
	result := false
	j, err := jsonic.New(data)
	if err != nil {
		return false, err // Not JSON?
	}

	for _, v := range jsonCheck.Value() {
		if strings.Contains(v, "==") {
			check := strings.Split(v, "==")
			value, je := j.GetString(strings.TrimSpace(check[0]))
			if je == nil {
				result = strings.EqualFold(value, strings.TrimSpace(check[1]))
			}
		}
		if strings.Contains(v, "!=") {
			check := strings.Split(v, "!=")
			value, je := j.GetString(strings.TrimSpace(check[0]))
			if je == nil {
				result = strings.EqualFold(value, strings.TrimSpace(check[1]))
				result = value != strings.TrimSpace(check[1])
			}
		}
		if strings.Contains(v, " in ") {
			check := strings.Split(v, " in ")
			value, je := j.GetString(strings.TrimSpace(check[0]))
			if je == nil {
				checkV := strings.Split(strings.TrimSpace(check[0]), ",")
				for _, v := range checkV {
					result = strings.EqualFold(value, strings.TrimSpace(v))
				}
			}
		}
	}
	return result, nil
}
