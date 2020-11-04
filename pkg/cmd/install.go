package cmd

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/benjlevesque/ghr/pkg/gh"
	"github.com/benjlevesque/ghr/pkg/release"
	"github.com/google/go-github/v32/github"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func parseArgs(args []string) (string, string, string, error) {
	var asset string
	if len(args) > 1 {
		asset = args[1]
	}
	if len(args) >= 1 {
		values := strings.Split(args[0], "/")
		if len(values) == 2 {
			return values[0], values[1], asset, nil
		}
	}
	return "", "", "", errors.New("cannot find OWNER/REPO in args")
}

func tagArgsCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	gh := github.NewClient(nil)
	owner, repo, _, err := parseArgs(args)
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	releases, _, err := gh.Repositories.ListReleases(context.TODO(), owner, repo, &github.ListOptions{})
	if err != nil {
		return nil, cobra.ShellCompDirectiveDefault
	}
	versions := make([]string, len(releases))
	for i, release := range releases {
		versions[i] = *release.TagName
	}
	return versions, cobra.ShellCompDirectiveNoFileComp
}

func installCommandCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	tag := viper.GetString("tag")
	owner, repo, _, err := parseArgs(args)
	if owner == "" && toComplete != "" {
		owner, _, _, err = parseArgs([]string{toComplete})
	}
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	if repo == "" {
		repos, err := gh.GetReposForOrg(owner)
		if err != nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		return repos, cobra.ShellCompDirectiveNoSpace
	}
	release, err := gh.GetReleaseByTag(owner, repo, tag)

	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	result := make([]string, len(release.Assets))
	for i, asset := range release.Assets {
		result[i] = *asset.Name
	}
	return result, cobra.ShellCompDirectiveNoFileComp
}

func makeInstallCmd() *cobra.Command {
	command := &cobra.Command{
		Use:               "install [OWNER/REPO] [asset]",
		Short:             "downloads & install an executable from a github release",
		Args:              cobra.RangeArgs(1, 2),
		ValidArgsFunction: installCommandCompletion,
		RunE:              runInstall,
	}
	command.Flags().StringP("tag", "t", "latest", "specify the release tag")
	// workaround due to https://github.com/spf13/viper/issues/233
	viper.BindPFlag("install.tag", command.Flags().Lookup("tag"))

	command.RegisterFlagCompletionFunc("tag", tagArgsCompletion)
	return command
}

func runInstall(cmd *cobra.Command, args []string) error {
	tag := viper.GetString("install.tag")
	fmt.Println(tag)

	owner, repo, assetName, err := parseArgs(args)
	if err != nil {
		return fmt.Errorf("Invalid arguments: %s", err)
	}
	releaseManager := &release.ReleaseManager{
		Owner: owner,
		Repo:  repo,
		Tag:   tag,
	}
	err = releaseManager.Install(assetName)
	if err != nil {
		return fmt.Errorf("Installation failed: %s", err)
	}
	return nil
}
