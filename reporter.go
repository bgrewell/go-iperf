package iperf

import (
	"fmt"
	"github.com/BGrewell/tail"
)

type Reporter struct {
	ReportingChannel chan *StreamIntervalReport
	LogFile          string
	running          bool
	tailer *tail.Tail
}

func (r *Reporter) Start() {
	r.running = true
	go r.runLogProcessor()
}

func (r *Reporter) Stop() {
	r.running = false
	r.tailer.Stop()
	r.tailer.Cleanup()
	close(r.ReportingChannel)
	fmt.Println("reporter stopped")
}

// runLogProcessor is OS specific because of differences in iperf on Windows and Linux
