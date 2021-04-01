package iperf

import (
	"github.com/BGrewell/tail"
	"time"
)

type Reporter struct {
	ReportingChannel chan *StreamIntervalReport
	LogFile          string
	running          bool
	tailer           *tail.Tail
}

func (r *Reporter) Start() {
	r.running = true
	go r.runLogProcessor()
}

func (r *Reporter) Stop() {
	r.running = false
	r.tailer.Stop()
	r.tailer.Cleanup()
	for {
		if len(r.ReportingChannel) == 0 {
			break
		}
		time.Sleep(1 * time.Millisecond)
	}
	close(r.ReportingChannel)
}

// runLogProcessor is OS specific because of differences in iperf on Windows and Linux
