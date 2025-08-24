package main

import (
	"fmt"
	"log/slog"
	"net"
	"runtime"
	"sync"
)

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
		return make([]byte, 0, 512)
	},
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	buf := bufferPool.Get().([]byte)
	defer bufferPool.Put(buf)

	n, err := conn.Read(buf[:512])
	if err != nil {
		return
	}
	activebuf := buf[:n]

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
	fmt.Println(string(path))
	return
}

func main() {
	serve()
}
