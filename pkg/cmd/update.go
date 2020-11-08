package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/benjlevesque/ghr/pkg/config"
	"github.com/benjlevesque/ghr/pkg/gh"
	"github.com/benjlevesque/ghr/pkg/release"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func updateCommandCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	config, err := config.Load()
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	if len(args) > 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	results := make([]string, len(config.Items))
	for i, item := range config.Items {
		results[i] = item.Name
	}
	return results, cobra.ShellCompDirectiveNoFileComp
}

func makeUpdateCmd() *cobra.Command {
	command := &cobra.Command{
		Use:               "update [OWNER/REPO]",
		Short:             "updates an application",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: updateCommandCompletion,
		RunE:              runUpdate,
	}

	command.Flags().BoolP("force", "f", false, "forces the update")
	viper.BindPFlag("force", command.Flags().Lookup("force"))

	command.Flags().StringP("tag", "t", "latest", "specify the release tag")
	// workaround due to https://github.com/spf13/viper/issues/233
	viper.BindPFlag("update.tag", command.Flags().Lookup("tag"))

	command.RegisterFlagCompletionFunc("tag", tagArgsCompletion)

	return command
}

func runUpdate(cmd *cobra.Command, args []string) error {
	tag := viper.GetString("update.tag")

	config, err := config.Load()
	if err != nil {
		return err
	}
	item := config.Get(args[0])
	owner := item.Owner()
	repo := item.Repository()

	rel, err := gh.GetReleaseByTag(owner, repo, tag)
	if err != nil {
		return err
	}

	force := viper.GetBool("force")

	if item.Version == *rel.TagName && !force {
		fmt.Printf("Already up-to-date, version is %s\n", item.Version)
		return nil
	}

	releaseManager := &release.ReleaseManager{
		Owner: owner,
		Repo:  repo,
		Tag:   *rel.TagName,
	}

	return releaseManager.Install("", filepath.Dir(item.Executable))
}
