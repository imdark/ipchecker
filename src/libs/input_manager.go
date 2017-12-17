package libs

import (
	"os"
	"bufio"
	"net"
)
type IpCheckTarget struct {
	TargetTcpAddress net.TCPAddr
	TargetIpMask string
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
