package cmd

import (
	"context"
	"log"
	"time"

	"github.com/google/go-github/v52/github"
	"github.com/spf13/cobra"
)

func init() {

	prGet.Flags().StringVar(&GithubOrgName, "githubOrgName", "", "name of the organization the repos will be deleted from")
	prGet.Flags().StringVar(&GithubRepo, "githubRepoName", "", "name of the repository to delete the branch from")
	prGet.Flags().StringVar(&GithubBranchName, "branchName", "", "the name of the branch to delete")
	prGet.MarkFlagRequired("githubOrgName")
	prGet.MarkFlagRequired("githubRepoName")
	prGet.MarkFlagRequired("branchName")

	prGet.Run = func(cmd *cobra.Command, args []string) {
		if err := GetPR(); err != nil {
			log.Fatalf("error when deleting branch: %v", err)
		}
	}

	prMerge.Flags().StringVar(&GithubOrgName, "githubOrgName", "", "name of the organization the repos will be deleted from")
	prMerge.Flags().StringVar(&GithubRepo, "githubRepoName", "", "name of the repository to delete the branch from")
	prMerge.Flags().StringVar(&GithubBranchName, "branchName", "", "the name of the branch to delete")
	prMerge.MarkFlagRequired("githubOrgName")
	prMerge.MarkFlagRequired("githubRepoName")
	prMerge.MarkFlagRequired("branchName")

	prMerge.Run = func(cmd *cobra.Command, args []string) {
		if err := MergePR(); err != nil {
			log.Fatalf("error when merging PR: %v", err)
		}
	}

	prComment.Flags().StringVar(&GithubOrgName, "githubOrgName", "", "name of the organization the repos will be deleted from")
	prComment.Flags().StringVar(&GithubRepo, "githubRepoName", "", "name of the repository to delete the branch from")
	prComment.MarkFlagRequired("githubOrgName")
	prComment.MarkFlagRequired("githubRepoName")

	prComment.Run = func(cmd *cobra.Command, args []string) {
		if err := CommentPR(); err != nil {
			log.Fatalf("error when merging PR: %v", err)
		}
	}
}

func GetPR() error {

	ctx := context.Background()

	var prNumber int
	list, _, err := GithubClient.PullRequests.List(ctx, GithubOrgName, GithubRepo, &github.PullRequestListOptions{})
	if err != nil {
		return err
	}
	for _, pr := range list {
		if pr.Head.GetRef() == GithubBranchName {
			prNumber = pr.GetNumber()
		}
	}

	log.Println("pr number:", prNumber)

	since := time.Now().Add(-10 * time.Minute)
	comments, _, err := GithubClient.Issues.ListComments(ctx, GithubOrgName, GithubRepo, prNumber,
		&github.IssueListCommentsOptions{Sort: github.String("created"), Since: &since},
	)
	if err != nil {
		return err
	}

	log.Printf("comments for org: %s, repo: %s, branch: %s, pr number: %d", GithubOrgName, GithubRepo, GithubBranchName, prNumber)
	for _, c := range comments {
		log.Println(c.GetBody())
	}
	return nil
}

func MergePR() error {

	ctx := context.Background()

	var prNumber int
	list, _, err := GithubClient.PullRequests.List(ctx, GithubOrgName, GithubRepo, &github.PullRequestListOptions{})
	if err != nil {
		return err
	}
	for _, pr := range list {
		if pr.Head.GetRef() == GithubBranchName {
			prNumber = pr.GetNumber()
		}
	}

	mergeResult, _, err := GithubClient.PullRequests.Merge(ctx, GithubOrgName, GithubRepo, prNumber, "commit message", &github.PullRequestOptions{})
	if err != nil {
		return err
	}

	log.Printf("pr %d for repo %s merge result: %+v\n", prNumber, GithubRepo, mergeResult)

	return nil
}

func CommentPR() error {

	ctx := context.Background()

	var prNumber int = 241
	_, _, err := GithubClient.Issues.CreateComment(ctx, GithubOrgName, GithubRepo, prNumber, &github.IssueComment{Body: github.String("/retest")})
	if err != nil {
		return err
	}

	return nil
}
