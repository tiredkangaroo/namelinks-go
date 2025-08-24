package main

import (
	"log/slog"
	"net"
	"runtime"
)

// server documentation:
// - GET /<name> -> will either redirect to the corresponding link or return a 404 error
// - PATCH / -> will update the link for the given name (drops connection, cannot Keep-Alive here)
// this is a fragile handler, it assumes the request is well-formed and valid and does not handle errors gracefully.

var NOT_FOUND = []byte("HTTP/1.1 404 Not Found\r\nContent-Length: 0\r\n\r\n")

var FOUND_PREFIX = []byte("HTTP/1.1 302 Found\r\nContent-Length: 0\r\nLocation: ")
var crlfcrlf = []byte("\r\n\r\n")

func serve() error {
	listener, err := net.Listen("tcp", ":80")
	if err != nil {
		return err
	}
	defer listener.Close()

	slog.Info("starting workers", "count", runtime.NumCPU())
	for range runtime.NumCPU() {
		go worker(listener)
	}

	select {}
}

func worker(listener net.Listener) {
	runtime.LockOSThread()
	buf := new([1024]byte) // pre-allocate buffer
	respBuf := new([1024]byte)
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		handleConnection(conn, buf, respBuf)
	}
}

func handleConnection(conn net.Conn, buf *[1024]byte, respBuf *[1024]byte) {
	defer conn.Close()

	for {
		n, err := conn.Read(buf[:]) // if client disconnects here it should eof/use of closed pipe so we're lowk chillin
		if err != nil {
			return
		}
		// t := time.Now()
		if n == 1024 || n == 0 {
			return // i don't want to deal w ts
		}
		activebuf := buf[:n]
		if activebuf[0] != 'G' { // not a GET request, client is trying to indicate that they want to refresh
			load()
			return
		}

		// parse path
		spaceCount := 0
		var start, end int
		for i, b := range activebuf {
			if b != ' ' {
				continue
			}
			spaceCount++
			if spaceCount == 1 {
				start = i
			} else if spaceCount == 2 {
				end = i
				break
			}
		}
		if start == 0 || end == 0 || start+1 > len(activebuf) {
			return
		}
		path := activebuf[start+1 : end]
		link := getlink(path[1:])

		if link == nil {
			_, err = conn.Write(NOT_FOUND)
		} else {
			pos := copy(respBuf[:], FOUND_PREFIX)
			pos += copy(respBuf[pos:], link)
			pos += copy(respBuf[pos:], crlfcrlf)
			_, err = conn.Write(respBuf[:pos])
		}
		if err != nil {
			return
		}

		// s := time.Since(t)
		// slog.Info("served", "link", string(link), "duration_microseconds", s.Microseconds())
		// loop continues for next request (keep-alive, just assume cause why not :| )
	}
}
