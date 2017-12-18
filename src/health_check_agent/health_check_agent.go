package main

import (
	"time"
	"github.com/imdark/ipchecker/src/libs"
	"fmt"
	"net"
	"flag"
	"encoding/gob"
)


var ipsFile string
func init() {
	flag.StringVar(&ipsFile, "filename", "./sample_ips.csv", "file containing list of ips to scan")
}

func HealthCheckConfiguredIps(conn net.Conn) {
	encoder := gob.NewEncoder(conn)
	fmt.Println("before tcpReport")
	//tcpReport := libs.GenerateReport(ipsFile)
	tcpReport := libs.TCPReport{1,2,3,4, time.Second}
	fmt.Println("after tcpReport")
	fmt.Println("tcpReport", tcpReport)
	encoder.Encode(tcpReport)
	conn.Close()
}

func main() {
	flag.Parse()
	configuredPort := libs.GetConfiguredPort()
	listener, err := net.Listen("tcp", ":" + configuredPort)
	if err != nil {
		fmt.Println("Error opening http server on port ", configuredPort, err)
		return
	}
	for {
		connection, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting", err)
			return
		}

		time.Sleep(time.Second * 50)
		HealthCheckConfiguredIps(connection)
	}
}
