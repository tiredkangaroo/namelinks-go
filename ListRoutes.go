package main

import (
	"fmt"
	"net/http"
)

func ListRoutesHandler(w http.ResponseWriter, _ *http.Request) {
	var bds string = "<ul>"
	for _, bond := range bonds {
		bds += fmt.Sprintf("<li>%s  -->  <a href='%s'>%s</a></li>", bond.short, bond.long, bond.long)
		bds += "<br><br>"
	}
	bds += "</ul>"

	w.Header().Add("Content-Type", "text/html")
	fmt.Fprintf(w, bds)
}
