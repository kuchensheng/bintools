package main

import (
	"net/http"
	"strings"
	"testing"
)

func TestCurl(t *testing.T) {
	http.Post("", "", strings.NewReader(""))
}
