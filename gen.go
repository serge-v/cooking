package main

import (
	"io"
	"os"
	"fmt"
	"log"
	"io/ioutil"
	"strings"
	"bufio"
	"flag"
	"net/http"
	"path/filepath"
)

const (
	page_break = "\n<div class=\"page-break\"></div>"
)

var (
	rec_page string
	httpFlag = flag.Bool("server", false, "Start debug server")
	lintFlag = flag.String("lint", "", "Lint html chunk file")
)

func getTopicListItem(dir_path string) string {
	if dir_path == "." {
		return ""
	}
	bytes, err := ioutil.ReadFile(dir_path + "/title.txt")
	if err != nil {
		panic(err)
	}
	title := strings.Trim(string(bytes), " \r\n")
	li := "<li><a href=\"" + dir_path + "\">" + title  + "</li>\n"
	return li
}

func dirItems(dir_path string) string {
	f, err := os.Open(dir_path)
	if err != nil {
		panic(err)
	}
	
	finfos, err := f.Readdir(100)
	if err != nil {
		panic(err)
	}

	os.Mkdir("build/" + dir_path, 0777)

	contents := ""
	items := ""

	for _, fi := range(finfos) {
		fname := dir_path + "/" + fi.Name()
		println(fname)

		if !strings.HasSuffix(fname, ".html") {
			continue
		}

		chunk, err := ioutil.ReadFile(fname)
		if err != nil {
			panic(err)
		}
		
		page := string(chunk)
		r := bufio.NewReader(strings.NewReader(page))
		title, err := r.ReadString('\n')
		title = strings.Replace(title, "<h3>", "", 1)
		title = strings.Replace(title, "</h3>", "", 1)
		if strings.IndexAny(title, "<>") >= 0 {
			panic(fmt.Sprintf("%s: error: invalid title. Should be in format <h3>Title</h3>.", fname))
		}
		
		li := "<li><a href=\"" + dir_path + "#" + fi.Name() + "\">" + title  + "</li>\n"
		contents += li

		items += "<a name=\"" + fi.Name() + "\"></a>"
		items += string(chunk)
		items = strings.Replace(items, "/recbook/images/", "/images/", -1)
		items += page_break
	}
	

	text := strings.Replace(string(rec_page), "{contents}", items, 1)

	err = ioutil.WriteFile("build/" + dir_path + "/index.html", []byte(text), 0666)
	if err != nil {
		panic(err)
	}
	
	return contents
}

func dirContents(dir_path string) string {

	println(dir_path)

	f, err := os.Open(dir_path)
	if err != nil {
		log.Fatal(err)
	}
	
	contents := getTopicListItem(dir_path)

	finfos, err := f.Readdir(100)
	if err != nil {
		log.Fatal(err)
	}

	if len(finfos) == 0 {
		return contents + "."
	}

	contents += "<ul>\n"
	for _, fi := range(finfos) {
		if fi.IsDir() && (fi.Name() == ".git" || fi.Name() == "images" || fi.Name() == "build" || fi.Name() == "templates") {
			continue
		}
		
		if !fi.IsDir() {
			continue
		}

		subdir_path := dir_path + "/" + fi.Name()		
		contents += dirContents(subdir_path)
	}
	
	contents += dirItems(dir_path)
	contents += "</ul>\n"
	return contents
}

func lintFile(fname string) {
	bytes, err := ioutil.ReadFile(fname)
	if err != nil {
		fmt.Println(fname, ": error: cannot open", err)
		return
	}
	
	changes := false

	if bytes[0] == 0xEF && bytes[1] == 0xBB && bytes[2] == 0xBF {
		fmt.Println(fname, ": warn: BOM")
		bytes = bytes[3:]
		changes = true
	}
	
	r := bufio.NewReader(strings.NewReader(string(bytes)))
	linenum := 0
	for {
		s, err := r.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		linenum++
		
		pos := strings.IndexAny(s, "\r")
		if pos >= 0 {
			fmt.Printf("%s:%d: warn: invalid charachter at pos: %d\n", fname, linenum, pos)
		}
	}

	if changes {
//		err = ioutil.WriteFile(fname, bytes, 0666)
		fmt.Println(fname, ": changed")
	}
}

func walk(path string, info os.FileInfo, err error) error {

	if strings.HasSuffix(path, ".git") || strings.HasSuffix(path, "build") {
		return filepath.SkipDir
	}

	if info.IsDir() {
		return nil
	}
	
	if strings.HasSuffix(path, ".html") {
		lintFile(path)
	}

	return nil
}

func lintFiles(root string) {
	err := filepath.Walk(root, walk)
	if err != nil {
		panic(err)
	}
}

func main() {
	flag.Parse()

	if *httpFlag {
		fmt.Println("starting server on http://localhost:9000")
		panic(http.ListenAndServe(":9000", http.FileServer(http.Dir("build"))))
		return
	}
	
	if *lintFlag != "" {
		lintFiles(*lintFlag)
		return
	}

	text, err := ioutil.ReadFile("templates/main.html")
	if err != nil {
		log.Fatal(err)
	}
	main_page := string(text)

	text, err = ioutil.ReadFile("templates/recpage.html")
	if err != nil {
		log.Fatal(err)
	}
	rec_page = string(text)

	contents := dirContents(".")

	out := strings.Replace(string(main_page), "{contents}", contents, 1)
	err = ioutil.WriteFile("build/index.html", []byte(out), 0666)
	if err != nil {
		log.Fatal(err)
	}
}
