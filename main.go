package main

import "fmt"
import "net"
import "time"

const TIMEOUT = 5 * time.Second
type TCPCheckResult struct {
	ip string
	is_up bool
}

func DialIp(ip string, results chan<- TCPCheckResult) {
	_, error := net.DialTimeout("tcp", ip, TIMEOUT)
	if error != nil {
		results <- TCPCheckResult{ip, true}
	} else {
		results <- TCPCheckResult{ip, false}
	}

}

const PORT = 449

func ParseIps() {
}

const IP_RANGE_BITS = 24

// 32 is the number of bits in ip v4
const IP_RANGE_MASK = net.CIDRMask(IP_RANGE_BITS, 32)

func ParseIpAndIpRange(ipString) IP, IP {

	ipv4Addr := net.ParseIP("192.0.2.1")
        return ipv4Addr, ipv4Addr.Mask(IP_RANGE_MASK)
}

type TCPReport struct {
	TotalIps int
	TotalIpsNotReachable int
	IpRangesNotReachable int /*ranges are under /24 subnet masks */
	IpRangesPartiallyReachable int /*More then 50 ips in range are not*/
	TotalRunningTime int /* Including report generation */
}
func main() {
	//TCPAddr{IP: ip.IP, Port: portnum}
	ips := []string{"127.0.0.1:449"}
        results := make(chan TCPCheckResult)
	for _, ip := range ips {
		go DialIp(ip, results)
	}
	ips_len := len(ips)
	for i := 0; i < ips_len; i++ {
		fmt.Println(<-results)
	}
	fmt.Println("done")
}
