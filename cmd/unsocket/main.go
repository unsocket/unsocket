package main

import (
	"errors"
	"github.com/spf13/cobra"
	"github.com/unsocket/unsocket"
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
	err := unsocket.Unsocket(&unsocket.Config{
		WebhookURL: args[0],
	})
	if err != nil {
		return err
	}

	return nil
}
