package main

import (
	_ "embed"
	"fmt"
	"net"

	"github.com/redis/go-redis/v9"
)

//go:embed version.txt
var version string

func getLocalIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return nil
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP
}

func main() {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	go startServer(client)
	fmt.Println(`Visit: http://[::1]:80`)
	select {}
}
