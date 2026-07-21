package main

import (
	"archive/zip"
	"flag"
	"log"
	"net/http"

	"golang.org/x/tools/godoc/vfs/httpfs"
	"golang.org/x/tools/godoc/vfs/zipfs"
)

var (
	gen  = flag.Bool("gen", false, "generate website")
	lint = flag.String("lint", "", "lint html chunk file")
	addr = flag.String("addr", ":8080", "address to listen")
)

func main() {
	flag.Parse()
	log.SetFlags(log.Lshortfile)

	if *lint != "" {
		lintFiles(*lint)
		return
	}

	if *gen {
		generateWebsite()
		return
	}

	rc, err := zip.OpenReader("cooking.zip")
	if err != nil {
		panic(err)
	}
	defer rc.Close()

	fs := httpfs.New(zipfs.New(rc, "cooking"))
	log.Println("listening on", "http://"+*addr)
	log.Fatal(http.ListenAndServe(*addr, http.FileServer(fs)))
}
