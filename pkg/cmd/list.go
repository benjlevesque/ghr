package cmd

import (
	"fmt"

	"github.com/benjlevesque/ghr/pkg/config"
	"github.com/spf13/cobra"
)

func makeListCmd() *cobra.Command {
	command := &cobra.Command{
		Use:   "list",
		Short: "lists existing applications",
		Args:  cobra.NoArgs,
		RunE:  runList,
	}

	return command
}

func runList(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	for _, app := range cfg.Items {
		fmt.Printf(" - %s  (%s)\n", app.Name, app.Version)
	}

	return nil
}
