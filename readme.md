

Design Goals
------------
50,000 IPs every minute

To Run:
-------

```
cd $GOPATH
git clone https://github.com/imdark/ipchecker

export TCP_PORT=8082
go run health_check_agent/health_check_agent.go <optional input/file/path>

export agent_ips_num=1
export agent_ip1=localhost:8082
go run health_check_master/health_check_master.go

```

code assumes input file is a list of ips seperated by break lines
if no parameters are specified the code will try to run on sample_ips.csv

To Test:
--------
```go test```

Next steps:
-----------
[] ip parsing
[] getting a list of 5000 ips

[] app currently assumes that list of ips will fit into memory
[] added the linux environment variable you need to change
[] currently the code assumes all ips read from input files are ipv4

