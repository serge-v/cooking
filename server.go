package main

import (
	"net/http"
	"net/http/cgi"
)

func main() {
	err := cgi.Serve(http.FileServer(http.Dir("build")))
	if err != nil {
		panic(err)
	}
}
