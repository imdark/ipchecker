package main

import (
	"testing"
	"github.com/imdark/ipchecker/src/libs"
	"net"
	"bytes"
//	"time"
)

/*func TestAgentExistingIps(t *testing.T) {
	expectedReport := libs.TCPReport{50003, 45640, 1, 1, time.Minute}
	actualReport := libs.GenerateReport("sample_ips.csv")

	if expectedReport.TotalIps != actualReport.TotalIps ||
	   expectedReport.TotalRunningTime <= actualReport.TotalRunningTime {
		t.Errorf("Generated report was incorrect, got: %v, want: %v.", actualReport, expectedReport)
	}
}

func TestAgentMixedIps(t *testing.T) {
	expectedReport := libs.TCPReport{50103, 49413, 200, 49, time.Minute}
	actualReport := libs.GenerateReport("mixed_ips.csv")

	if expectedReport.TotalIps != actualReport.TotalIps ||
	   expectedReport.TotalRunningTime <= actualReport.TotalRunningTime {
		t.Errorf("Generated report was incorrect, got: %v, want: %v.", actualReport, expectedReport)
	}
}*/

func TestParseIp(t *testing.T) {
	expectedIp := net.IPv4(127, 0, 0, 1)
	expectedIpMask := "127.0.0.0"
	actualIp, actualIpMask := libs.ParseIpAndIpRange("127.0.0.1")
	if !bytes.Equal(expectedIp, actualIp) || expectedIpMask != actualIpMask {
		t.Errorf("IP was incorrect, got ip: %v, want ip: %v.", actualIp, expectedIp)
	}

	if expectedIpMask != actualIpMask {
		t.Errorf("IP mask was incorrect, got ip mask: %v, want ip mask: %v.", actualIpMask, expectedIpMask)
	}
}

func TestCountIpRangesFrequencies(t *testing.T) {
	expectedIpRangesNotReachable := 1
	expectedIpRangesPartiallyReachable := 2

	ipRangeFrequencyMap := make(map[string]*libs.IpRangeFrequncy)
	ipRangeFrequencyMap["127.0.0.1"] = &libs.IpRangeFrequncy{10, 20}
	ipRangeFrequencyMap["127.0.1.1"] = &libs.IpRangeFrequncy{20, 20}
	ipRangeFrequencyMap["127.0.2.1"] = &libs.IpRangeFrequncy{51, 100}
	ipRangeFrequencyMap["127.0.3.1"] = &libs.IpRangeFrequncy{51, 100}

	actualIpRangesNotReachable, actualIpRangesPartiallyReachable :=
		libs.CountIpRangesFrequencies(ipRangeFrequencyMap)

	if actualIpRangesNotReachable != expectedIpRangesNotReachable {
		t.Errorf("IpRangesNotReachable was incorrect, " +
		"got ipRangesNotReachable: %v, want ipRangesNotReachable: %v.",
		actualIpRangesNotReachable, expectedIpRangesNotReachable)
	}

	if actualIpRangesPartiallyReachable != expectedIpRangesPartiallyReachable {
		t.Errorf("IpRangesPartiallyReachable was incorrect, " +
		"got IpRangesPartiallyReachable : %v, want IpRangesPartiallyReachable: %v.",
		actualIpRangesPartiallyReachable, expectedIpRangesPartiallyReachable)
	}
}
