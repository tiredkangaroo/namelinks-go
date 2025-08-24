package main

import (
	"bytes"
	"hash/crc32"
	"log/slog"
	"os"
)

var namelinks = make([][]byte, 256*256) // 65536 slots of name links :)

func hash16crc(b []byte) uint16 { // collision generator
	return uint16(crc32.ChecksumIEEE(b))
}

func getlink(name []byte) []byte {
	h := hash16crc(name)
	return namelinks[h]
}

func putlink(name, link []byte) {
	h := hash16crc(name)
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
		putlink(name, link)
		slog.Info("loaded", "name", string(name), "link", string(link))
	}
}
