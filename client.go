package iperf

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func NewClient(host string) *Client {
	json := true
	proto := Protocol(PROTO_TCP)
	time := 10
	length := "128KB"
	streams := 1
	c := &Client{
		options: &ClientOptions{
			JSON:    &json,
			Proto:   &proto,
			TimeSec: &time,
			Length:  &length,
			Streams: &streams,
			Host:    &host,
		},
	}
	c.Id = uuid.New().String()
	c.Done = make(chan bool)
	return c
}

type ClientOptions struct {
	Port          *int      `json:"port" yaml:"port" xml:"port"`
	Format        *rune     `json:"format" yaml:"format" xml:"format"`
	Interval      *int      `json:"interval" yaml:"interval" xml:"interval"`
	JSON          *bool     `json:"json" yaml:"json" xml:"json"`
	LogFile       *string   `json:"log_file" yaml:"log_file" xml:"log_file"`
	Host          *string   `json:"host" yaml:"host" xml:"host"`
	Proto         *Protocol `json:"proto" yaml:"proto" xml:"proto"`
	Bandwidth     *string   `json:"bandwidth" yaml:"bandwidth" xml:"bandwidth"`
	TimeSec       *int      `json:"time_sec" yaml:"time_sec" xml:"time_sec"`
	Bytes         *string   `json:"bytes" yaml:"bytes" xml:"bytes"`
	BlockCount    *string   `json:"block_count" yaml:"block_count" xml:"block_count"`
	Length        *string   `json:"length" yaml:"length" xml:"length"`
	Streams       *int      `json:"streams" yaml:"streams" xml:"streams"`
	Reverse       *bool     `json:"reverse" yaml:"reverse" xml:"reverse"`
	Window        *string   `json:"window" yaml:"window" xml:"window"`
	MSS           *int      `json:"mss" yaml:"mss" xml:"mss"`
	NoDelay       *bool     `json:"no_delay" yaml:"no_delay" xml:"no_delay"`
	Version4      *bool     `json:"version_4" yaml:"version_4" xml:"version_4"`
	Version6      *bool     `json:"version_6" yaml:"version_6" xml:"version_6"`
	TOS           *int      `json:"tos" yaml:"tos" xml:"tos"`
	ZeroCopy      *bool     `json:"zero_copy" yaml:"zero_copy" xml:"zero_copy"`
	OmitSec       *int      `json:"omit_sec" yaml:"omit_sec" xml:"omit_sec"`
	Prefix        *string   `json:"prefix" yaml:"prefix" xml:"prefix"`
	IncludeServer *bool     `json:"include_server" yaml:"include_server" xml:"include_server"`
}

type Client struct {
	Id            string                     `json:"id" yaml:"id" xml:"id"`
	Running       bool                       `json:"running" yaml:"running" xml:"running"`
	Done          chan bool                  `json:"done" yaml:"done" xml:"done"`
	options       *ClientOptions             `json:"options" yaml:"options" xml:"options"`
	exitCode      *int                       `json:"exit_code" yaml:"exit_code" xml:"exit_code"`
	report        *TestReport                `json:"report" yaml:"report" xml:"report"`
	outputStream  io.ReadCloser              `json:"output_stream" yaml:"output_stream" xml:"output_stream"`
	errorStream   io.ReadCloser              `json:"error_stream" yaml:"error_stream" xml:"error_stream"`
	cancel        context.CancelFunc         `json:"cancel" yaml:"cancel" xml:"cancel"`
	mode          TestMode                   `json:"mode" yaml:"mode" xml:"mode"`
	live          bool                       `json:"live" yaml:"live" xml:"live"`
	reportingChan chan *StreamIntervalReport `json:"reporting_chan" yaml:"reporting_chan" xml:"reporting_chan"`
	reportingFile string                     `json:"reporting_file" yaml:"reporting_file" xml:"reporting_file"`
}

func (c *Client) LoadOptionsJSON(jsonStr string) (err error) {
	return json.Unmarshal([]byte(jsonStr), c.options)
}

func (c *Client) LoadOptions(options *ClientOptions) {
	c.options = options
}

