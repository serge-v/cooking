package main

import (
	"os"
	"fmt"
	"log"
	"io/ioutil"
	"strings"
	"bufio"
	"flag"
	"net/http"
)

var (
	rec_page string
	httpFlag = flag.Bool("server", false, "Start debug server")
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
	li := "<li><a href=\"" + dir_path + "/index.html\">" + title  + "</li>\n"
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

	os.Mkdir("gen/" + dir_path, 0777)

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
		
		li := "<li><a href=\"" + dir_path + "/index.html#" + fi.Name() + "\">" + title  + "</li>\n"
		contents += li

		items += "<a name=\"" + fi.Name() + "\"></a>"
		items += string(chunk)
		items = strings.Replace(items, "/recbook/images/", "/images/", 1)
		items += "\n<div class=\"page-break\"></div>"
	}
	

	text := strings.Replace(string(rec_page), "{contents}", items, 1)

	err = ioutil.WriteFile("gen/" + dir_path + "/index.html", []byte(text), 0666)
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
		if fi.IsDir() && (fi.Name() == ".git" || fi.Name() == "images" || fi.Name() == "gen" || fi.Name() == "templates") {
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

func main() {
	flag.Parse()

	if *httpFlag {
		fmt.Println("starting server on http://localhost:9000")
		panic(http.ListenAndServe(":9000", http.FileServer(http.Dir("gen"))))
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
	err = ioutil.WriteFile("gen/index.html", []byte(out), 0666)
	if err != nil {
		log.Fatal(err)
	}
}
