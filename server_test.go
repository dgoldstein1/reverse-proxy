package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
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
