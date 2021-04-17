package cmd

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/rijdendetreinen/ndov-receiver/output"
	"github.com/rijdendetreinen/ndov-receiver/receiver"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var exitReceiverChannel = make(chan bool)

var listenCmd = &cobra.Command{
	Use:   "listen",
	Short: "Listen for new messages",
	Run: func(cmd *cobra.Command, args []string) {
		initLogger(cmd)

		log.Infof("%v starting", RootCmd.Version)

		signalChan := make(chan os.Signal, 1)
		shutdownFinished := make(chan struct{})

		signal.Notify(signalChan, os.Interrupt)
		signal.Notify(signalChan, syscall.SIGTERM)

		go func() {
			sig := <-signalChan
			log.Errorf("Received signal: %+v, shutting down", sig)
			signal.Reset()
			shutdown()
			close(shutdownFinished)
		}()

		output.SetupOutputs()

		go receiver.ReceiveData(exitReceiverChannel)

		<-shutdownFinished
		log.Error("Exiting")
	},
}

func init() {
	RootCmd.AddCommand(listenCmd)
}

func initLogger(cmd *cobra.Command) {
	// TODO: setup logger

	if verbose, _ := cmd.Flags().GetBool("verbose"); verbose {
		log.SetLevel(log.DebugLevel)
		log.Debug("Verbose logging enabled")
	}
}

func shutdown() {
	log.Warn("Shutting down")
	exitReceiverChannel <- true
	<-exitReceiverChannel
}
