package iperf

import (
	"io"
	"os/exec"
	"strings"
	"syscall"
)

func Execute(cmd string, outPipe io.ReadCloser, errPipe io.ReadCloser, exit chan <- int) (err error) {
	cmdParts := strings.Fields(cmd)
	binary, err := exec.LookPath(cmdParts[0])
	if err != nil {
		return err
	}
	exe := exec.Command(binary, cmdParts[1:]...)
	outPipe, err = exe.StdoutPipe()
	if err != nil {
		return err
	}
	errPipe, err = exe.StderrPipe()
	if err != nil {
		return err
	}
	err = exe.Start()
	if err != nil {
		return err
	}
	go func() {
		if err := exe.Wait(); err != nil {
			if exiterr, ok := err.(*exec.ExitError); ok {
				if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
					exit <- status.ExitStatus()
				}
			}
		} else {
			exit <- 0
		}
	}()
	return nil
}
