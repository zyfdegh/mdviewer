package server

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/russross/blackfriday"
)

var defaultPort = 8080

var mdPath string

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
	io.WriteString(w, string(render(content)))
}

func render(content []byte) string {
	return string(blackfriday.MarkdownCommon(content))
}
