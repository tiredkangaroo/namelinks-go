package main

import (
	"fmt"
	"net"
	"os"
	"strings"
	"time"
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

func loadBonds() (bonds, error) {
	data, err := os.ReadFile("bonds.txt")
	if err != nil {
		return nil, err
	}
	localIP := getLocalIP().String()
	dataLines := strings.Split(string(data), "\n")
	bp := make(map[string]string)
	for i, line := range dataLines {
		sl := strings.Split(line, " ")
		if len(sl) == 1 { // empty line
			continue
		}
		if len(sl) < 2 {
			return nil, fmt.Errorf("bonds syntax error: invalid row at line %d in bonds.txt", i)
		}
		sl1 := strings.ReplaceAll(sl[1], "$SERVERIP", localIP)
		bp[sl[0]] = sl1
	}
	return &bp, nil
}

func reloadBonds(b bonds) {
	for {
		time.Sleep(time.Second * 60)
		newB, err := loadBonds()
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		b = newB
	}
}

func main() {
	b, err := loadBonds()
	if err != nil {
		fmt.Println(err)
		return
	}

	go reloadBonds(b)
	go startServer(b)
	fmt.Println(`
  __                 __             .___
  _______/  |______ ________/  |_  ____   __| _/   ______ ______________  __ ___________
 /  ___/\   __\__  \\_  __ \   __\/ __ \ / __ |   /  ___// __ \_  __ \  \/ // __ \_  __ \
 \___ \  |  |  / __ \|  | \/|  | \  ___// /_/ |   \___ \\  ___/|  | \/\   /\  ___/|  | \/
/____  > |__| (____  /__|   |__|  \___  >____ |  /____  >\___  >__|    \_/  \___  >__|
     \/            \/                 \/     \/       \/     \/                 \/
        __                          __      ______ _______
_____ _/  |_  ______   ____________/  |_   /  __  \\   _  \
\__  \\   __\ \____ \ /  _ \_  __ \   __\  >      </  /_\  \
 / __ \|  |   |  |_> >  <_> )  | \/|  |   /   --   \  \_/   \
(____  /__|   |   __/ \____/|__|   |__|   \______  /\_____  /
     \/       |__|                               \/       \/

     Visit: http://[::1]:80
	`)
	select {}
}
