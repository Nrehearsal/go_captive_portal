package signal

import (
	"github.com/Nrehearsal/go_captive_portal/environment"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func ListenException(done chan int) {
	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGSEGV)

	sig := <-sigs
	//clean up the battlefield
	environment.Clean()

	log.Println("interrupted by signal: ", sig)
	done <- 1
}
