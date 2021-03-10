package iperf

type Reporter struct {
	ReportingChannel chan *StreamIntervalReport
	LogFile          string
	running          bool
}

func (r *Reporter) Start() {
	r.running = true
	go r.runLogProcessor()
}

func (r *Reporter) Stop() {
	r.running = false
	close(r.ReportingChannel)
}

// runLogProcessor is OS specific because of differences in iperf on Windows and Linux
