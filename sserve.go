package main

import (
	"compress/gzip"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"runtime"
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

//mkcert
func mkcert() {
	// set the right executable according to the system
	mkcert := "mkcert/"
	switch runtime.GOOS {
	case "darwin":
		mkcert = mkcert + "mkcert-v1.2.0-darwin-amd64"
	case "linux":
		mkcert = mkcert + "mkcert-v1.2.0-linux-amd64"
	case "windows":
		mkcert = mkcert + "mkcert-v1.2.0-windows-amd64.exe"
	default:
		log.Fatal("Your system is not supported. Sorry.")
		os.Exit(1)
	}

	// generate the certificate
	if _, err := exec.Command(mkcert, "-install", "-cert-file", "localhost.crt",
		"-key-file", "localhost.key", "localhost").Output(); err != nil {
		log.Fatal(err.Error())
		os.Exit(1)
	}
}

// check if file exists
func exist(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	} else if os.IsNotExist(err) {
		return false
	} else {
		os.Exit(1)
		return false
	}
}

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
func serve(path string, port string, minify bool, compression bool) {
	// file server
	fs := http.FileServer(http.Dir(path))

	// minifier middleware
	if minify {
		fs = getMinifier(fs)
		log.Println("Code minified.")
	}

	// gzip middleware
	if compression {
		log.Println("gzip compression activated.")
		fs = getGzipper(fs)
	}

	// http handler
	http.Handle("/", fs)

	// start server
	address := "https://localhost"
	if port != "443" {
		address = address + ":" + port
	}
	log.Println("Serving " + path + " on port " + port + ". Checkout at " + address + ".")
	err := http.ListenAndServeTLS(":"+port, "localhost.crt", "localhost.key", nil)
	log.Fatal(err)
}

// CLI interface
func main() {
	// parse CLI arguments
	log.SetFlags(0)
	var portFlag = flag.String("port", "443", "Port number of the server.")
	var redirectFlag = flag.Bool("redirect", true, "If true activate the http redirect.")
	var minifyFlag = flag.Bool("minify", true, "If true minify the code.")
	var compressionFlag = flag.Bool("compression", true, "If true activate gzip compression.")
	flag.Parse()

	// read the static path from the cli args or use the working directory
	path := "./"
	if len(flag.Args()) > 0 {
		path = flag.Args()[0]
	}

	// ensure that the certificate files exists
	if !exist("localhost.crt") || !exist("localhost.key") {
		mkcert()
	}

	// activate the redirect
	if *redirectFlag {
		redirectHTTP()
		log.Println("http redirect activated.")
	}

	// start the server
	serve(path, *portFlag, *minifyFlag, *compressionFlag)
}
