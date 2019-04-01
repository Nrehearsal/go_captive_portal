package main

import (
	"flag"
	"log"

	"github.com/Nrehearsal/go_captive_portal/config"
	"github.com/Nrehearsal/go_captive_portal/signal"
	"github.com/Nrehearsal/go_captive_portal/webserver"
	"github.com/Nrehearsal/go_captive_portal/environment"
)

func main() {
	log.SetFlags(log.Lshortfile | log.Ltime)

	configFile := flag.String("c", "config.json", "path of configuration file")
	flag.Parse()

	err := config.ParseConfigFile(*configFile)
	if err != nil {
		return
	}

	done := make(chan bool)
	gwHttpConf := config.GetGatewayHttp()
	//make sure the http service runs successfully
	go func() {
		go webserver.Run(gwHttpConf)
		done <- true
	}()
	<-done

	err = environment.Init()
	if err != nil {
		return
	}

	keepRunning := make(chan int, 1)
	signal.ListenException(keepRunning)

	<-keepRunning
	return
}
