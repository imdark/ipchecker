package main

import (
	"fmt"
	"net"
	"flag"
	"os"
	"bufio"
	"time"
	"strconv"
	"encoding/gob"
)


var ipsFile string
func init() {
	flag.StringVar(&ipsFile, "filename", "./sample_ips.csv", "file containing list of ips to scan")
}

const TIMEOUT = 5 * time.Second
type TCPCheckResult struct {
	targetIpMask string
	ip string
	is_up bool
}

func DialTCP(tcp string, targetIpMask string, results chan<- TCPCheckResult) {
	_, err := net.DialTimeout("tcp", tcp, TIMEOUT)
	if err == nil {
		results <- TCPCheckResult{targetIpMask, tcp, true}
	} else {
		results <- TCPCheckResult{targetIpMask, tcp, false}
	}
}

type IpCheckTarget struct {
	targetTcpAddress net.TCPAddr
	targetIpMask string
}

// code assumes input is a stored in a file the program has permission to read
// and its content is a list of ipv4s seperated by line breakers, not dns names 
func ReadInputFile(filePath string, lines chan<- IpCheckTarget) (err error) {

	defer close(lines)
	inputFile, err := os.Open(filePath)
	if err != nil {

		return err
	}

	// should only happen once the function is done, and unless the file 
	// could not be opened
	defer inputFile.Close()


	inputLineByLineScanner := bufio.NewScanner(inputFile)
	// this is an optimization, the list of ips file might be too big so
	// we start reading a row and pinging one ip at a time 

	for inputLineByLineScanner.Scan() {
		line := inputLineByLineScanner.Text()
		// TODO: this should probably moved outside and processed in a seprate pipe
		targetIp, ipRange := ParseIpAndIpRange(line)
		targetTCP := GetConfiguredTcpAddressForIp(targetIp)
		lines <- IpCheckTarget{targetTCP, ipRange}
	}

	return inputLineByLineScanner.Err()
}

const IP_RANGE_BITS = 24

// 32 is the number of bits in ip v4

func ParseIpAndIpRange(ipString string) (net.IP, string) {

	IP_RANGE_MASK := net.CIDRMask(IP_RANGE_BITS, 32)
	ipv4Addr := net.ParseIP(ipString)

        return ipv4Addr, ipv4Addr.Mask(IP_RANGE_MASK).String()
}

const STANDART_IP_CHECKER_PORT = 80
func GetConfiguredTcpAddressForIp(ip net.IP) net.TCPAddr {
	return net.TCPAddr{IP: ip, Port: STANDART_IP_CHECKER_PORT }
}

type TCPReport struct {
	TotalIps int
	TotalIpsNotReachable int
	IpRangesNotReachable int /*ip ranges are under /24 subnet mask */
	IpRangesPartiallyReachable int /*More then 50 ips in range are not*/
	TotalRunningTime time.Duration /* Including report generation */
}

type IpRangeFrequncy struct {
	IpsNotReachable int
	IpsInRange int
}

const IP_RANGE_THRUSHSHOLD = 50
func GetConfiguredPort() string {
	configuredPort := os.Getenv("TCP_PORT")
	if configuredPort != "" {
		return configuredPort
	} else {
		return "8081"
	}
}

func main() {
	fmt.Println("")
	//TODO: move port to configuration
	configuredPort := GetConfiguredPort()
	listener, err := net.Listen("tcp", ":" + configuredPort)
	if err != nil {
		fmt.Println("Error opening http server on port ", configuredPort, err)
		return
	}
	for {
		connection, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting", err)
		}

		go GenerateReport(connection)
	}
}

func GenerateReport(conn net.Conn) {
	startTime := time.Now()
        ipCheckTargets := make(chan IpCheckTarget)

	go ReadInputFile(ipsFile, ipCheckTargets)
	/*if err != nil {
		fmt.Println("Error reading input file", err)
	}*/

        results := make(chan TCPCheckResult)
	ips_len := 0

	for ipCheckTarget := range ipCheckTargets {
                ips_len++
		go DialTCP(ipCheckTarget.targetTcpAddress.IP.String() + ":" + strconv.Itoa(ipCheckTarget.targetTcpAddress.Port), ipCheckTarget.targetIpMask, results)
	}
	ipRangeFrequencyMap := make(map[string]*IpRangeFrequncy)
	totalIpsNotReachable := 0

	for i := 0; i < ips_len; i++ {
		health_check_result := <-results


		ipRangeFrequncy, ok := ipRangeFrequencyMap[health_check_result.targetIpMask]
		if !ok {
			ipRangeFrequncy = &IpRangeFrequncy{}


			ipRangeFrequencyMap[health_check_result.targetIpMask] = ipRangeFrequncy
		}
		ipRangeFrequncy.IpsInRange++

		if !health_check_result.is_up {
			totalIpsNotReachable++
			ipRangeFrequncy.IpsNotReachable++
		}
	}

	ipRangesNotReachable := 0
	ipRangesPartiallyReachable := 0
	for ipRange := range ipRangeFrequencyMap {
		item := ipRangeFrequencyMap[ipRange]
		if item.IpsNotReachable == item.IpsInRange {
			ipRangesNotReachable++
		} else if item.IpsNotReachable > IP_RANGE_THRUSHSHOLD {
			ipRangesPartiallyReachable++
		}
	}

	endTime := time.Now()

	encoder := gob.NewEncoder(conn)
	tcpReport := &TCPReport{TotalIps: ips_len,
			      TotalIpsNotReachable: totalIpsNotReachable,
			      IpRangesNotReachable: ipRangesNotReachable,
			      IpRangesPartiallyReachable: ipRangesPartiallyReachable,
			      TotalRunningTime: endTime.Sub(startTime)}
	encoder.Encode(tcpReport)
	conn.Close()
}
