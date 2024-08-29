package main

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/valyala/fasthttp"
)

type bonds *map[string]string

func notFound(c *fasthttp.RequestCtx) {
	c.Response.Header.Set("Content-Type", "text/html")
	fmt.Fprintf(c, `
    <h1> Not Found </h1>
    <pre>There is no registry for the short your provided.</pre>
  `)
}

func listPaths(c *fasthttp.RequestCtx, b bonds) {
	cmd := exec.Command("git", "rev-parse", "HEAD")
	output, err := cmd.CombinedOutput()
	version := ""
	if err == nil {
		version = string(output)
	} else {
		version = fmt.Sprintf("There was an issue reading the version: %s.", err.Error())
	}

	listElements := ""
	for short, long := range *b {
		listElements += fmt.Sprintf(`<li>%s: <a href="%s">%s</a></li><br>`, short, long, long)
	}
	c.Response.Header.Set("Content-Type", "text/html")
	fmt.Fprintf(c, `
	<p>Version: %s</p>
    <ul>
      %s
    </ul>
  `, version, listElements)
}

func handleRedirection(c *fasthttp.RequestCtx, b bonds, path string) {
	s := strings.Split(path, "/")
	shortReq := strings.ToLower(s[1])
	longFromReq := "/not-found"
	for short, long := range *b {
		if short == shortReq {
			longFromReq = long
			break
		}
	}
	if len(s) > 2 && longFromReq != "/not-found" {
		longFromReq += "/" + strings.Join(s[2:], "/")
	}
	c.Redirect(longFromReq, 301)
}

// startServer starts the FastHTTP server at port 80.
func startServer(b *map[string]string) {
	handler := func(c *fasthttp.RequestCtx) {
		path := string(c.Path())
		switch path {
		case "/not-found":
			notFound(c)
		case "/list":
			listPaths(c, b)
		default:
			handleRedirection(c, b, path)
		}
	}
	fasthttp.ListenAndServe(":80", handler)
}
