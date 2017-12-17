

Design Goals
------------
50,000 IPs every minute

To Run:
-------
```go run main.go input/file/path```
code assumes input file is a list of ips seperated by break lines
if no parameters are specified the code will try to run on sample_ips.csv

To Test:
--------
```go test```

To profile performance 
----------------------
```go test -run=none -bench=ClientServerParallel4 -cpuprofile=cprof net/http```
Next steps:
-----------
[] ip parsing
[] getting a list of 5000 ips

[] app currently assumes that list of ips will fit into memory
[] added the linux environment variable you need to change
[] currently the code assumes all ips read from input files are ipv4

