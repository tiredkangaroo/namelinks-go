package main

import (
	"fmt"
	"net"
	"net/http"
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

func getLongByShort(bonds []Bond, short string) (string, error) {
	for _, bond := range bonds {
		if bond.short == short {
			return bond.long, nil
		}
	}
	return "", fmt.Errorf("Bond with short %s not found", short)
}
func direct(w http.ResponseWriter, req *http.Request) {
	bonds := []Bond{
		{"pihole", fmt.Sprintf("http://%s:8080/admin", ServerLocalIP)},
	}
	short := strings.Split(req.URL.Path, "/")[1]
	long, err := getLongByShort(bonds, short)
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
	fmt.Println(ServerLocalIP)
	http.HandleFunc("/", direct)
	fmt.Println(http.ListenAndServe(":80", nil))
}
