package server

import (
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/russross/blackfriday"
)

var defaultPort = 8080

var mdPath string

type WebContent struct {
	MarkdownHtml interface{}
}

// lanuch a http server with default port
func Serve(path string) {
	mdPath = path

	http.HandleFunc("/", handleRoot)

	s := &http.Server{Addr: fmt.Sprintf(":%d", defaultPort)}

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
	t, err := template.ParseFiles("static/html/index.html")
	if err != nil {
		log.Fatalf("parse file to html template error: %v", err)
		io.WriteString(w, string(content))
	}

	err = t.Execute(w, wc)
	if err != nil {
		log.Fatalf("execute html template error: %v", err)
		io.WriteString(w, string(content))
	}
}
