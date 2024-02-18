package main

import (
	"encoding/csv"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
)

type Bond struct {
	short string
	long  string
}

func getLocalIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return nil
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP
}

var ServerLocalIP net.IP = getLocalIP()
var bonds []Bond

func getLongByShort(short string) (string, error) {
	for _, bond := range bonds {
		if bond.short == short {
			return bond.long, nil
		}
	}
	return "", fmt.Errorf("Bond with short %s not found", short)
}
func direct(w http.ResponseWriter, req *http.Request) {
	short := strings.Split(req.URL.Path, "/")[1]
	long, err := getLongByShort(short)
	if err != nil {
		fmt.Fprintf(w, "There is no go route for %s.", short)
		return
	}
	http.Redirect(w, req, long, http.StatusPermanentRedirect)
}

func main() {
	if ServerLocalIP == nil {
		fmt.Println("Unable to access local ip.")
		return
	}
	file, err := os.Open("bonds.csv")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer file.Close()
	reader := csv.NewReader(file)
	reader.FieldsPerRecord = 2
	data, err := reader.ReadAll()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	for _, row := range data {
		lng := row[1]
		lng = strings.Replace(lng, "SERVERIP", ServerLocalIP.String(), 1)
		bonds = append(bonds, Bond{row[0], lng})
	}
	http.HandleFunc("/", direct)
	fmt.Println(http.ListenAndServe(":80", nil))
}
