package main

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
)

type proxyConfig struct {
	incomingPath string
	outgoingURL  *url.URL
	name         string
}

var logFatal = log.Fatalf

// exists on error
func readInConfig() (cfg []proxyConfig) {
	services := strings.Split(os.Getenv("services"), ",")
	log.Printf("read in services %v", services)
	for _, s := range services {
		p := fmt.Sprintf("%s_incoming_path", s)
		d := fmt.Sprintf("%s_outgoing_url", s)
		if os.Getenv(p) == "" || os.Getenv(d) == "" {
			logFatal("%s and %s must be defined as environment variables", p, d)
		}
		remote, err := url.Parse(os.Getenv(d))
		if err != nil {
			logFatal("bad _outgoing_url: %v", err)
		}
		cfg = append(cfg, proxyConfig{os.Getenv(p), remote, s})
	}
	return cfg
}

func createOutgoingURL(c proxyConfig, incomingURL *url.URL) (outgoing url.URL) {
	// get scheme, host, and path from config
	splitUrl := strings.Split(c.outgoingURL.String(), "://")
	outgoing.Scheme = splitUrl[0]
	outgoing.Host = splitUrl[1]
	// add in path, minus incoming path
	outgoing.Path = strings.TrimPrefix(incomingURL.String(), c.incomingPath)
	// remove everything before query
	splitQuery := strings.Split(outgoing.Path, "?")
	outgoing.Path = splitQuery[0]
	if len(splitQuery) == 2 {

		outgoing.RawQuery = splitQuery[1]
		log.Printf("%v", outgoing.String())
	}
	return outgoing
}

func handler(p *httputil.ReverseProxy, c proxyConfig) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// remove incoming path
		newUrl := createOutgoingURL(c, r.URL)
		r.URL = &newUrl
		r.Host = newUrl.Host
		r.RequestURI = ""
		log.Printf("handler: %s, path: %v, redirect: %v", c.name, c.incomingPath, newUrl.String())
		p.ServeHTTP(w, r)
	}
}

func serveReverseProxy(cfg []proxyConfig) {
	// loop through config creating routes
	for _, c := range cfg {
		proxy := httputil.NewSingleHostReverseProxy(c.outgoingURL)
		http.HandleFunc(c.incomingPath, handler(proxy, c))
	}

	port := fmt.Sprintf(":%s", os.Getenv("PORT"))
	log.Printf("Serving on port %s", port)
	logFatal("%v", http.ListenAndServe(port, nil))

}

func main() {
	cfg := readInConfig()
	log.Printf("read in config %s", spew.Sdump(cfg))
	serveReverseProxy(cfg)
}
