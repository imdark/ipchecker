package main

import (
	"fmt"
	"net"
	"time"
	"flag"
	"os"
	"bufio"
)


var ipsFile string
func init() {
	flag.StringVar(&ipsFile, "filename", "./sample_ips.csv", "file containing list of ips to scan")
}

const TIMEOUT = 5 * time.Second
type TCPCheckResult struct {
	ip string
	is_up bool
}

func DialTCP(tcp string, results chan<- TCPCheckResult) {
	_, error := net.DialTimeout("tcp", tcp, TIMEOUT)
	if error != nil {
		results <- TCPCheckResult{tcp, true}
	} else {
		results <- TCPCheckResult{tcp, false}
	}

}

type IpCheckTarget struct {
	targetTcpAddress net.TCPAddr
	targetIpMask net.IP
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
		// TODO: this should probably moved outside and processed in a diferent pipe
		targetIp, ipRange := ParseIpAndIpRange(line)
		targetTCP := GetConfiguredTcpAddressForIp(targetIp)
		lines <- IpCheckTarget{targetTCP, ipRange}
	}

	return inputLineByLineScanner.Err()
}

const IP_RANGE_BITS = 24



// 32 is the number of bits in ip v4

func ParseIpAndIpRange(ipString string) (net.IP, net.IP) {

	IP_RANGE_MASK := net.CIDRMask(IP_RANGE_BITS, 32)
	ipv4Addr := net.ParseIP(ipString)
        return ipv4Addr, ipv4Addr.Mask(IP_RANGE_MASK)
}

const STANDART_IP_CHECKER_PORT = 449
func GetConfiguredTcpAddressForIp(ip net.IP) net.TCPAddr {
	return net.TCPAddr{IP: ip, Port: STANDART_IP_CHECKER_PORT }
}

type TCPReport struct {
	TotalIps int
	TotalIpsNotReachable int
	IpRangesNotReachable int /*ranges are under /24 subnet masks */
	IpRangesPartiallyReachable int /*More then 50 ips in range are not*/
	TotalRunningTime int /* Including report generation */
}

func main() {
	fmt.Println(ipsFile)

        ipCheckTargets := make(chan IpCheckTarget)

	go ReadInputFile(ipsFile, ipCheckTargets)
	/*if err != nil {
		fmt.Println("Error reading input file", err)
	}*/

	//ips := []string{"127.0.0.1:449"}
        results := make(chan TCPCheckResult)
	ips_len := 0

	for ipCheckTarget := range ipCheckTargets {
                ips_len++
		go DialTCP(string(ipCheckTarget.targetTcpAddress.IP) + ":" + string(ipCheckTarget.targetTcpAddress.Port), results)
	}
	for i := 0; i < ips_len; i++ {
		fmt.Println(<-results)
	}
	fmt.Println("done")
}
