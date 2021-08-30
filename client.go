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
		Options: &ClientOptions{
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
	Id            string         `json:"id" yaml:"id" xml:"id"`
	Running       bool           `json:"running" yaml:"running" xml:"running"`
	Done          chan bool      `json:"-" yaml:"-" xml:"-"`
	Options       *ClientOptions `json:"options" yaml:"options" xml:"options"`
	Debug         bool           `json:"-" yaml:"-" xml:"-"`
	StdOut        bool           `json:"-" yaml:"-" xml:"-"`
	exitCode      *int
	report        *TestReport
	outputStream  io.ReadCloser
	errorStream   io.ReadCloser
	cancel        context.CancelFunc
	mode          TestMode
	live          bool
	reportingChan chan *StreamIntervalReport
	reportingFile string
}

func (c *Client) LoadOptionsJSON(jsonStr string) (err error) {
	return json.Unmarshal([]byte(jsonStr), c.Options)
}

func (c *Client) LoadOptions(options *ClientOptions) {
	c.Options = options
}

func (c *Client) commandString() (cmd string, err error) {
	builder := strings.Builder{}
	if c.Options.Host == nil || *c.Options.Host == "" {
		return "", errors.New("unable to execute client. The field 'host' is required")
	}
	fmt.Fprintf(&builder, "%s -c %s", binaryLocation, c.Host())

	if c.Options.Port != nil {
		fmt.Fprintf(&builder, " -p %d", c.Port())
	}

	if c.Options.Format != nil && *c.Options.Format != ' ' {
		fmt.Fprintf(&builder, " -f %c", c.Format())
	}

	if c.Options.Interval != nil {
		fmt.Fprintf(&builder, " -i %d", c.Interval())
	}

	if c.Options.Proto != nil && *c.Options.Proto == PROTO_UDP {
		fmt.Fprintf(&builder, " -u")
	}

	if c.Options.Bandwidth != nil {
		fmt.Fprintf(&builder, " -b %s", c.Bandwidth())
	}

	if c.Options.TimeSec != nil {
		fmt.Fprintf(&builder, " -t %d", c.TimeSec())
	}

	if c.Options.Bytes != nil {
		fmt.Fprintf(&builder, " -n %s", c.Bytes())
	}

	if c.Options.BlockCount != nil {
		fmt.Fprintf(&builder, " -k %s", c.BlockCount())
	}

	if c.Options.Length != nil {
		fmt.Fprintf(&builder, " -l %s", c.Length())
	}

	if c.Options.Streams != nil {
		fmt.Fprintf(&builder, " -P %d", c.Streams())
	}

	if c.Options.Reverse != nil && *c.Options.Reverse {
		builder.WriteString(" -R")
	}

	if c.Options.Window != nil {
		fmt.Fprintf(&builder, " -w %s", c.Window())
	}

	if c.Options.MSS != nil {
		fmt.Fprintf(&builder, " -M %d", c.MSS())
	}

	if c.Options.NoDelay != nil && *c.Options.NoDelay {
		builder.WriteString(" -N")
	}

	if c.Options.Version6 != nil && *c.Options.Version6 {
		builder.WriteString(" -6")
	}

	if c.Options.TOS != nil {
		fmt.Fprintf(&builder, " -S %d", c.TOS())
	}

	if c.Options.ZeroCopy != nil && *c.Options.ZeroCopy {
		builder.WriteString(" -Z")
	}

	if c.Options.OmitSec != nil {
		fmt.Fprintf(&builder, " -O %d", c.OmitSec())
	}

	if c.Options.Prefix != nil {
		fmt.Fprintf(&builder, " -T %s", c.Prefix())
	}

	if c.Options.LogFile != nil && *c.Options.LogFile != "" {
		fmt.Fprintf(&builder, " --logfile %s", c.LogFile())
	}

	if c.Options.JSON != nil && *c.Options.JSON {
		builder.WriteString(" -J")
	}

	if c.Options.IncludeServer != nil && *c.Options.IncludeServer {
		builder.WriteString(" --get-server-output")
	}

	return builder.String(), nil
}

