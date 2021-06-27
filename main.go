package main

import (
	"flag"
	"fmt"
	"github.com/prd-fox/ynab-tools/log"
	"github.com/prd-fox/ynab-tools/ui"
	"gitlab.com/NebulousLabs/go-upnp"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prd-fox/ynab-tools/config"
)

func main() {

	port := uint16(50000)

	c := register(int(port))
	defer close(c)

	//d := upnpSetup(18000)
	//defer d.Clear(port)
	//defer d.Clear(3000)

	var configPath string
	flag.StringVar(&configPath, "configfile", "", "file path for the configuration")
	flag.Parse()

	if configPath == "" {
		fmt.Println("No config detected")
		os.Exit(1)
	}


	userConfig, err := config.LoadConfig(configPath)
	if err != nil {
		fmt.Println("err: " + err.Error())
		os.Exit(1)
	}

	fmt.Println(userConfig.ApiKey)

	uiHandler := ui.NewUIHandler(int(port))
	if err := uiHandler.Start(); err != nil {
		log.Error("Unable to start the UI", "err", err)
		return
	}
	defer uiHandler.Stop()

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(sigc)
	select {
	case <-sigc:
	}
	log.Info("Received interrupt signal, shutting down...")

	time.Sleep(1 * time.Second)

	return
}

func upnpSetup(port uint16) *upnp.IGD {
	// connect to router
	d, err := upnp.Load("http://192.168.0.1:5431/dyndev/uuid:00099eca-80ba-4642-9d85-6bbeb424bf48")
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	// discover external IP
	ip, err := d.ExternalIP()
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	fmt.Println("Your external IP is:", ip)

	// forward a port
	err = d.Forward(port, "upnp test")
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	fmt.Println(d.Location())

	return d
}