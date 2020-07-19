package main

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/unsocket/unsocket"
	"os"
	"os/signal"
)

func main() {
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

	var verbose bool

	cmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")

	err := cmd.Execute()
	if err != nil {

	}
}

func run(cmd *cobra.Command, args []string) error {
	unsock, err := unsocket.NewUnsocket(&unsocket.Config{
		WebhookURL: args[0],
	})
	if err != nil {
		return fmt.Errorf("unable to create unsocket: %w", err)
	}

	// handle interrupt signals
	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, os.Interrupt)
		<-signals
		fmt.Println("interrupt received")
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
