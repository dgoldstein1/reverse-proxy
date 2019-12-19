package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/url"
	"os"
	"testing"
)

func TestReadInConfig(t *testing.T) {
	originLogFatalf := logFatal
	defer func() { logFatal = originLogFatalf }()
	logs := []string{}
	logFatal = func(format string, args ...interface{}) {
		if len(args) > 0 {
			logs = append(logs, fmt.Sprintf(format, args))
		} else {
			logs = append(logs, format)
		}
	}

	type Test struct {
		name              string
		before            func()
		after             func()
		expectedErrLength int
	}

	tests := []Test{
		Test{
			name: "correctly validates good config",
			before: func() {
				os.Setenv("services", "biggraph")
				os.Setenv("biggraph_incoming_path", "/services/biggraph")
				os.Setenv("biggraph_outgoing_url", "http://google.com")
			},
			after: func() {
				os.Unsetenv("services")
				os.Unsetenv("biggraph_incoming_path")
				os.Unsetenv("biggraph_outgoing_url")
			},
			expectedErrLength: 0,
		},
		Test{
			name: "fails when outgoing destination isn't supplied for service",
			before: func() {
				os.Setenv("services", "biggraph")
				os.Setenv("biggraph_incoming_path", "/services/biggraph")
			},
			after: func() {
				os.Unsetenv("services")
				os.Unsetenv("biggraph_incoming_path")
				os.Unsetenv("biggraph_outgoing_url")
			},
			expectedErrLength: 1,
		},
		Test{
			name: "fails on bad url for _outgoing_url",
			before: func() {
				os.Setenv("services", "biggraph")
				os.Setenv("biggraph_incoming_path", "/services/biggraph")
				os.Setenv("biggraph_outgoing_url", "http://[fe80::1%en0]/")
			},
			after: func() {
				os.Unsetenv("services")
				os.Unsetenv("biggraph_incoming_path")
				os.Unsetenv("biggraph_outgoing_url")
			},
			expectedErrLength: 1,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			logs = []string{}
			test.before()
			cfg := readInConfig()
			if test.expectedErrLength == 0 {
				assert.Equal(t, logs, []string{})
				assert.NotEqual(t, cfg, nil)
			} else {
				assert.Equal(t, len(logs), test.expectedErrLength)
			}
			test.after()
		})
	}
}

func TestCreateOutgoingURL(t *testing.T) {
	type Test struct {
		name              string
		config            proxyConfig
		incoming          *url.URL
		expectedError     string
		expectedOutputURL url.URL
	}

	example, _ := url.Parse("http://example.com")
	exampleIncoming, _ := url.Parse("/services/example/")
	exampleWithPort, _ := url.Parse("http://localhost:5000")
	exampleIncomingWithPath, _ := url.Parse("/services/example/test/1")
	exampleWithQuery, _ := url.Parse("/services/example/random?n=5")

	tests := []Test{
		Test{
			name: "able to set scheme and domain of URL",
			config: proxyConfig{
				incomingPath: "/services/example/",
				outgoingURL:  example,
				name:         "example",
			},
			incoming:      exampleIncoming,
			expectedError: "",
			expectedOutputURL: url.URL{
				Scheme: "http",
				Host:   "example.com",
				Path:   "",
			},
		},
		Test{
			name: "adds path back in successfully",
			config: proxyConfig{
				incomingPath: "/services/example/",
				outgoingURL:  example,
				name:         "example",
			},
			incoming:      exampleIncomingWithPath,
			expectedError: "",
			expectedOutputURL: url.URL{
				Scheme: "http",
				Host:   "example.com",
				Path:   "test/1",
			},
		},
		Test{
			name: "adds in port if present",
			config: proxyConfig{
				incomingPath: "/services/example/",
				outgoingURL:  exampleWithPort,
				name:         "example",
			},
			incoming:      exampleIncomingWithPath,
			expectedError: "",
			expectedOutputURL: url.URL{
				Scheme: "http",
				Host:   "localhost:5000",
				Path:   "test/1",
			},
		},
		Test{

			name: "supports encoding query strings",
			config: proxyConfig{
				incomingPath: "/services/example/",
				outgoingURL:  example,
				name:         "example",
			},
			incoming:      exampleWithQuery,
			expectedError: "",
			expectedOutputURL: url.URL{
				Scheme: "http",
				Host:   "example.com",
				Path:   "random?n=5",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			u := createOutgoingURL(test.config, test.incoming)
			assert.Equal(t, test.expectedError, "")
			assert.Equal(t, test.expectedOutputURL.Scheme, u.Scheme)
			assert.Equal(t, test.expectedOutputURL.Host, u.Host)
			assert.Equal(t, test.expectedOutputURL.Path, u.Path)

		})
	}
}
