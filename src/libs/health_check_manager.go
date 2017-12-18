package libs
import (
	"time"
	"strconv"
	"net"
)


const TIMEOUT = 5 * time.Second
func DialTCP(tcp string, targetIpMask string, results chan<- TCPCheckResult) {
	_, err := net.DialTimeout("tcp", tcp, TIMEOUT)
	if err == nil {
		results <- TCPCheckResult{targetIpMask, tcp, true}
	} else {
		results <- TCPCheckResult{targetIpMask, tcp, false}
	}
}

type IpRangeFrequncy struct {
	IpsNotReachable int
	IpsInRange int
}

type TCPCheckResult struct {
	TargetIpMask string
	Ip string
	IsUp bool
}

func GenerateReport(ipsFile string) *TCPReport {
	startTime := time.Now()
        ipCheckTargets := make(chan IpCheckTarget)

	go ReadInputFile(ipsFile, ipCheckTargets)

        results := make(chan TCPCheckResult)
	ipsCount := 0

	for ipCheckTarget := range ipCheckTargets {
                ipsCount++
		targetIp := ipCheckTarget.TargetTcpAddress.IP.String()
		targetPort := strconv.Itoa(ipCheckTarget.TargetTcpAddress.Port)
		targetAddress := targetIp  + ":" + targetPort
		go DialTCP(targetAddress, ipCheckTarget.TargetIpMask, results)
	}

	ipRangeFrequencyMap, totalIpsNotReachable := CountIps(ipsCount, results)
	ipRangesNotReachable, ipRangesPartiallyReachable := CountIpRangesFrequencies(ipRangeFrequencyMap)
	endTime := time.Now()

	return &TCPReport{TotalIps: ipsCount,
			      TotalIpsNotReachable: totalIpsNotReachable,
			      IpRangesNotReachable: ipRangesNotReachable,
			      IpRangesPartiallyReachable: ipRangesPartiallyReachable,
			      TotalRunningTime: endTime.Sub(startTime)}
}

func CountIps(ipsCount int, results chan TCPCheckResult) (map[string]*IpRangeFrequncy, int) {
	ipRangeFrequencyMap := make(map[string]*IpRangeFrequncy)
	totalIpsNotReachable := 0

	for i := 0; i < ipsCount; i++ {
		healthCheckResult := <-results
		ipRangeFrequncy, ok := ipRangeFrequencyMap[healthCheckResult.TargetIpMask]
		if !ok {
			ipRangeFrequncy = &IpRangeFrequncy{}
			ipRangeFrequencyMap[healthCheckResult.TargetIpMask] = ipRangeFrequncy
		}
		ipRangeFrequncy.IpsInRange++

		if !healthCheckResult.IsUp {
			totalIpsNotReachable++
			ipRangeFrequncy.IpsNotReachable++
		}
	}
	return ipRangeFrequencyMap, totalIpsNotReachable
}

const IP_RANGE_THRUSHSHOLD = 50
func CountIpRangesFrequencies(ipRangeFrequencyMap map[string]*IpRangeFrequncy) (int, int) {
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
	return ipRangesNotReachable, ipRangesPartiallyReachable
}
