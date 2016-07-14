package server

import (
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"

	"github.com/kardianos/osext"
	"github.com/russross/blackfriday"
)

// find available port from startPort to endPort
const startPort = 8080
const endPort = 8089

// relative to binary file 'mdv'
const htmlPath = "static/html/index.html"

var mdPath string

// WebContent is an inject interface for html template
type WebContent struct {
	MarkdownHtml interface{}
}

// Server lanuchs a http server with an available port
func Serve(path string) {
	mdPath = path

	http.HandleFunc("/", handleRoot)

	// find available port from startPort to endPort
	port := -1
	for p := startPort; p <= endPort; p++ {
		log.Printf("try listen on localhost:%d...\n", p)
		occupied, err := portInUse(p)
		if occupied {
			log.Printf("error: %v\n", err)
			continue
		}
		port = p
		break
	}

	if port == -1 {
		log.Fatalf("Sorry, all ports[%d~%d] are in use. I am tired.", startPort, endPort)
	}

	s := &http.Server{Addr: fmt.Sprintf(":%d", port)}
	log.Printf("server start on localhost:%d", port)
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

func portInUse(port int) (bool, error) {
	s := &http.Server{Addr: fmt.Sprintf(":%d", port)}
	ln, err := net.Listen("tcp", s.Addr)
	if err != nil {
		return true, err
	}

	err = ln.Close()
	if err != nil {
		return true, err
	}
	return false, nil
}