func (c *Client) commandString() (cmd string, err error) {
	builder := strings.Builder{}
	if c.options.Host == nil || *c.options.Host == "" {
		return "", errors.New("unable to execute client. The field 'host' is required")
	}
	fmt.Fprintf(&builder, "%s -c %s", binaryLocation, c.Host())

	if c.options.Port != nil {
		fmt.Fprintf(&builder, " -p %d", c.Port())
	}

	if c.options.Format != nil && *c.options.Format != ' ' {
		fmt.Fprintf(&builder, " -f %c", c.Format())
	}

	if c.options.Interval != nil {
		fmt.Fprintf(&builder, " -i %d", c.Interval())
	}

	if c.options.Proto != nil && *c.options.Proto == PROTO_UDP {
		fmt.Fprintf(&builder, " -u")
	}

	if c.options.Bandwidth != nil {
		fmt.Fprintf(&builder, " -b %s", c.Bandwidth())
	}

	if c.options.TimeSec != nil {
		fmt.Fprintf(&builder, " -t %d", c.TimeSec())
	}

	if c.options.Bytes != nil {
		fmt.Fprintf(&builder, " -n %s", c.Bytes())
	}

	if c.options.BlockCount != nil {
		fmt.Fprintf(&builder, " -k %s", c.BlockCount())
	}

	if c.options.Length != nil {
		fmt.Fprintf(&builder, " -l %s", c.Length())
	}

	if c.options.Streams != nil {
		fmt.Fprintf(&builder, " -P %d", c.Streams())
	}

	if c.options.Reverse != nil && *c.options.Reverse {
		builder.WriteString(" -R")
	}

	if c.options.Window != nil {
		fmt.Fprintf(&builder, " -w %s", c.Window())
	}

	if c.options.MSS != nil {
		fmt.Fprintf(&builder, " -M %d", c.MSS())
	}

	if c.options.NoDelay != nil && *c.options.NoDelay {
		builder.WriteString(" -N")
	}

	if c.options.Version6 != nil && *c.options.Version6 {
		builder.WriteString(" -6")
	}

	if c.options.TOS != nil {
		fmt.Fprintf(&builder, " -S %d", c.TOS())
	}

	if c.options.ZeroCopy != nil && *c.options.ZeroCopy {
		builder.WriteString(" -Z")
	}

	if c.options.OmitSec != nil {
		fmt.Fprintf(&builder, " -O %d", c.OmitSec())
	}

	if c.options.Prefix != nil {
		fmt.Fprintf(&builder, " -T %s", c.Prefix())
	}

	if c.options.LogFile != nil && *c.options.LogFile != "" {
		fmt.Fprintf(&builder, " --logfile %s", c.LogFile())
	}

	if c.options.JSON != nil && *c.options.JSON {
		builder.WriteString(" -J")
	}

	if c.options.IncludeServer != nil && *c.options.IncludeServer {
		builder.WriteString(" --get-server-output")
	}

	return builder.String(), nil
}

func (c *Client) Host() string {
	if c.options.Host == nil {
		return ""
	}
	return *c.options.Host
}

func (c *Client) SetHost(host string) {
	c.options.Host = &host
}

func (c *Client) Port() int {
	if c.options.Port == nil {
		return 5201
	}
	return *c.options.Port
}

func (c *Client) SetPort(port int) {
	c.options.Port = &port
}

func (c *Client) Format() rune {
	if c.options.Format == nil {
		return ' '
	}
	return *c.options.Format
}

func (c *Client) SetFormat(format rune) {
	c.options.Format = &format
}

func (c *Client) Interval() int {
	if c.options.Interval == nil {
		return 1
	}
	return *c.options.Interval
}

func (c *Client) SetInterval(interval int) {
	c.options.Interval = &interval
}

func (c *Client) Proto() Protocol {
	if c.options.Proto == nil {
		return PROTO_TCP
	}
	return *c.options.Proto
}

func (c *Client) SetProto(proto Protocol) {
	c.options.Proto = &proto
}

func (c *Client) Bandwidth() string {
	if c.options.Bandwidth == nil && c.Proto() == PROTO_TCP {
		return "0"
	} else if c.options.Bandwidth == nil && c.Proto() == PROTO_UDP {
		return "1M"
	}
	return *c.options.Bandwidth
}

func (c *Client) SetBandwidth(bandwidth string) {
	c.options.Bandwidth = &bandwidth
}

func (c *Client) TimeSec() int {
	if c.options.TimeSec == nil {
		return 10
	}
	return *c.options.TimeSec
}

func (c *Client) SetTimeSec(timeSec int) {
	c.options.TimeSec = &timeSec
}

func (c *Client) Bytes() string {
	if c.options.Bytes == nil {
		return ""
	}
	return *c.options.Bytes
}

func (c *Client) SetBytes(bytes string) {
	c.options.Bytes = &bytes
}

func (c *Client) BlockCount() string {
	if c.options.BlockCount == nil {
		return ""
	}
	return *c.options.BlockCount
}

func (c *Client) SetBlockCount(blockCount string) {
	c.options.BlockCount = &blockCount
}

func (c *Client) Length() string {
	if c.options.Length == nil {
		if c.Proto() == PROTO_UDP {
			return "1460"
		} else {
			return "128K"
		}
	}
	return *c.options.Length
}

func (c *Client) SetLength(length string) {
	c.options.Length = &length
}

func (c *Client) Streams() int {
	if c.options.Streams == nil {
		return 1
	}
	return *c.options.Streams
}

func (c *Client) SetStreams(streamCount int) {
	c.options.Streams = &streamCount
}

func (c *Client) Reverse() bool {
	if c.options.Reverse == nil {
		return false
	}
	return *c.options.Reverse
}

func (c *Client) SetReverse(reverse bool) {
	c.options.Reverse = &reverse
}

