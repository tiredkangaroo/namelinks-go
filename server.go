package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/redis/go-redis/v9"
	"github.com/valyala/fasthttp"
)

func notFound(c *fasthttp.RequestCtx) {
	c.Response.Header.Set("Content-Type", "text/html")
	fmt.Fprintf(c, `
    <h1> Not Found </h1>
    <pre>There is no registry for the short your provided.</pre>
  `)
}

func listPaths(c *fasthttp.RequestCtx, client *redis.Client) {
	listElements := ""

	res := client.HGetAll(context.Background(), "urlshortner")
	b := res.Val()

	for short, long := range b {
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

func handleRedirection(c *fasthttp.RequestCtx, client *redis.Client, path string) {
	s := strings.Split(path, "/")
	shortReq := strings.ToLower(s[1])
	long := ""

	res := client.HGet(c, "urlshortner", shortReq)
	if res == nil {
		c.Redirect("/not-found", 301)
		return
	}
	if res.Err() != nil {
		c.Redirect("/not-found", 301)
		return
	}
	long = res.Val()
	if long == "" {
		c.Redirect("/not-found", 301)
		return
	}
	long = strings.ReplaceAll(long, "$SERVERIP", getLocalIP().String())

	if len(s) > 2 {
		long += "/" + strings.Join(s[2:], "/")
	}
	c.Redirect(long, 301)
}

// startServer starts the FastHTTP server at port 80.
func startServer(client *redis.Client) {
	handler := func(c *fasthttp.RequestCtx) {
		c.Response.Header.Set("Cache-Control", "no-store")
		path := string(c.Path())
		switch path {
		case "/not-found":
			notFound(c)
		case "/list":
			listPaths(c, client)
		default:
			handleRedirection(c, client, path)
		}
	}
	fasthttp.ListenAndServe(":80", handler)
}
