package iperf

import (
	"context"
	"io"
	"os/exec"
	"strings"
	"syscall"
	"time"
)

func ExecuteAsync(cmd string) (outPipe io.ReadCloser, errPipe io.ReadCloser, exitCode chan int, err error) {
	exitCode = make(chan int)
	cmdParts := strings.Fields(cmd)
	binary, err := exec.LookPath(cmdParts[0])
	if err != nil {
		return nil, nil, nil, err
	}
	exe := exec.Command(binary, cmdParts[1:]...)
	outPipe, err = exe.StdoutPipe()
	if err != nil {
		return nil, nil, nil, err
	}
	errPipe, err = exe.StderrPipe()
	if err != nil {
		return nil, nil, nil, err
	}
	err = exe.Start()
	if err != nil {
		return nil, nil, nil, err
	}
	go func() {
		if err := exe.Wait(); err != nil {
			if exiterr, ok := err.(*exec.ExitError); ok {
				if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
					exitCode <- status.ExitStatus()
				}
			}
		} else {
			exitCode <- 0
		}
	}()
	return outPipe, errPipe, exitCode, nil
}

func ExecuteAsyncWithCancel(cmd string) (stdOut io.ReadCloser, stdErr io.ReadCloser, exitCode chan int, cancelToken context.CancelFunc, err error) {
	exitCode = make(chan int)
	ctx, cancel := context.WithCancel(context.Background())
	cmdParts := strings.Fields(cmd)
	binary, err := exec.LookPath(cmdParts[0])
	if err != nil {
		defer cancel()
		return nil, nil, nil, nil, err
	}
	exe := exec.CommandContext(ctx, binary, cmdParts[1:]...)
	stdOut, err = exe.StdoutPipe()
	if err != nil {
		defer cancel()
		return nil, nil, nil, nil, err
	}
	stdErr, err = exe.StderrPipe()
	if err != nil {
		defer cancel()
		return nil, nil, nil, nil, err
	}
	err = exe.Start()
	if err != nil {
		defer cancel()
		return nil, nil, nil, nil, err
	}
	go func() {
		time.Sleep(30 * time.Second)
		if err := exe.Wait(); err != nil {
			if exiterr, ok := err.(*exec.ExitError); ok {
				if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
					exitCode <- status.ExitStatus()
				}
			}
		} else {
			exitCode <- 0
		}
	}()
	return stdOut, stdErr, exitCode, cancel, nil
}
