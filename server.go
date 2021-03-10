package iperf

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"io"
	"time"
)

func NewServer() *Server {
	defaultPort := 5201
	defaultInterval := 1
	s := &Server{
		Port:     &defaultPort,
		Interval: &defaultInterval,
	}
	s.Id = uuid.New().String()
	return s
}

type Server struct {
	Id            string             `json:"id" yaml:"id" xml:"id"`
	OneOff        *bool              `json:"one_off" yaml:"one_off" xml:"one_off"`
	ExitCode      *int               `json:"exit_code" yaml:"exit_code" xml:"exit_code"`
	Port          *int               `json:"port" yaml:"port" xml:"port"`
	Format        *rune              `json:"format" yaml:"format" xml:"format"`
	Interval      *int               `json:"interval" yaml:"interval" xml:"interval"`
	JSON          *bool              `json:"json" yaml:"json" xml:"json"`
	LogFile       *string            `json:"log_file" yaml:"log_file" xml:"log_file"`
	Running       bool               `json:"running" yaml:"running" xml:"running"`
	outputStream  io.ReadCloser      `json:"output_stream" yaml:"output_stream" xml:"output_stream"`
	errorStream   io.ReadCloser      `json:"error_stream" yaml:"error_stream" xml:"error_stream"`
	cancel        context.CancelFunc `json:"cancel" yaml:"cancel" xml:"cancel"`
}

func (s *Server) Start() (err error) {
	var exit chan int
	s.outputStream, s.errorStream, exit, s.cancel, err = ExecuteAsyncWithCancel(fmt.Sprintf("%s -s -J", binaryLocation))
	if err != nil {
		return err
	}
	s.Running = true
	go func() {
		ds := DebugScanner{Silent: true}
		ds.Scan(s.outputStream)
	}()
	go func() {
		ds := DebugScanner{Silent: true}
		ds.Scan(s.errorStream)
	}()
	go func() {
		exitCode := <-exit
		s.ExitCode = &exitCode
		s.Running = false
	}()
	return nil
}

func (s *Server) Stop() {
	if s.Running && s.cancel != nil {
		s.cancel()
		time.Sleep(100 * time.Millisecond)
	}
}
