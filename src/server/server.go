package server

import (
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/kardianos/osext"
	"github.com/russross/blackfriday"
)

const defaultPort = 8080

// relative to binary file 'mdv'
const htmlPath = "static/html/index.html"

var mdPath string

type WebContent struct {
	MarkdownHtml interface{}
}

// lanuch a http server with default port
func Serve(path string) {
	mdPath = path

	http.HandleFunc("/", handleRoot)

	s := &http.Server{Addr: fmt.Sprintf(":%d", defaultPort)}

	log.Printf("server start on localhost:%d", defaultPort)
	if err := s.ListenAndServe(); err != nil {
		log.Fatalf("start server error: %v", err)
	}
}

func handleRoot(w http.ResponseWriter, req *http.Request) {
	content, err := ioutil.ReadFile(mdPath)
	if err != nil {
		log.Printf("read file error: %v", err)
		return
	}
	// convert *.md file to basic html
	markdownHtml := blackfriday.MarkdownCommon(content)
	// convert html to rendered html
	wc := WebContent{MarkdownHtml: template.HTML(markdownHtml)}

	// check if template html file exist
	// like "/usr/local/bin/mdv"
	dir, err := osext.ExecutableFolder()
	if err != nil {
		log.Printf("get mdv binary file directory error: %v", err)
		return
	}

	path := filepath.Join(dir, htmlPath)

	// like "/usr/local/bin/static/html/index.html"
	if _, err := os.Stat(path); err != nil {
		log.Printf("get template html error: %v", err)
		return
	}

	// parse template
	t, err := template.ParseFiles(path)
	if err != nil {
		log.Fatalf("parse file to html template error: %v", err)
		io.WriteString(w, string(content))
	}

	// inject content to template
	err = t.Execute(w, wc)
	if err != nil {
		log.Printf("execute html template error: %v", err)
		io.WriteString(w, string(content))
	}
}
