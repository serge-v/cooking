package main

//go:generate go run server.go gen.go -gen

import (
	"os/exec"
	//	"io"
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	//	"bytes"
)

const (
	pageBreak = "\n<div class=\"page-break\"></div>"
)

var (
	recPage string
)

func getTopicListItem(dirPath string) string {
	if dirPath == "." {
		return ""
	}
	bytes, err := ioutil.ReadFile(dirPath + "/title.txt")
	if err != nil {
		panic(err)
	}
	title := strings.Trim(string(bytes), " \r\n")
	li := "<li><a href=\"" + dirPath + "\">" + title + "</li>\n"
	return li
}

type byName []os.FileInfo

func (f byName) Len() int      { return len(f) }
func (f byName) Swap(i, j int) { f[i], f[j] = f[j], f[i] }
func (f byName) Less(i, j int) bool {
	return f[i].Name() < f[j].Name()
}

func dirItems(dirPath string) string {
	f, err := os.Open(dirPath)
	if err != nil {
		panic(err)
	}

	finfos, err := f.Readdir(100)
	if err != nil {
		panic(err)
	}

	sort.Sort(byName(finfos))
	os.Mkdir("build/"+dirPath, 0777)

	contents := ""
	items := ""

	for _, fi := range finfos {
		fname := dirPath + "/" + fi.Name()
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

		tag := strings.Replace(fi.Name(), ".html", "", 1)
		tag = strings.ToLower(tag)

		li := "<li><a href=\"" + dirPath + "#" + tag + "\">" + title + "</li>\n"
		contents += li

		items += "<a name=\"" + tag + "\"></a>"
		items += string(chunk)
		items = strings.Replace(items, "/recbook/images/", "/images/", -1)
		items += pageBreak
	}

	text := strings.Replace(string(recPage), "{contents}", items, 1)

	_, err = os.Stat("build/" + dirPath)
	if os.IsNotExist(err) {
		os.MkdirAll("build/"+dirPath, 0755)
	}
	err = ioutil.WriteFile("build/"+dirPath+"/index.html", []byte(text), 0666)
	if err != nil {
		panic(err)
	}

	return contents
}

func dirContents(dirPath string) string {

	println(dirPath)

	f, err := os.Open(dirPath)
	if err != nil {
		log.Fatal(err)
	}

	contents := getTopicListItem(dirPath)

	finfos, err := f.Readdir(100)
	if err != nil {
		log.Fatal(err)
	}
	sort.Sort(byName(finfos))

	if len(finfos) == 0 {
		return contents + "."
	}

	contents += "<ul>\n"
	for _, fi := range finfos {
		if fi.IsDir() && (fi.Name() == ".git" || fi.Name() == "images" || fi.Name() == "build" || fi.Name() == "templates") {
			continue
		}

		if !fi.IsDir() {
			continue
		}

		subdirPath := dirPath + "/" + fi.Name()
		contents += dirContents(subdirPath)
	}

	contents += dirItems(dirPath)
	contents += "</ul>\n"
	return contents
}

func lintFile(fname string) {
	buf, err := ioutil.ReadFile(fname)
	if err != nil {
		fmt.Println(fname, ": error: cannot open", err)
		return
	}

	changes := false

	if buf[0] == 0xEF && buf[1] == 0xBB && buf[2] == 0xBF {
		fmt.Println(fname, ": warn: BOM")
		buf = buf[3:]
		changes = true
	}

	//	r := bufio.NewReader(strings.NewReader(string(bytes)))
	//	linenum := 0
	//	has_crlf := false

	s := string(buf)

	pos := strings.Index(s, "\r\n")
	if pos >= 0 {
		fmt.Println(fname, ": warn: CRLF")
		s = strings.Replace(s, "\r\n", "\n", -1)
		changes = true
	}

	/*	for {
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
	*/
	if changes {
		err = ioutil.WriteFile(fname, []byte(s), 0666)
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

	if strings.HasSuffix(path, ".html") || strings.HasSuffix(path, "/title.txt") {
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

func generateWebsite() {
	text, err := ioutil.ReadFile("templates/main.html")
	if err != nil {
		log.Fatal(err)
	}
	mainPage := string(text)

	text, err = ioutil.ReadFile("templates/recpage.html")
	if err != nil {
		log.Fatal(err)
	}
	recPage = string(text)

	contents := dirContents(".")

	out := strings.Replace(string(mainPage), "{contents}", contents, 1)
	err = ioutil.WriteFile("build/index.html", []byte(out), 0666)
	if err != nil {
		log.Fatal(err)
	}

	cmd := exec.Command("cp", "-R", "images", "build/")
	if out, err := cmd.CombinedOutput(); err != nil {
		log.Fatalln(string(out), err)
	}

	os.Remove("cooking.zip")

	cmd = exec.Command("zip", "-r", "../cooking.zip", ".")
	cmd.Dir = "build"
	if err = cmd.Run(); err != nil {
		log.Fatal(err)
	}
	log.Println("zip created")
}
