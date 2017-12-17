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

