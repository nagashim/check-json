package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mackerelio/checkers"
	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"result":{"OK":"1"}}`)
	}))
	defer ts.Close()

	ckr := run([]string{"-u", ts.URL, "-p", "/result/OK"})
	assert.Equal(t, checkers.OK, ckr.Status, "chr.Status should be CRITICAL")
	assert.Equal(t, "/result/OK: 1", ckr.Message, "something went wrong")

	ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `{"result":{"OK":1}}`)
	}))
	defer ts.Close()

	ckr = run([]string{"-u", ts.URL, "-p", "/result/OK"})
	assert.Equal(t, checkers.OK, ckr.Status, "chr.Status should be CRITICAL")
	assert.Equal(t, "/result/OK: 1", ckr.Message, "something went wrong")

	ckr = run([]string{"-u", ts.URL, "-p", "/timestamp"})
	assert.Equal(t, checkers.WARNING, ckr.Status, "chr.Status should be WARNING")
	assert.Equal(t, `Invalid JSON pointer: "/timestamp": reflect: call of reflect.Value.Interface on zero Value`, ckr.Message, "cannot get pointer value")
}

func TestNotJsonResponse(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `not json response`)
	}))
	defer ts.Close()

	ckr := run([]string{"-u", ts.URL, "-p", "/result/OK"})
	assert.Equal(t, checkers.CRITICAL, ckr.Status, "chr.Status should be CRITICAL")
	assert.Equal(t, "invalid character 'o' in literal null (expecting 'u')", ckr.Message, "response body decode error")
}

func TestInvalidUrl(t *testing.T) {
	ckr := run([]string{"-u", "hoge", "-p", "/result/ok"})
	assert.Equal(t, checkers.CRITICAL, ckr.Status, "chr.Status should be CRITICAL")
	assert.Equal(t, `Get hoge: unsupported protocol scheme ""`, ckr.Message, "something went wrong")
}

func TestNoCheckCertificate(t *testing.T) {
	ckr := run([]string{"-u", "hoge", "-p", "/result/ok", "--no-check-certificate"})
	assert.Equal(t, checkers.CRITICAL, ckr.Status, "chr.Status should be CRITICAL")
	assert.Equal(t, `Get hoge: unsupported protocol scheme ""`, ckr.Message, "something went wrong")
}
