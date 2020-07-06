package main

import "net/http"

// LoggingMiddleware allows to log any request whenever the log level is equal to Debug.
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		sugar.Debugf("%s %s", r.Method, r.URL.EscapedPath())
		next.ServeHTTP(rw, r)
	})
}