func (c *Client) Host() string {
	if c.Options.Host == nil {
		return ""
	}
	return *c.Options.Host
}

func (c *Client) SetHost(host string) {
	c.Options.Host = &host
}

func (c *Client) Port() int {
	if c.Options.Port == nil {
		return 5201
	}
	return *c.Options.Port
}

func (c *Client) SetPort(port int) {
	c.Options.Port = &port
}

func (c *Client) Format() rune {
	if c.Options.Format == nil {
		return ' '
	}
	return *c.Options.Format
}

func (c *Client) SetFormat(format rune) {
	c.Options.Format = &format
}

func (c *Client) Interval() int {
	if c.Options.Interval == nil {
		return 1
	}
	return *c.Options.Interval
}

func (c *Client) SetInterval(interval int) {
	c.Options.Interval = &interval
}

func (c *Client) Proto() Protocol {
	if c.Options.Proto == nil {
		return PROTO_TCP
	}
	return *c.Options.Proto
}

func (c *Client) SetProto(proto Protocol) {
	c.Options.Proto = &proto
}

func (c *Client) Bandwidth() string {
	if c.Options.Bandwidth == nil && c.Proto() == PROTO_TCP {
		return "0"
	} else if c.Options.Bandwidth == nil && c.Proto() == PROTO_UDP {
		return "1M"
	}
	return *c.Options.Bandwidth
}

func (c *Client) SetBandwidth(bandwidth string) {
	c.Options.Bandwidth = &bandwidth
}

func (c *Client) TimeSec() int {
	if c.Options.TimeSec == nil {
		return 10
	}
	return *c.Options.TimeSec
}

func (c *Client) SetTimeSec(timeSec int) {
	c.Options.TimeSec = &timeSec
}

func (c *Client) Bytes() string {
	if c.Options.Bytes == nil {
		return ""
	}
	return *c.Options.Bytes
}

func (c *Client) SetBytes(bytes string) {
	c.Options.Bytes = &bytes
}

func (c *Client) BlockCount() string {
	if c.Options.BlockCount == nil {
		return ""
	}
	return *c.Options.BlockCount
}

func (c *Client) SetBlockCount(blockCount string) {
	c.Options.BlockCount = &blockCount
}

func (c *Client) Length() string {
	if c.Options.Length == nil {
		if c.Proto() == PROTO_UDP {
			return "1460"
		} else {
			return "128K"
		}
	}
	return *c.Options.Length
}

func (c *Client) SetLength(length string) {
	c.Options.Length = &length
}

func (c *Client) Streams() int {
	if c.Options.Streams == nil {
		return 1
	}
	return *c.Options.Streams
}

func (c *Client) SetStreams(streamCount int) {
	c.Options.Streams = &streamCount
}

func (c *Client) Reverse() bool {
	if c.Options.Reverse == nil {
		return false
	}
	return *c.Options.Reverse
}

func (c *Client) SetReverse(reverse bool) {
	c.Options.Reverse = &reverse
}

func (c *Client) Window() string {
	if c.Options.Window == nil {
		return ""
	}
	return *c.Options.Window
}

func (c *Client) SetWindow(window string) {
	c.Options.Window = &window
}

func (c *Client) MSS() int {
	if c.Options.MSS == nil {
		return 1460
	}
	return *c.Options.MSS
}

func (c *Client) SetMSS(mss int) {
	c.Options.MSS = &mss
}

func (c *Client) NoDelay() bool {
	if c.Options.NoDelay == nil {
		return false
	}
	return *c.Options.NoDelay
}

func (c *Client) SetNoDelay(noDelay bool) {
	c.Options.NoDelay = &noDelay
}

func (c *Client) Version4() bool {
	if c.Options.Version6 == nil && c.Options.Version4 == nil {
		return true
	} else if c.Options.Version6 != nil && *c.Options.Version6 == true {
		return false
	}
	return *c.Options.Version4
}

