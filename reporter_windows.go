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

func (r *Reporter) runLogProcessor() {
	var err error
	r.tailer, err = tail.TailFile(r.LogFile, tail.Config{
		Follow:    true,
		ReOpen:    true,
		Poll:      true,
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
				if len(fields) >= 6 {
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
					omitted := false
					if len(fields) >= 7 && fields[6] == "(omitted)" {
						omitted = true
					}
					report := &StreamIntervalReport{
						Socket:        stream,
						StartInterval: float32(start),
						EndInterval:   float32(end),
						Seconds:       float32(end - start),
						Bytes:         int(transferedBytes),
						BitsPerSecond: float64(rate),
						Omitted:       omitted,
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
