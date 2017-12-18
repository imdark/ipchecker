package main

import (
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
	tcpReport := libs.GenerateReport(ipsFile)
	fmt.Println("request recived", tcpReport)
	encoder.Encode(tcpReport)
	conn.Close()
}

func main() {
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

		go HealthCheckConfiguredIps(connection)
	}
}
