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

func createOutgoingURL(c proxyConfig, incomingURL *url.URL) (outgoing url.URL, err error) {
	// get scheme, host, and path from config
	splitUrl := strings.Split(c.outgoingURL.String(), ":")
	outgoing.Scheme = splitUrl[0]
	outgoing.Host = strings.TrimPrefix(splitUrl[1], "//")
	// add in path, minus incoming path
	outgoing.Path = strings.TrimPrefix(incomingURL.String(), c.incomingPath)
	return outgoing, err
}

func handler(p *httputil.ReverseProxy, c proxyConfig) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// remove incoming path
		newUrl, err := createOutgoingURL(c, r.URL)
		if err != nil {
			log.Printf("Could not create new URL: %v", err)
		}
		r.Host = r.URL.Host
		log.Printf("%s -- /%v", c.name, newUrl)
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
