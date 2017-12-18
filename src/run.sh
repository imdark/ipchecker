#!/bin/bash
export TCP_PORT=8082
go run health_check_agent/health_check_agent.go

export agent_ips_num=1
export agent_ip1=localhost:8082
go run health_check_master/health_check_master.go

