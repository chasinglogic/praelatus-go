// Copyright 2017 Mathew Robinson <mrobinson@praelatus.io>. All rights reserved.
// Use of this source code is governed by the AGPLv3 license that can be found in
// the LICENSE file.

package middleware

import (
	"log"
	"net/http"
	"time"

	"github.com/praelatus/praelatus/config"
)

// LoggedResponseWriter wraps http.ResponseWriter so we can capture the status
// code for logging
type LoggedResponseWriter struct {
	status int
	error  []byte
	http.ResponseWriter
}

// Status will return the status code used in this request.
func (w *LoggedResponseWriter) Status() int {
	return w.status
}

// WriteHeader implements http.ResponseWriter adding our custom functionality
// to it
func (w *LoggedResponseWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *LoggedResponseWriter) Write(b []byte) (int, error) {
	if w.status == 0 {
		w.WriteHeader(200)
	} else if w.status >= 400 {
		w.error = b
	}

	return w.ResponseWriter.Write(b)
}

// Logger will log a request and any information about the request, it should
// be the first middleware in any chain.
func Logger(next http.Handler) http.Handler {
	var requestLog = log.New(config.LogWriter(), "[REQUEST] ", log.LstdFlags)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		lrw := &LoggedResponseWriter{status: 0, ResponseWriter: w}
		next.ServeHTTP(lrw, r)

		if lrw.error != nil {
			requestLog.Printf("|%s| [%d] %s %s %s",
				r.Method, lrw.Status(), r.URL.Path, time.Since(start).String(), string(lrw.error))
		} else {
			requestLog.Printf("|%s| [%d] %s %s",
				r.Method, lrw.Status(), r.URL.Path, time.Since(start).String())
		}
	})
}
