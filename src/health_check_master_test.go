package main

import (
	"testing"
	"libs"
)

func TestAgentExistingIps(t *testing.T) {
	libs.GenerateReport("sample_ips.csv")
}

func TestAgentMixedIps(t *testing.T) {
	libs.GenerateReport("mixed_ips.csv")
}

func TestSummeryCombination(t *testing.T) {

}


func TestSummeryGeneration(t *testing.T) {

}

func TestIpParsing(t *testing.T) {

}
