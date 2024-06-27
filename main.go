package main

import (
	"fmt"
	"net/http"
)

type Bond struct {
	short string
	long  string
}

var bonds []Bond

func main() {
	if ServerIP == nil {
		fmt.Println("Accessing Local IP failed.")
		return
	}

	err := CreateBonds()
	if err != nil {
		fmt.Println(err)
		return
	}

	http.HandleFunc("/list", ListRoutesHandler)
	http.HandleFunc("/", RedirectionHandler)

	fmt.Println(http.ListenAndServe(":80", nil))
}
