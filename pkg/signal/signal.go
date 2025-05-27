package signal

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/sirupsen/logrus"
)

var (
	onlyOneSignalHandler = make(chan struct{})
	shutdownHandler      chan os.Signal
	shutdownSignals      = []os.Signal{os.Interrupt, syscall.SIGTERM}
	shutdownCallbacks    = make([]func() error, 0)
	wg                   = sync.WaitGroup{}
)

// SetupSignalContext is same as SetupSignalHandler, but a context.Context is returned.
// Only one of SetupSignalContext and SetupSignalHandler should be called, and only can
// be called once.
func SetupSignalContext() context.Context {
	close(onlyOneSignalHandler) // panics when called twice

	shutdownHandler = make(chan os.Signal, 2)

	ctx, cancel := context.WithCancel(context.Background())
	signal.Notify(shutdownHandler, shutdownSignals...)
	go func() {
		<-shutdownHandler
		wg.Add(1)
		cancel()
		fmt.Println()
		logrus.Infof("Stopping MCP server")
		if len(shutdownCallbacks) != 0 {
			for _, cb := range shutdownCallbacks {
				if cb == nil {
					continue
				}
				if err := cb(); err != nil {
					logrus.Error(err)
				}
			}
		}
		wg.Done()
		<-shutdownHandler

		// second signal. Exit directly.
		logrus.Warnf("forced to stop.")
		os.Exit(130)
	}()

	return ctx
}

func RegisterOnShutdown(cb func() error) {
	shutdownCallbacks = append(shutdownCallbacks, cb)
}

func Flush() {
	wg.Wait()
}
