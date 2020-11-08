package cmd

import (
	"fmt"

	"github.com/benjlevesque/ghr/pkg/config"
	"github.com/benjlevesque/ghr/pkg/gh"
	"github.com/blang/semver"
	"github.com/spf13/cobra"
)

func makeCheckCmd() *cobra.Command {
	command := &cobra.Command{
		Use:   "check",
		Short: "lists available updates",
		Args:  cobra.NoArgs,
		RunE:  runCheck,
	}

	return command
}

func runCheck(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	for _, app := range cfg.Items {
		release, err := gh.GetReleaseByTag(app.Owner(), app.Repository(), "latest")
		if err != nil {
			fmt.Printf("Error getting latest release for %s: %s\n", app.Name, err)
			continue
		}
		currentver, err := semver.Make(app.Version)
		if err != nil {
			fmt.Printf("Not a semver: %s %s\n", app.Name, app.Version)
			continue
		}
		latest := *release.TagName
		latestver, err := semver.Make(latest)
		if err != nil {
			fmt.Printf("Not a semver: %s %s\n", app.Name, *release.TagName)
			continue
		}
		if currentver.Compare(latestver) < 0 {
			fmt.Printf(" - %s  %s => %s\n", app.Name, app.Version, latest)
		} else {
			fmt.Printf(" - %s  up to date\n", app.Name)
		}
	}

	return nil
}