func (c *Client) Window() string {
	if c.options.Window == nil {
		return ""
	}
	return *c.options.Window
}

func (c *Client) SetWindow(window string) {
	c.options.Window = &window
}

func (c *Client) MSS() int {
	if c.options.MSS == nil {
		return 1460
	}
	return *c.options.MSS
}

func (c *Client) SetMSS(mss int) {
	c.options.MSS = &mss
}

func (c *Client) NoDelay() bool {
	if c.options.NoDelay == nil {
		return false
	}
	return *c.options.NoDelay
}

func (c *Client) SetNoDelay(noDelay bool) {
	c.options.NoDelay = &noDelay
}

func (c *Client) Version4() bool {
	if c.options.Version6 == nil && c.options.Version4 == nil {
		return true
	} else if c.options.Version6 != nil && *c.options.Version6 == true {
		return false
	}
	return *c.options.Version4
}

func (c *Client) SetVersion4(set bool) {
	c.options.Version4 = &set
}

func (c *Client) Version6() bool {
	if c.options.Version6 == nil {
		return false
	}
	return *c.options.Version6
}

func (c *Client) SetVersion6(set bool) {
	c.options.Version6 = &set
}

func (c *Client) TOS() int {
	if c.options.TOS == nil {
		return 0
	}
	return *c.options.TOS
}

func (c *Client) SetTOS(value int) {
	c.options.TOS = &value
}

func (c *Client) ZeroCopy() bool {
	if c.options.ZeroCopy == nil {
		return false
	}
	return *c.options.ZeroCopy
}

func (c *Client) SetZeroCopy(set bool) {
	c.options.ZeroCopy = &set
}

func (c *Client) OmitSec() int {
	if c.options.OmitSec == nil {
		return 0
	}
	return *c.options.OmitSec
}

func (c *Client) SetOmitSec(value int) {
	c.options.OmitSec = &value
}

func (c *Client) Prefix() string {
	if c.options.Prefix == nil {
		return ""
	}
	return *c.options.Prefix
}

func (c *Client) SetPrefix(prefix string) {
	c.options.Prefix = &prefix
}

func (c *Client) LogFile() string {
	if c.options.LogFile == nil {
		return ""
	}
	return *c.options.LogFile
}

func (c *Client) SetLogFile(logfile string) {
	c.options.LogFile = &logfile
}

func (c *Client) JSON() bool {
	if c.options.JSON == nil {
		return false
	}
	return *c.options.JSON
}

func (c *Client) SetJSON(set bool) {
	c.options.JSON = &set
}

func (c *Client) IncludeServer() bool {
	if c.options.IncludeServer == nil {
		return false
	}
	return *c.options.IncludeServer
}

func (c *Client) SetIncludeServer(set bool) {
	c.options.IncludeServer = &set
}

func (c *Client) ExitCode() *int {
	return c.exitCode
}

func (c *Client) Report() *TestReport {
	return c.report
}

func (c *Client) Mode() TestMode {
	return c.mode
}

func (c *Client) SetModeJson() {
	c.SetJSON(true)
	c.reportingChan = nil
	c.reportingFile = ""
}

func (c *Client) SetModeLive() <-chan *StreamIntervalReport {
	c.SetJSON(false) // having JSON == true will cause reporting to fail
	c.live = true
	c.reportingChan = make(chan *StreamIntervalReport, 10000)
	f, err := ioutil.TempFile("", "iperf_")
	if err != nil {
		log.Fatalf("failed to create logfile: %v", err)
	}
	c.reportingFile = f.Name()
	c.SetLogFile(c.reportingFile)
	return c.reportingChan
}

func (c *Client) Start() (err error) {
	//todo: Need to build the string based on the options above that are set
	cmd, err := c.commandString()
	if err != nil {
		return err
	}
	fmt.Println(cmd)
	var exit chan int
	c.outputStream, c.errorStream, exit, c.cancel, err = ExecuteAsyncWithCancel(cmd)
	if err != nil {
		return err
	}
	c.Running = true
	//go func() {
	//	ds := DebugScanner{Silent: false}
	//	ds.Scan(c.outputStream)
	//}()
	//go func() {
	//	ds := DebugScanner{Silent: false}
	//	ds.Scan(c.errorStream)
	//}()
	go func() {
		var reporter *Reporter
		if c.live {
			reporter = &Reporter{
				ReportingChannel: c.reportingChan,
				LogFile:          c.reportingFile,
			}
			reporter.Start()
		} else {
			testOutput, err := ioutil.ReadAll(c.outputStream)
			if err != nil {
				return
			}
			c.report, err = Loads(string(testOutput))
		}
		exitCode := <-exit
		c.exitCode = &exitCode
		c.Running = false
		c.Done <- true
		if reporter != nil {
			reporter.Stop()
		}
	}()
	return nil
}

func (c *Client) Stop() {
	if c.Running && c.cancel != nil {
		c.cancel()
		os.Remove(c.reportingFile)
		c.Done <- true
	}
}
