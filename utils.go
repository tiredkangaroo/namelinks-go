package main

import (
	"encoding/csv"
	"fmt"
	"net"
	"os"
	"strings"
)

func getLocalIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return nil
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP
}

func GetLong(short string) (string, error) {
	for _, bond := range bonds {
		if bond.short == short {
			return bond.long, nil
		}
	}
	return "", fmt.Errorf("Bond with short %s not found", short)
}

func CreateBonds() error {
	file, err := os.Open("bonds.csv")
	if err != nil {
		return fmt.Errorf("Opening bonds.csv failed: %s", err.Error())
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = 2
	data, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("Reading CSV from bonds.csv failed with: %s", err.Error())
	}

	for _, row := range data {
		lng := row[1]
		lng = strings.Replace(lng, "SERVERIP", ServerIP.String(), 1)
		bonds = append(bonds, Bond{row[0], lng})
	}
	return nil
}

var ServerIP = getLocalIP()
