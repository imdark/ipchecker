package main

import (
	"testing"
	"github.com/imdark/ipchecker/src/libs"
)

func TestAgentExistingIps(t *testing.T) {
	libs.GenerateReport("sample_ips.csv")
}

func TestAgentMixedIps(t *testing.T) {
	libs.GenerateReport("mixed_ips.csv")
}

