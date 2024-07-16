package main

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/unsocket/unsocket"
	"os"
	"os/signal"
)

var (
	verbose = false
	webhookSecret = ""
)

func main() {
	log.SetFormatter(&log.JSONFormatter{})

	var cmd = &cobra.Command{
		Use:   "unsocket [webhook]",
		Short: "unsocket",
		Long:  `unsocket`,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errors.New("requires the webhook url")
			}

			return nil
		},
		RunE: run,
	}

	cmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	cmd.PersistentFlags().StringVar(&webhookSecret, "webhook-secret", envOrDefault("WEBHOOK_SECRET", webhookSecret), "webhook client secret")

	err := cmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true

	log.Infof("running %s", cmd.Name())

	if verbose {
		log.SetLevel(log.DebugLevel)
	}

	unsock, err := unsocket.NewUnsocket(&unsocket.Config{
		WebhookURL: args[0],
		WebhookSecret: webhookSecret,
	})
	if err != nil {
		return fmt.Errorf("unable to create unsocket: %w", err)
	}

	// handle interrupt signals
	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, os.Interrupt)
		s := <-signals
		log.Infof("interrupt %s received", s)
		_ = unsock.Stop()
	}()

	// run and block
	err = unsock.RunAndWait()
	if err != nil {
		return fmt.Errorf("unsocket failed: %w", err)
	}

	// signal successful execution
	return nil
}

func envOrDefault(key string, def string) string {
	value := os.Getenv(key)
	if value == "" {

		return def
	}
	return value
}
