# go-iperf
A Go based wrapper around iperf3

## Basic Usage

basic client setup
```go
func main() {
	
	c := iperf.NewClient("192.168.0.10")
	c.SetJSON(true)
	c.SetIncludeServer(true)
	c.SetStreams(4)
	c.SetTimeSec(30)
	c.SetInterval(1)
	
	err := c.Start()
	if err != nil {
        fmt.Printf("failed to start client: %v\n", err)
        os.Exit(-1)
	}
	
	<- c.Done
	
	fmt.Println(c.Report().String())
}
```

basic server setup
```go
func main() {
	
	s := iperf.NewServer()
	err := s.Start()
	if err != nil {
        fmt.Printf("failed to start server: %v\n", err)
        os.Exit(-1)
    }
    
    for s.Running() {
    	time.Sleep(100 * time.Millisecond)
    }
    
    fmt.Println("server finished")
}
```

client with live results printing
```go
func main() {
	
	c := iperf.NewClient("192.168.0.10")
	c.SetJSON(true)
	c.SetIncludeServer(true)
	c.SetStreams(4)
	c.SetTimeSec(30)
	c.SetInterval(1)
	liveReports := c.SetModeLive()
	
	go func() {
	    for report := range liveReports {
	        fmt.Println(report.String())	
            }   	
        }   
	
	err := c.Start()
	if err != nil {
            fmt.Printf("failed to start client: %v\n", err)
            os.Exit(-1)
	}
	
	<- c.Done
	
	fmt.Println(c.Report().String())
}
```

building binary data package with iperf binaries
```
go-bindata -pkg iperf -prefix "embedded/" embedded/
```




## Sample Output

Linux iperf3 3.11 - client side

```shell
Connecting to host 127.0.0.1, port 5201
[  5] local 127.0.0.1 port 37338 connected to 127.0.0.1 port 5201
[ ID] Interval           Transfer     Bitrate         Retr  Cwnd
[  5]   0.00-1.00   sec  3.77 GBytes  32.3 Gbits/sec    0   1.50 MBytes       
[  5]   1.00-2.00   sec  4.06 GBytes  34.9 Gbits/sec    0   1.56 MBytes       
[  5]   2.00-3.00   sec  4.06 GBytes  34.9 Gbits/sec    0   1.56 MBytes       
[  5]   3.00-4.00   sec  3.99 GBytes  34.3 Gbits/sec    0   1.75 MBytes       
[  5]   4.00-5.00   sec  4.07 GBytes  35.0 Gbits/sec    0   2.25 MBytes       
[  5]   5.00-6.00   sec  4.03 GBytes  34.6 Gbits/sec    0   2.25 MBytes       
[  5]   6.00-7.00   sec  3.90 GBytes  33.5 Gbits/sec    0   5.18 MBytes       
[  5]   7.00-8.00   sec  3.89 GBytes  33.4 Gbits/sec    0   5.18 MBytes       
[  5]   8.00-9.00   sec  3.98 GBytes  34.2 Gbits/sec    0   5.18 MBytes       
[  5]   9.00-10.00  sec  4.07 GBytes  35.0 Gbits/sec    0   5.18 MBytes       
- - - - - - - - - - - - - - - - - - - - - - - - -
[ ID] Interval           Transfer     Bitrate         Retr
[  5]   0.00-10.00  sec  39.8 GBytes  34.2 Gbits/sec    0             sender
[  5]   0.00-10.00  sec  39.8 GBytes  34.2 Gbits/sec                  receiver

iperf Done.
```

Linux iperf3 3.11 - server side

```shell
-----------------------------------------------------------
Server listening on 5201 (test #1)
-----------------------------------------------------------
Accepted connection from 127.0.0.1, port 37336
[  5] local 127.0.0.1 port 5201 connected to 127.0.0.1 port 37338
[ ID] Interval           Transfer     Bitrate
[  5]   0.00-1.00   sec  3.77 GBytes  32.3 Gbits/sec                  
[  5]   1.00-2.00   sec  4.06 GBytes  34.9 Gbits/sec                  
[  5]   2.00-3.00   sec  4.06 GBytes  34.9 Gbits/sec                  
[  5]   3.00-4.00   sec  3.99 GBytes  34.3 Gbits/sec                  
[  5]   4.00-5.00   sec  4.07 GBytes  35.0 Gbits/sec                  
[  5]   5.00-6.00   sec  4.03 GBytes  34.6 Gbits/sec                  
[  5]   6.00-7.00   sec  3.90 GBytes  33.5 Gbits/sec                  
[  5]   7.00-8.00   sec  3.89 GBytes  33.4 Gbits/sec                  
[  5]   8.00-9.00   sec  3.98 GBytes  34.2 Gbits/sec                  
[  5]   9.00-10.00  sec  4.07 GBytes  35.0 Gbits/sec                  
[  5]  10.00-10.00  sec   384 KBytes  33.1 Gbits/sec                  
- - - - - - - - - - - - - - - - - - - - - - - - -
[ ID] Interval           Transfer     Bitrate
[  5]   0.00-10.00  sec  39.8 GBytes  34.2 Gbits/sec                  receiver
-----------------------------------------------------------
Server listening on 5201 (test #2)
-----------------------------------------------------------
```