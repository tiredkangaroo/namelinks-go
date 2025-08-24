package main

import (
	"bytes"
	"hash/crc32"
	"log/slog"
	"net"
	"os"
)

var namelinks = make([][]byte, 256*256) // 65536 slots of name links :)
var serverIP []byte

func init() {
	conn, err := net.Dial("tcp", "1.1.1.1:80")
	if err != nil {
		slog.Error("could not determine server IP", "err", err.Error())
		os.Exit(1)
	}
	defer conn.Close()
	serverIP = []byte(conn.LocalAddr().(*net.TCPAddr).IP.String()) // don't add port to serverIP
}

func hash16crc(b []byte) uint16 { // collision generator
	return uint16(crc32.ChecksumIEEE(b))
}

func getlink(name []byte) []byte {
	h := hash16crc(name)
	if link := namelinks[h]; link != nil {
		return link
	}
	if link, err := getAILink(name); err == nil {
		putlink(name, link)
		return link
	}
	return nil
}

func putlink(name, link []byte) {
	h := hash16crc(name)
	if namelinks[h] != nil {
		panic("hash collision for " + string(name))
	}
	namelinks[h] = link
}

func load() {
	data, err := os.ReadFile("namelinks.txt")
	if err != nil {
		slog.Error("namelinks.txt read error", "err", err.Error())
		return
	}
	for _, line := range bytes.Split(data, []byte{'\n'}) {
		parts := bytes.SplitN(line, []byte{' '}, 2)
		if len(parts) != 2 {
			continue
		}
		name := parts[0]
		link := parts[1]
		link = bytes.ReplaceAll(link, []byte("$SERVER_IP"), serverIP)
		putlink(name, link)
		slog.Info("loaded", "name", string(name), "link", string(link))
	}
}
