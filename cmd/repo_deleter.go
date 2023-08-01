package cmd

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/google/go-github/v52/github"
	"github.com/spf13/cobra"
)

func init() {

	repoDelete.Flags().StringVar(&GithubOrgName, "githubOrgName", "", "name of the organization the repos will be deleted from")
	repoDelete.Flags().StringVar(&RepoFilter, "repoFilter", "", "the filter used for a selection of the repos that will be deleted")
	repoDelete.MarkFlagRequired("githubOrgName")
	repoDelete.MarkFlagRequired("repoFilter")

	repoDelete.Run = func(cmd *cobra.Command, args []string) {
		if err := DeleteRepo(); err != nil {
			log.Fatalf("error when deleting repo: %v", err)
		}
	}
}

func DeleteRepo() error {

	ctx := context.Background()

	reps, _, err := GithubClient.Repositories.ListByOrg(ctx, GithubOrgName, &github.RepositoryListByOrgOptions{Type: "all", ListOptions: github.ListOptions{PerPage: 500}})
	if err != nil {
		return fmt.Errorf("error when listing repositories: %v", err)
	}

	fmt.Printf("total number of repos in org '%s': %d\n", GithubOrgName, len(reps))
	for _, repo := range reps {
		if strings.Contains(repo.GetName(), RepoFilter) {
			fmt.Printf("about to delete a repo '%s' from org '%s'\n", repo.GetName(), GithubOrgName)
			_, err := GithubClient.Repositories.Delete(ctx, GithubOrgName, repo.GetName())
			if err != nil {
				return err
			}
			fmt.Printf("repository '%s' deleted successfully\n", repo.GetName())
		}
	}
	return nil
}
