package main

import (
	"fmt"
	"github.com/BGrewell/go-iperf"
	"log"
)

func main() {

	// First create a controller for the "client" side (would be on a different computer)
	cc, err := iperf.NewController(6802)
	if err != nil {
		log.Fatal(err)
	}
	cc.Port = 6801 //Note: this is just a hack because we start a listener on the port when we get a new controller and in this case since we are on the same pc we would have a port conflict so we start our listener on the port+1 and then fix that after the listener has started since we won't use it anyway

	// Second create a controller for the "server" side (again would normally be on a different computer)
	sc, err := iperf.NewController(6801)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("[+] Server side controller listening on %d\n", sc.Port)

	// Test iperf
	iperfCli, err := cc.NewClient("127.0.0.1")
	if err != nil {
		log.Fatal(err)
	}

	err = iperfCli.Start()
	if err != nil {
		log.Fatal(err)
	}

	<- iperfCli.Done

	fmt.Println(iperfCli.Report().String())
}
