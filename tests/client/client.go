package main

import (
	"fmt"
	"github.com/BGrewell/go-iperf"
	"os"
)

func main() {

	includeServer := true
	proto := "tcp"
	runTime := 30
	omitSec := 10
	length := "65500"

	c := iperf.NewClient("10.254.100.100")
	c.SetIncludeServer(includeServer)
	c.SetTimeSec(runTime)
	c.SetOmitSec(omitSec)
	c.SetProto((iperf.Protocol)(proto))
	c.SetLength(length)
	c.SetJSON(false)
	c.SetIncludeServer(false)
	c.SetTimeSec(20)
	c.SetStreams(2)
	reports := c.SetModeLive()

	go func() {
		for report := range reports {
			fmt.Println(report.String())
		}
	}()

	err := c.Start()
	if err != nil {
		fmt.Println("failed to start client")
		os.Exit(-1)
	}

	// Method 1: Wait for the test to finish by pulling from the 'Done' channel which will block until something is put in or it's closed
	<-c.Done

	// Method 2: Poll the c.Running state and wait for it to be 'false'
	//for c.Running {
	//	time.Sleep(100 * time.Millisecond)
	//}

	//if c.Report().Error != "" {
	//	fmt.Println(c.Report().Error)
	//} else {
	//	for _, entry := range c.Report().End.Streams {
	//		fmt.Println(entry.String())
	//	}
	//	for _, entry := range c.Report().ServerOutputJson.End.Streams {
	//		fmt.Println(entry.String())
	//	}
	//	fmt.Printf("DL Rate: %s\n", conversions.IntBitRateToString(int64(c.Report().End.SumReceived.BitsPerSecond)))
	//	fmt.Printf("UL Rate: %s\n", conversions.IntBitRateToString(int64(c.Report().End.SumSent.BitsPerSecond)))
	//}
}
