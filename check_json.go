package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	flags "github.com/jessevdk/go-flags"
	"github.com/mackerelio/checkers"
	jsonpointer "github.com/mattn/go-jsonpointer"
)

var opts struct {
	URL                string `short:"u" long:"url" required:"true" description:"A URL to connect to"`
	Pointer            string `short:"p" long:"pointer" required:"true" description:"JSON Pointer"`
	NoCheckCertificate bool   `long:"no-check-certificate" description:"Do Not check certificate"`
}

func main() {
	ckr := run(os.Args[1:])
	ckr.Name = "JSON API Response"
	ckr.Exit()
}

func run(args []string) *checkers.Checker {
	_, err := flags.ParseArgs(&opts, args)
	if err != nil {
		os.Exit(1)
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: opts.NoCheckCertificate,
		},
	}
	client := &http.Client{Transport: tr}

	resp, err := client.Get(opts.URL)
	if err != nil {
		return checkers.Critical(err.Error())
	}
	defer resp.Body.Close()

	var i interface{}
	decoder := json.NewDecoder(resp.Body)
	if err = decoder.Decode(&i); err != nil {
		return checkers.Critical(err.Error())
	}

	value, err := jsonpointer.Get(i, opts.Pointer)
	if err != nil {
		return checkers.Warning(err.Error())
	}

	checkSt := checkers.OK
	msg := fmt.Sprintf("%s: %s", opts.Pointer, value)

	return checkers.NewChecker(checkSt, msg)
}
