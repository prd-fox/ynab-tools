package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/p2p/nat"
)

func register(port int) chan struct{} {


	upnpInterface := nat.Any() //nat.UPnP()

	ip, _ := upnpInterface.ExternalIP()
	fmt.Println(upnpInterface.String())

	closeChan := make(chan struct{})

	go nat.Map(upnpInterface, closeChan, "tcp", port+1000, port, "ynab tools")
	//nat.Map(upnpInterface, closeChan, "udp", port, port, "ynab tools")

	fmt.Println(ip.String())
	return closeChan
}
