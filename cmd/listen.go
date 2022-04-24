package cmd

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/evalphobia/logrus_sentry"
	"github.com/rijdendetreinen/ndov-receiver/output"
	"github.com/rijdendetreinen/ndov-receiver/receiver"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var exitReceiverChannel = make(chan bool)

var listenCmd = &cobra.Command{
	Use:   "listen",
	Short: "Listen for new messages",
	Run: func(cmd *cobra.Command, args []string) {
		initLogger(cmd)

		log.Infof("ndov-receiver %v starting", RootCmd.Version)

		signalChan := make(chan os.Signal, 1)
		shutdownFinished := make(chan struct{})

		signal.Notify(signalChan, os.Interrupt)
		signal.Notify(signalChan, syscall.SIGTERM)

		go func() {
			sig := <-signalChan
			log.Warnf("Received signal: %+v, shutting down", sig)
			signal.Reset()
			shutdown()
			close(shutdownFinished)
		}()

		output.SetupOutputs()

		go receiver.ReceiveData(exitReceiverChannel)

		<-shutdownFinished
		log.Warn("Exiting")
	},
}

func init() {
	RootCmd.AddCommand(listenCmd)
}

func initLogger(cmd *cobra.Command) {
	if verbose, _ := cmd.Flags().GetBool("verbose"); verbose {
		log.SetLevel(log.DebugLevel)
		log.Debug("Verbose logging enabled")
	}

	if viper.GetString("sentry.dsn") != "" {
		hook, err := logrus_sentry.NewSentryHook(viper.GetString("sentry.dsn"), []log.Level{
			log.PanicLevel,
			log.FatalLevel,
			log.ErrorLevel,
		})
		hook.SetRelease(Version.Version)

		if err == nil {
			log.AddHook(hook)
			log.WithField("dsn", viper.GetString("sentry.dsn")).Debug("Sentry logging enabled")
		} else {
			log.Error(err)
		}
	}
}

func shutdown() {
	log.Warn("Shutting down")
	exitReceiverChannel <- true
	<-exitReceiverChannel
}
