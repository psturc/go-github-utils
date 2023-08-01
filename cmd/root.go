package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/google/go-github/v52/github"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

var (
	GithubTokenKey string = "github_token"
	GithubToken    string
	GithubClient   *github.Client

	GithubOrgName string
	GithubRepo    string
	RepoFilter    string

	FilePath    string
	FileContent string

	GithubBranchName      string
	GithubNewBranchName   string
	GithubBaseBranchName  string
	GithubBaseBranchSHA   string
	GithubBranchNameRegex string
)

var rootCmd = &cobra.Command{
	Use:   "ggh",
	Short: "ggh helper to do github stuff via cli",
}

var repoDelete = &cobra.Command{
	Use:   "repo-delete",
	Short: "Delete Github repo from org",
	// Run: func(cmd *cobra.Command, args []string) {
	// },
}

var branchDelete = &cobra.Command{
	Use:   "branch-delete",
	Short: "Delete Github branch from the repo",
	// Run: func(cmd *cobra.Command, args []string) {
	// },
}

var branchCreate = &cobra.Command{
	Use:   "branch-create",
	Short: "Create Github branch in the repo",
	// Run: func(cmd *cobra.Command, args []string) {
	// },
}

var branchList = &cobra.Command{
	Use:   "branch-list",
	Short: "List branches from the GitHub repo",
	// Run: func(cmd *cobra.Command, args []string) {
	// },
}

var branchListChecks = &cobra.Command{
	Use:   "branch-list-checks",
	Short: "List checks for Github branch from the repo",
	// Run: func(cmd *cobra.Command, args []string) {
	// },
}

var prGet = &cobra.Command{
	Use:   "pr-get",
	Short: "Get Github PR from the repo",
	// Run: func(cmd *cobra.Command, args []string) {
	// },
}

var prMerge = &cobra.Command{
	Use:   "pr-merge",
	Short: "Merge Github PR created from a specified branch",
	// Run: func(cmd *cobra.Command, args []string) {
	// },
}

var prComment = &cobra.Command{
	Use:   "pr-comment",
	Short: "Comment on Github PR",
	// Run: func(cmd *cobra.Command, args []string) {
	// },
}

var webhookConfig = &cobra.Command{
	Use:   "webhook-config",
	Short: "Configure Github webhook for a repo",
	// Run: func(cmd *cobra.Command, args []string) {
	// },
}

var webhookList = &cobra.Command{
	Use:   "webhook-list",
	Short: "List Github webhooks for a repo",
	// Run: func(cmd *cobra.Command, args []string) {
	// },
}

var fileCreate = &cobra.Command{
	Use:   "file-create",
	Short: "create a file in github repo",
	// Run: func(cmd *cobra.Command, args []string) {
	// },
}

var fileUpdate = &cobra.Command{
	Use:   "file-update",
	Short: "push a change to a github repo",
	// Run: func(cmd *cobra.Command, args []string) {
	// },
}

var fileDelete = &cobra.Command{
	Use:   "file-delete",
	Short: "delete a file from a github repo",
	// Run: func(cmd *cobra.Command, args []string) {
	// },
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(repoDelete)
	rootCmd.AddCommand(webhookConfig)
	rootCmd.AddCommand(webhookList)
	rootCmd.AddCommand(fileCreate)
	rootCmd.AddCommand(fileUpdate)
	rootCmd.AddCommand(fileDelete)
	rootCmd.AddCommand(branchDelete)
	rootCmd.AddCommand(branchCreate)
	rootCmd.AddCommand(prGet)
	rootCmd.AddCommand(prMerge)
	rootCmd.AddCommand(prComment)
	rootCmd.AddCommand(branchListChecks)
	rootCmd.AddCommand(branchList)

	rootCmd.PersistentFlags().StringVarP(&GithubToken, "token", "t", "", fmt.Sprintf("Github access token. Can be set via the %s env var.", strings.ToUpper(GithubTokenKey)))
	viper.BindPFlag(GithubTokenKey, rootCmd.PersistentFlags().Lookup("token"))

	cobra.OnInitialize(initGithubClient)
}

func initGithubClient() {
	viper.AutomaticEnv() // read in environment variables that match
	token := viper.GetString(GithubTokenKey)
	if token == "" {
		log.Fatalln("Github token not defined. See usage.")
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)

	GithubClient = github.NewClient(tc)
}
