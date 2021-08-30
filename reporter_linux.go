package iperf

import (
	"fmt"
	"github.com/BGrewell/go-conversions"
	"github.com/BGrewell/tail"
	"log"
	"strconv"
	"strings"
	"time"
)

/*
Connecting to host 10.254.100.100, port 5201
[  4] local 192.168.3.182 port 54104 connected to 10.254.100.100 port 5201
[ ID] Interval           Transfer     Bandwidth       Retr  Cwnd
[  4]   0.00-1.00   sec   109 MBytes   913 Mbits/sec   13    634 KBytes       (omitted)
[  4]   1.00-2.00   sec   110 MBytes   927 Mbits/sec    7    550 KBytes       (omitted)
[  4]   2.00-3.00   sec   109 MBytes   918 Mbits/sec    6    559 KBytes       (omitted)
[  4]   3.00-4.00   sec   111 MBytes   930 Mbits/sec    6    690 KBytes       (omitted)
[  4]   4.00-5.00   sec   111 MBytes   933 Mbits/sec    0    803 KBytes       (omitted)
[  4]   5.00-6.00   sec   111 MBytes   933 Mbits/sec    6    673 KBytes       (omitted)
[  4]   6.00-7.00   sec   111 MBytes   932 Mbits/sec    6    605 KBytes       (omitted)
[  4]   7.00-8.00   sec   110 MBytes   925 Mbits/sec    0    732 KBytes       (omitted)
[  4]   8.00-9.00   sec   111 MBytes   932 Mbits/sec    0    840 KBytes       (omitted)
[  4]   9.00-10.00  sec   110 MBytes   923 Mbits/sec    6    690 KBytes       (omitted)
[  4]   0.00-1.00   sec   111 MBytes   928 Mbits/sec    6    618 KBytes
[  4]   1.00-2.00   sec   111 MBytes   931 Mbits/sec    0    745 KBytes
[  4]   2.00-3.00   sec   111 MBytes   929 Mbits/sec   11    614 KBytes
[  4]   3.00-4.00   sec   110 MBytes   922 Mbits/sec    6    551 KBytes
[  4]   4.00-5.00   sec   111 MBytes   933 Mbits/sec    6    519 KBytes
[  4]   5.00-6.00   sec   111 MBytes   928 Mbits/sec    0    663 KBytes
[  4]   6.00-7.00   sec   111 MBytes   932 Mbits/sec    0    783 KBytes
[  4]   7.00-8.00   sec   111 MBytes   933 Mbits/sec    6    656 KBytes
[  4]   8.00-9.00   sec   111 MBytes   933 Mbits/sec    6    598 KBytes
[  4]   9.00-10.00  sec   110 MBytes   925 Mbits/sec    0    728 KBytes
[  4]  10.00-11.00  sec   111 MBytes   933 Mbits/sec    0    839 KBytes
[  4]  11.00-12.00  sec   109 MBytes   918 Mbits/sec    6    680 KBytes
[  4]  12.00-12.24  sec  25.0 MBytes   888 Mbits/sec    0    711 KBytes
- - - - - - - - - - - - - - - - - - - - - - - - -
[ ID] Interval           Transfer     Bandwidth       Retr
[  4]   0.00-12.24  sec  1.32 GBytes   928 Mbits/sec   47             sender
[  4]   0.00-12.24  sec  0.00 Bytes  0.00 bits/sec                  receiver

iperf Done.

*/
func (r *Reporter) runLogProcessor() {
	var err error
	r.tailer, err = tail.TailFile(r.LogFile, tail.Config{
		Follow:    true,
		ReOpen:    true,
		Poll:      false, // on linux we don't need to poll as the fsnotify works properly
		MustExist: true,
	})
	if err != nil {
		log.Fatalf("failed to tail log file: %v", err)
	}

	for {
		select {
		case line := <-r.tailer.Lines:
			if line == nil {
				continue
			}
			if len(line.Text) > 5 {
				id := line.Text[1:4]
				stream, err := strconv.Atoi(strings.TrimSpace(id))
				if err != nil {
					continue
				}
				fields := strings.Fields(line.Text[5:])
				if len(fields) >= 9 {
					if fields[0] == "local" {
						continue
					}
					timeFields := strings.Split(fields[0], "-")
					start, err := strconv.ParseFloat(timeFields[0], 32)
					if err != nil {
						log.Printf("failed to convert start time: %s\n", err)
					}
					end, err := strconv.ParseFloat(timeFields[1], 32)
					transferedStr := fmt.Sprintf("%s%s", fields[2], fields[3])
					transferedBytes, err := conversions.StringBitRateToInt(transferedStr)
					if err != nil {
						log.Printf("failed to convert units: %s\n", err)
					}
					transferedBytes = transferedBytes / 8
					rateStr := fmt.Sprintf("%s%s", fields[4], fields[5])
					rate, err := conversions.StringBitRateToInt(rateStr)
					if err != nil {
						log.Printf("failed to convert units: %s\n", err)
					}
					retrans, err := strconv.Atoi(fields[6])
					if err != nil {
						log.Printf("failed to convert units: %s\n", err)
					}
					cwndStr := fmt.Sprintf("%s%s", fields[7], fields[8])
					cwnd, err := conversions.StringBitRateToInt(cwndStr)
					if err != nil {
						log.Printf("failed to convert units: %s\n", err)
					}
					cwnd = cwnd / 8
					omitted := false
					if len(fields) >= 10 && fields[9] == "(omitted)" {
						omitted = true
					}
					report := &StreamIntervalReport{
						Socket:           stream,
						StartInterval:    float32(start),
						EndInterval:      float32(end),
						Seconds:          float32(end - start),
						Bytes:            int(transferedBytes),
						BitsPerSecond:    float64(rate),
						Retransmissions:  retrans,
						CongestionWindow: int(cwnd),
						Omitted:          omitted,
					}
					r.ReportingChannel <- report
				}
			}
		case <-time.After(100 * time.Millisecond):
			if !r.running {
				return
			}
		}
	}
}
