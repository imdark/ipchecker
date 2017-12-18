

Design Goals
------------
50,000 IPs every minute

To Run:
-------

Run child agent
```
go get github.com/imdark/ipchecker/src/health_check_agent

export TCP_PORT=8082
$GOPATH/bin/health_check_agent $GOPATH/src/github.com/imdark/ipchecker/src/sample_ips.csv

```

Run another agent
```
go get github.com/imdark/ipchecker/src/health_check_agent

export TCP_PORT=8083
$GOPATH/bin/health_check_agent $GOPATH/src/github.com/imdark/ipchecker/src/mixed_ips.csv

```

Run master agent

```
go get github.com/imdark/ipchecker/src/health_check_master

export agent_ips_num=2
export agent_ip1=localhost:8082
export agent_ip2=localhost:8083
$GOPATH/bin/health_check_master

```

code assumes input file is a list of ips seperated by break lines
if no parameters are specified the code will try to run on ./sample_ips.csv

To Test:
--------
```go test```

Next steps:
-----------

- app currently assumes that list of ips will fit into memory
- more unit tests
- improve perfomance by only doing SYN and reciving ACK
