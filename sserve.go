package main

import (
	"compress/gzip"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"

	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
	"github.com/tdewolff/minify/v2/html"
	"github.com/tdewolff/minify/v2/js"
	"github.com/tdewolff/minify/v2/json"
	"github.com/tdewolff/minify/v2/svg"
	"github.com/tdewolff/minify/v2/xml"
)

// AUXILIARY FUNCTIONS

// http to https rederect handler
func redirect(w http.ResponseWriter, req *http.Request) {
	target := "https://" + req.Host + req.URL.Path
	if len(req.URL.RawQuery) > 0 {
		target += "?" + req.URL.RawQuery
	}
	log.Printf("redirect to: %s", target)
	http.Redirect(w, req, target, http.StatusTemporaryRedirect)
}

// minify middleware
func getMinifier(httpHandler http.Handler) http.Handler {
	m := minify.New()
	m.AddFunc("text/css", css.Minify)
	m.AddFunc("text/html", html.Minify)
	m.AddFunc("image/svg+xml", svg.Minify)
	m.AddFuncRegexp(regexp.MustCompile("^(application|text)/(x-)?(java|ecma)script$"), js.Minify)
	m.AddFuncRegexp(regexp.MustCompile("[/+]json$"), json.Minify)
	m.AddFuncRegexp(regexp.MustCompile("[/+]xml$"), xml.Minify)
	return m.Middleware(httpHandler)
}

// gzip middleware
var gzPool = sync.Pool{
	New: func() interface{} {
		w := gzip.NewWriter(ioutil.Discard)
		return w
	},
}

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w *gzipResponseWriter) WriteHeader(status int) {
	w.Header().Del("Content-Length")
	w.ResponseWriter.WriteHeader(status)
}

func (w *gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func getGzipper(httpHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			httpHandler.ServeHTTP(w, r)
			return
		}
		w.Header().Set("Content-Encoding", "gzip")
		gz := gzPool.Get().(*gzip.Writer)
		defer gzPool.Put(gz)
		gz.Reset(w)
		defer gz.Close()
		httpHandler.ServeHTTP(&gzipResponseWriter{ResponseWriter: w, Writer: gz}, r)
	})
}

// MAIN FUNCTIONS

// Redirect every http request to https
func redirectHTTP() {
	go http.ListenAndServe(":80", http.HandlerFunc(redirect))
}

// Serve a static content
func serve(path string, port string) {
	// file server
	fs := http.FileServer(http.Dir(path))
	// minifier middleware
	fsmin := getMinifier(fs)
	// gzip middleware
	fsmingzip := getGzipper(fsmin)
	// http handler
	http.Handle("/", fsmingzip)
	// start server
	log.Println("Serving " + path + " on port " + port)
	err := http.ListenAndServeTLS(":"+port, "localhost.crt", "localhost.key", nil)
	log.Fatal(err)
}

// CLI interface
func main() {
	// read the static path from the cli args or use the working directory
	path := "./"
	if len(os.Args) > 1 {
		path = os.Args[1]
	}
	// start the server
	serve(path, "443")
	redirectHTTP()
}
