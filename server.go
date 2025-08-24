package main

import (
	"log/slog"
	"net"
	"runtime"
	"sync"
)

var NOT_FOUND = []byte("HTTP/1.1 404 Not Found\r\nContent-Length: 0\r\n\r\n")
var FOUND = []byte("HTTP/1.1 302 Found\r\nContent-Length: 0\r\nLocation: ")

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
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		handleConnection(conn)
	}
}

var bufferPool = sync.Pool{
	New: func() interface{} {
		buf := new([1024]byte)
		return buf
	},
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	for {
		buf := bufferPool.Get().(*[1024]byte)
		n, err := conn.Read(buf[:]) // if client disconnects here it should eof/use of closed pipe so we're lowk chillin
		if err != nil {
			bufferPool.Put(buf)
			return
		}
		if n == 1024 {
			return // i don't want to deal w ts
		}
		activebuf := buf[:n]

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
			bufferPool.Put(buf)
			return
		}
		path := activebuf[start+1 : end]
		link := getlink(path[1:])
		bufferPool.Put(buf)

		// response
		resp := bufferPool.Get().(*[1024]byte)
		r := resp[:0]

		if link == nil {
			r = append(r, NOT_FOUND...)
		} else {
			r = append(r, FOUND...)
			r = append(r, link...)
			r = append(r, '\r', '\n', '\r', '\n')
		}

		_, err = conn.Write(r)
		bufferPool.Put(resp)

		if err != nil {
			return
		}
		// loop continues for next request (keep-alive, just assume cause why not :| )
	}
}
