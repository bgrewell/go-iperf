package iperf

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"io"
	"time"
)

func NewServer() *Server {
	s := &Server{}
	s.Id = uuid.New().String()
	return s
}

type Server struct {
	SharedOptions
	Id           string
	OneOff       *bool
	ExitCode     *int
	Running      bool
	outputStream io.ReadCloser
	errorStream  io.ReadCloser
	cancel       context.CancelFunc
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