func (c *Client) SetVersion4(set bool) {
	c.Options.Version4 = &set
}

func (c *Client) Version6() bool {
	if c.Options.Version6 == nil {
		return false
	}
	return *c.Options.Version6
}

func (c *Client) SetVersion6(set bool) {
	c.Options.Version6 = &set
}

func (c *Client) TOS() int {
	if c.Options.TOS == nil {
		return 0
	}
	return *c.Options.TOS
}

func (c *Client) SetTOS(value int) {
	c.Options.TOS = &value
}

func (c *Client) ZeroCopy() bool {
	if c.Options.ZeroCopy == nil {
		return false
	}
	return *c.Options.ZeroCopy
}

func (c *Client) SetZeroCopy(set bool) {
	c.Options.ZeroCopy = &set
}

func (c *Client) OmitSec() int {
	if c.Options.OmitSec == nil {
		return 0
	}
	return *c.Options.OmitSec
}

func (c *Client) SetOmitSec(value int) {
	c.Options.OmitSec = &value
}

func (c *Client) Prefix() string {
	if c.Options.Prefix == nil {
		return ""
	}
	return *c.Options.Prefix
}

func (c *Client) SetPrefix(prefix string) {
	c.Options.Prefix = &prefix
}

func (c *Client) LogFile() string {
	if c.Options.LogFile == nil {
		return ""
	}
	return *c.Options.LogFile
}

func (c *Client) SetLogFile(logfile string) {
	c.Options.LogFile = &logfile
}

func (c *Client) JSON() bool {
	if c.Options.JSON == nil {
		return false
	}
	return *c.Options.JSON
}

func (c *Client) SetJSON(set bool) {
	c.Options.JSON = &set
}

func (c *Client) IncludeServer() bool {
	if c.Options.IncludeServer == nil {
		return false
	}
	return *c.Options.IncludeServer
}

func (c *Client) SetIncludeServer(set bool) {
	c.Options.IncludeServer = &set
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
	_, err = c.start()
	return err
}

func (c *Client) StartEx() (pid int, err error) {
	return c.start()
}

func (c *Client) start() (pid int, err error) {
	read := make(chan interface{})
	cmd, err := c.commandString()
	if err != nil {
		return -1, err
	}
	var exit chan int

	if c.Debug {
		fmt.Printf("executing command: %s\n", cmd)
	}
	c.outputStream, c.errorStream, exit, c.cancel, pid, err = ExecuteAsyncWithCancelReadIndicator(cmd, read)

	if err != nil {
		return -1, err
	}
	c.Running = true

	//go func() {
	//	ds := DebugScanner{Silent: !c.StdOut}
	//	ds.Scan(c.outputStream)
	//}()
	go func() {
		ds := DebugScanner{Silent: !c.Debug}
		ds.Scan(c.errorStream)
	}()

	go func() {
		var reporter *Reporter
		if c.live {
			reporter = &Reporter{
				ReportingChannel: c.reportingChan,
				LogFile:          c.reportingFile,
			}
			reporter.Start()
		} else {
			if c.Debug {
				fmt.Println("reading output")
			}
			testOutput, err := ioutil.ReadAll(c.outputStream)
			read <- true
			if err != nil {
				if c.Debug {
					fmt.Println(err.Error())
				}
			}
			if c.Debug {
				fmt.Println("parsing output")
			}
			c.report, err = Loads(string(testOutput))
			if err != nil && c.Debug {
				fmt.Println(err.Error())
			}
		}
		if c.Debug {
			fmt.Println("complete")
		}
		exitCode := <-exit
		c.exitCode = &exitCode
		c.Running = false
		c.Done <- true
		if reporter != nil {
			reporter.Stop()
		}
	}()
	return pid, nil
}

func (c *Client) Stop() {
	if c.Running && c.cancel != nil {
		c.cancel()
		os.Remove(c.reportingFile)
	}
}
