package iperf

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"runtime"
)

var (
	Debug          = false
	binaryDir      = ""
	binaryLocation = ""
)

func init() {
	// Extract the binaries
	if runtime.GOOS == "windows" {
		err := extractWindowsEmbeddedBinaries()
		if err != nil {
			log.Fatalf("error initializing iperf: %v", err)
		}
	} else if runtime.GOOS == "darwin" {
		err := extractMacEmbeddedBinaries()
		if err != nil {
			log.Fatalf("error initializing iperf: %v\n", err)
		}
	} else {
		err := extractLinuxEmbeddedBinaries()
		if err != nil {
			log.Fatalf("error initializing iperf: %v", err)
		}
	}
}

func Cleanup() {
	os.RemoveAll(binaryDir)
}

func ExtractBinaries() (err error) {
	files := []string{"cygwin1.dll", "iperf3.exe", "iperf3", "iperf3.app"}
	err = extractEmbeddedBinaries(files)
	fmt.Printf("files extracted to %s\n", binaryDir)
	return err
}

func extractWindowsEmbeddedBinaries() (err error) {
	files := []string{"cygwin1.dll", "iperf3.exe"}
	err = extractEmbeddedBinaries(files)
	binaryLocation = path.Join(binaryDir, "iperf3.exe")
	return err
}

func extractLinuxEmbeddedBinaries() (err error) {
	files := []string{"iperf3"}
	err = extractEmbeddedBinaries(files)
	binaryLocation = path.Join(binaryDir, "iperf3")
	return err
}

func extractMacEmbeddedBinaries() (err error) {
	files := []string{"iperf3.app"}
	err = extractEmbeddedBinaries(files)
	binaryLocation = path.Join(binaryDir, "iperf3.app")
	return err
}

func extractEmbeddedBinaries(files []string) (err error) {
	binaryDir, err = ioutil.TempDir("", "goiperf")
	if err != nil {
		return fmt.Errorf("failed to create temporary iperf directory: %v", err)
	}
	for _, file := range files {
		data, err := Asset(file)
		if err != nil {
			return fmt.Errorf("failed to extract embedded iperf: %v", err)
		}
		err = ioutil.WriteFile(path.Join(binaryDir, file), data, 0755)
		if err != nil {
			return fmt.Errorf("failed to save embedded iperf: %v", err)
		}
		if Debug {
			log.Printf("extracted file: %s\n", path.Join(binaryDir, file))
		}
	}
	return nil
}
