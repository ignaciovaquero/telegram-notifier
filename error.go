package main

import (
	"fmt"
	"io"
	"net/http"
)

// HandleError This function returns an error to the caller, along with a 400 Bad Request
// status code.
func HandleError(err error, rw http.ResponseWriter) bool {
	return HandleErrorWithStatusCode(err, rw, http.StatusBadRequest)
}

// HandleErrorWithStatusCode This function returns an error to the caller, along with a
// custom status code.
func HandleErrorWithStatusCode(err error, rw http.ResponseWriter, status int) bool {
	if err != nil {
		rw.WriteHeader(status)
		_, er2 := io.WriteString(rw, fmt.Sprintf("%s\n", err.Error()))
		if er2 != nil {
			rw.WriteHeader(http.StatusInternalServerError)
		}
		return true
	}
	return false
}
