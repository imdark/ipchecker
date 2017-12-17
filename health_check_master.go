package main

import (
	"fmt"
	"net"
	"encoding/gob"
	"time"
	"os"
	"strconv"
)

func GetAgentReport(url string, reports chan<-HealthCheckAgentResult) {
	conn, err := net.Dial("tcp", url)
	if err != nil {
		reports<-HealthCheckAgentResult{nil, err}
		return
	}
	encoder := gob.NewEncoder(conn)
	agentReport := &TCPReport{}
	encoder.Encode(agentReport)
	conn.Close()

	reports<-HealthCheckAgentResult{agentReport, nil}

}

type CombinedTCPReport struct {
	TotalIps int
	TotalIpsNotReachable int
	IpRangesNotReachable int /*This is sum per agent, the code currently assumes that different agents have separate agents*/
	IpRangesPartiallyReachable int
	TotalRunningTime time.Duration /* Including report generation */
	AvrerageRunningTime time.Duration /*Average across all agents*/
	MaxTimePerAgent time.Duration /*Time on slowest agent*/
}

type HealthCheckAgentResult struct {
	TcpReport *TCPReport
	Err error
}
type TCPReport struct {
	TotalIps int
	TotalIpsNotReachable int
	IpRangesNotReachable int /*ip ranges are under /24 subnet mask */
	IpRangesPartiallyReachable int /*More then 50 ips in range are not*/
	TotalRunningTime time.Duration /* Including report generation */
}

func Max(x, y time.Duration) time.Duration {
	if x > y {
		return x
	} else {
		return y
	}
}

func ReadAgentsUrls() []string {
	ips_num_string := os.Getenv("agent_ips_num")
	if ips_num_string == "" {
		ips_num_string = "0"
	}

	ips_num, _ := strconv.Atoi(ips_num_string)

	agentUrls := make([]string, ips_num)
	for i := 0; i < ips_num; i++ {
		agentUrls[i] = os.Getenv("agent_ip" + strconv.Itoa(i + 1))
	}

	return agentUrls

}

func main() {
	agentUrls := ReadAgentsUrls()

	reports := make(chan HealthCheckAgentResult)

	agentUrlsLen := 0
	for _, agentUrl := range agentUrls {
		agentUrlsLen++
		go GetAgentReport(agentUrl, reports)
	}

	totalIps := 0
	totalIpsNotReachable := 0
	ipRangesNotReachable := 0
	ipRangesPartiallyReachable := 0
	var sumAgentsRunningTime time.Duration
	sumAgentsRunningTime = 0
	var maxTimePerAgent time.Duration
	maxTimePerAgent = 0


	for i := 0; i < agentUrlsLen; i++ {
		result := <-reports
		if result.Err != nil {
			fmt.Println("Error during connecting to client ", result.Err)
			continue
		}
		tcpReport := result.TcpReport

		totalIps += tcpReport.TotalIps
		totalIpsNotReachable += tcpReport.TotalIpsNotReachable
		ipRangesNotReachable += tcpReport.IpRangesNotReachable
		ipRangesPartiallyReachable += tcpReport.IpRangesPartiallyReachable
		sumAgentsRunningTime += tcpReport.TotalRunningTime
		maxTimePerAgent = Max(maxTimePerAgent, tcpReport.TotalRunningTime)
	}

	var avrageAgentsRunningTime time.Duration
	avrageAgentsRunningTime = time.Duration(int64(sumAgentsRunningTime) / int64(agentUrlsLen))
	combinedReport := CombinedTCPReport{
		totalIps,
		totalIpsNotReachable,
		ipRangesNotReachable,
		ipRangesPartiallyReachable,
		sumAgentsRunningTime,
		avrageAgentsRunningTime,
		maxTimePerAgent}

	fmt.Printf("%+v", combinedReport)
}
