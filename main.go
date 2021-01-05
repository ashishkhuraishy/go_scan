package main

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

var wg sync.WaitGroup

// MaxRoutines store the max number
// of goroutines that are spawned
const MaxRoutines = 1000

func main() {
	ScanAll("localhost", 5*time.Second)
	wg.Wait()
}

// ScanAll will scan all ports
func ScanAll(addr string, timeout time.Duration) {

	gaurd := make(chan int, MaxRoutines)

	for i := 0; i < 65000; i++ {
		wg.Add(1)
		gaurd <- 1
		go func(n int) {
			ScanPort(addr, n, timeout)
			<-gaurd
		}(i)

	}

}

// QuickScan scans through the most common known ports
func QuickScan(addr string, timeout time.Duration) {
	commonPorts := map[int]string{
		7:    "echo",
		20:   "ftp",
		21:   "ftp",
		22:   "ssh",
		23:   "telnet",
		25:   "smtp",
		43:   "whois",
		53:   "dns",
		67:   "dhcp",
		68:   "dhcp",
		80:   "http",
		110:  "pop3",
		123:  "ntp",
		137:  "netbios",
		138:  "netbios",
		139:  "netbios",
		143:  "imap4",
		443:  "https",
		513:  "rlogin",
		540:  "uucp",
		554:  "rtsp",
		587:  "smtp",
		873:  "rsync",
		902:  "vmware",
		989:  "ftps",
		990:  "ftps",
		1194: "openvpn",
		3306: "mysql",
		5000: "unpn",
		8080: "https-proxy",
		8443: "https-alt",
	}

	for i := range commonPorts {
		wg.Add(1)
		go ScanPort(addr, i, timeout)
	}
}

// ScanPort -> to write
func ScanPort(addr string, port int, timeout time.Duration) {
	defer wg.Done()
	address := addr + ":" + strconv.Itoa(port)

	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		// fmt.Println("Error while connecting", err.Error())
		if strings.Contains(err.Error(), "too may open files") {
			time.Sleep(timeout)
			ScanPort(addr, port, timeout)
		}
		return
	}
	defer conn.Close()

	fmt.Println(port, "open")
	return
}
