package main

import (
	"flag"
	"github.com/Nrehearsal/go_captive_portal/config"
	"github.com/Nrehearsal/go_captive_portal/environment"
	"github.com/Nrehearsal/go_captive_portal/signal"
	"github.com/Nrehearsal/go_captive_portal/webserver"
	"log"
	"time"
)

func main() {
	log.SetFlags(log.Lshortfile | log.Ltime)

	configFile := flag.String("c", "config.json", "path of configuration file")
	flag.Parse()

	err := config.ParseConfigFile(*configFile)
	if err != nil {
		return
	}

	gwHttpConf := config.GetGatewayHttp()
	go webserver.Run(gwHttpConf)
	//make sure the http service runs successfully
	time.Sleep(5 * time.Second)

	err = environment.Init()
	if err != nil {
		return
	}

	keepRunning := make(chan int, 1)
	signal.ListenException(keepRunning)

	<-keepRunning
	return
}
