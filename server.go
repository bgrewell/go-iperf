package iperf

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io"
	"strings"
	"time"
)

var (
	defaultPort     = 5201
	defaultInterval = 1
	defaultJSON     = true
)

func NewServer() *Server {
	s := &Server{
		Options: &ServerOptions{
			Port:     &defaultPort,
			Interval: &defaultInterval,
			JSON:     &defaultJSON,
		},
	}
	s.Id = uuid.New().String()
	return s
}

type ServerOptions struct {
	OneOff   *bool   `json:"one_off" yaml:"one_off" xml:"one_off"`
	Port     *int    `json:"port" yaml:"port" xml:"port"`
	Format   *rune   `json:"format" yaml:"format" xml:"format"`
	Interval *int    `json:"interval" yaml:"interval" xml:"interval"`
	JSON     *bool   `json:"json" yaml:"json" xml:"json"`
	LogFile  *string `json:"log_file" yaml:"log_file" xml:"log_file"`
}

type Server struct {
	Id           string             `json:"id" yaml:"id" xml:"id"`
	Running      bool               `json:"running" yaml:"running" xml:"running"`
	Options      *ServerOptions     `json:"-" yaml:"-" xml:"-"`
	ExitCode     *int               `json:"exit_code" yaml:"exit_code" xml:"exit_code"`
	Debug        bool               `json:"-" yaml:"-" xml:"-"`
	StdOut       bool               `json:"-" yaml:"-" xml:"-"`
	outputStream io.ReadCloser      `json:"output_stream" yaml:"output_stream" xml:"output_stream"`
	errorStream  io.ReadCloser      `json:"error_stream" yaml:"error_stream" xml:"error_stream"`
	cancel       context.CancelFunc `json:"cancel" yaml:"cancel" xml:"cancel"`
}

func (s *Server) LoadOptionsJSON(jsonStr string) (err error) {
	return json.Unmarshal([]byte(jsonStr), s.Options)
}

func (s *Server) LoadOptions(options *ServerOptions) {
	s.Options = options
}

func (s *Server) commandString() (cmd string, err error) {
	builder := strings.Builder{}
	fmt.Fprintf(&builder, "%s -s", binaryLocation)

	if s.Options.OneOff != nil && s.OneOff() == true {
		builder.WriteString(" --one-off")
	}

	if s.Options.Port != nil {
		fmt.Fprintf(&builder, " --port %d", s.Port())
	}

	if s.Options.Format != nil {
		fmt.Fprintf(&builder, " --format %c", s.Format())
	}

	if s.Options.Interval != nil {
		fmt.Fprintf(&builder, " --interval %d", s.Interval())
	}

	if s.Options.JSON != nil && s.JSON() == true {
		builder.WriteString(" --json")
	}

	if s.Options.LogFile != nil && s.LogFile() != "" {
		fmt.Fprintf(&builder, " --logfile %s --forceflush", s.LogFile())
	}

	return builder.String(), nil
}

func (s *Server) OneOff() bool {
	if s.Options.OneOff == nil {
		return false
	}
	return *s.Options.OneOff
}

func (s *Server) SetOneOff(oneOff bool) {
	s.Options.OneOff = &oneOff
}

func (s *Server) Port() int {
	if s.Options.Port == nil {
		return defaultPort
	}
	return *s.Options.Port
}

func (s *Server) SetPort(port int) {
	s.Options.Port = &port
}

func (s *Server) Format() rune {
	if s.Options.Format == nil {
		return ' '
	}
	return *s.Options.Format
}

func (s *Server) SetFormat(format rune) {
	s.Options.Format = &format
}

func (s *Server) Interval() int {
	if s.Options.Interval == nil {
		return defaultInterval
	}
	return *s.Options.Interval
}

func (s *Server) JSON() bool {
	if s.Options.JSON == nil {
		return false
	}
	return *s.Options.JSON
}

func (s *Server) SetJSON(json bool) {
	s.Options.JSON = &json
}

func (s *Server) LogFile() string {
	if s.Options.LogFile == nil {
		return ""
	}
	return *s.Options.LogFile
}

func (s *Server) SetLogFile(filename string) {
	s.Options.LogFile = &filename
}

func (s *Server) Start() (err error) {
	_, err = s.start()
	return err
}

func (s *Server) StartEx() (pid int, err error) {
	return s.start()
}

func (s *Server) start() (pid int, err error) {
	cmd, err := s.commandString()
	if err != nil {
		return -1, err
	}
	var exit chan int

	if s.Debug {
		fmt.Printf("executing command: %s\n", cmd)
	}
	s.outputStream, s.errorStream, exit, s.cancel, pid, err = ExecuteAsyncWithCancel(cmd)

	if err != nil {
		return -1, err
	}
	s.Running = true

	go func() {
		ds := DebugScanner{Silent: !s.StdOut}
		ds.Scan(s.outputStream)
	}()
	go func() {
		ds := DebugScanner{Silent: !s.Debug}
		ds.Scan(s.errorStream)
	}()

	go func() {
		exitCode := <-exit
		s.ExitCode = &exitCode
		s.Running = false
	}()
	return pid,nil
}

func (s *Server) Stop() {
	if s.Running && s.cancel != nil {
		s.cancel()
		time.Sleep(100 * time.Millisecond)
	}
}
