package iperf

import (
	"fmt"
	"github.com/google/uuid"
	"io"
)

func NewServer() *Server {
	s := &Server{
	}
	s.Id = uuid.New().String()
	return s
}

type Server struct {
	SharedOptions
	OneOff *bool
	ExitCode *int
	Running bool
	outputStream io.ReadCloser
	errorStream io.ReadCloser
}

func (s *Server) Start() (err error) {
	exit := make(chan int, 0)
	err = Execute(fmt.Sprintf("%s -s", binaryLocation), s.outputStream, s.errorStream, exit)
	if err != nil {
		return err
	}
	s.Running = true
	go func() {
		ds := DebugScanner{}
		ds.Scan(s.outputStream)
	}()
	go func() {
		ds := DebugScanner{}
		ds.Scan(s.errorStream)
	}()
	go func() {
		exitCode := <- exit
		s.ExitCode = &exitCode
		s.Running = false
	}()
	return nil
}

func (s *Server) Stop() {

}
