package main

import (
	"fmt"
	"net/http"
	"strings"
)

func RedirectionHandler(w http.ResponseWriter, req *http.Request) {
	sp := strings.Split(req.URL.Path, "/")
	short := sp[1]
	rest := "/" + strings.Join(sp[2:], "/")
	long, err := GetLong(short)
	if err != nil {
		fmt.Fprintf(w, "There is no go route for %s.", short)
		return
	}
	long += rest
	http.Redirect(w, req, long, http.StatusPermanentRedirect)
}
