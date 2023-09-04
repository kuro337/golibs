package utils

import (
	"context"
	"os"
	"os/signal"
	"testing"
	"time"
)

func TestWaitForExitSignalThenCleanup(t *testing.T) {
	cleanupRan := false
	cleanupFn := func() error {
		cleanupRan = true
		return nil
	}

	// Create a context that can be canceled
	ctx, cancel := context.WithCancel(context.Background())

	// Create a channel to simulate receiving an OS signal
	signalChannel := make(chan os.Signal, 1)
	defer close(signalChannel)

	// Cancel the context after a short delay
	go func() {
		time.Sleep(50 * time.Millisecond)
		cancel()
	}()

	// Simulate receiving the cancel signal using signal.Notify
	signal.Notify(signalChannel, os.Interrupt)

	// Execute the WaitForExitSignalThenCleanup function with the cancellable context
	WaitForExitSignalThenCleanup(cleanupFn, ctx)

	// Wait for a short period to allow cleanup and logging to finish
	time.Sleep(100 * time.Millisecond)

	// Verify that the cleanup function ran
	if !cleanupRan {
		t.Error("Cleanup function did not run as expected")
	}
}
