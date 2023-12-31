package utils

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type SignalNotifier interface {
	Notify(c chan<- os.Signal, sig ...os.Signal)
}

func WaitForExitSignalThenCleanup(f func() error, ctx context.Context) { // creates pointer to function

	// waitExitThenCleanup(obj.Stop) -> obj is the actual object

	// Create a channel to receive OS signals
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)

	// Wait for a termination signal
	sig := <-signalChannel
	log.Printf("Received terminate signal: %s", sig.String())

	// Run the arbitrary function passed in
	err := f()
	if err != nil {
		log.Printf("Cleanup Function passed in returned an error:\n %v", err)
	} else {
		log.Printf("Exit Cleanup Ran Successfully")
	}
}
