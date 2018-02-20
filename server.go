package main

import (
	"archive/zip"
	"flag"
	"fmt"
	"net/http"
	"net/http/cgi"

	"golang.org/x/tools/godoc/vfs/httpfs"
	"golang.org/x/tools/godoc/vfs/zipfs"
)

var (
	gen       = flag.Bool("gen", false, "generate website")
	cgiServer = flag.Bool("cgi", false, "start debug cgi server")
	lint      = flag.String("lint", "", "lint html chunk file")
)

func runStandalone() {
	if *lint != "" {
		lintFiles(*lint)
		return
	}

	if *cgiServer {
		fmt.Println("starting server on http://localhost:9001")
		api := cgi.Handler{}
		api.Path = "cooking"
		err := http.ListenAndServe(":9001", &api)
		if err != nil {
			panic(err)
		}
	}

	if *gen {
		generateWebsite()
	}
}

func main() {
	flag.Parse()
	if flag.NFlag() != 0 {
		runStandalone()
		return
	}

	rc, err := zip.OpenReader("cooking.zip")
	if err != nil {
		panic(err)
	}
	defer rc.Close()
	fs := httpfs.New(zipfs.New(rc, "cooking"))
	err = cgi.Serve(http.FileServer(fs))
	if err != nil {
		panic(err)
	}
}
