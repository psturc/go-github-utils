package cmd

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"regexp"
	"strconv"
	"sync"
	"time"

	"github.com/google/go-github/v52/github"
	"github.com/spf13/cobra"
)

func init() {

	branchDelete.Flags().StringVar(&GithubOrgName, "githubOrgName", "", "name of the organization the repos will be deleted from")
	branchDelete.Flags().StringVar(&GithubRepo, "githubRepoName", "", "name of the repository to delete the branch from")
	branchDelete.Flags().StringVar(&GithubBranchName, "branchName", "", "the name of the branch to delete")
	branchDelete.Flags().StringVar(&GithubBranchNameRegex, "regex", "", "regex to match for the branches about to be deleted")
	branchDelete.MarkFlagRequired("githubOrgName")
	branchDelete.MarkFlagRequired("githubRepoName")

	branchDelete.Run = func(cmd *cobra.Command, args []string) {
		if err := DeleteBranch(); err != nil {
			log.Fatalf("error when deleting branch: %v", err)
		}
	}

	branchCreate.Flags().StringVar(&GithubOrgName, "githubOrgName", "", "name of the organization where to create a branch")
	branchCreate.Flags().StringVar(&GithubRepo, "githubRepoName", "", "name of the repository to create the branch in")
	branchCreate.Flags().StringVar(&GithubNewBranchName, "newBranchName", "", "the name of the branch to create")
	branchCreate.Flags().StringVar(&GithubBaseBranchName, "baseBranchName", "", "the name of the branch the new branch will be based on")
	branchCreate.Flags().StringVar(&GithubBaseBranchSHA, "baseBranchSHA", "", "the commit sha of the base branch - optional")
	branchCreate.MarkFlagRequired("githubOrgName")
	branchCreate.MarkFlagRequired("githubRepoName")
	branchCreate.MarkFlagRequired("newBranchName")
	branchCreate.MarkFlagRequired("baseBranchName")

	branchCreate.Run = func(cmd *cobra.Command, args []string) {
		if err := CreateBranch(); err != nil {
			log.Fatalf("error when creating branch: %v", err)
		}
	}

	branchListChecks.Flags().StringVar(&GithubOrgName, "githubOrgName", "", "name of the organization the repos will be deleted from")
	branchListChecks.Flags().StringVar(&GithubRepo, "githubRepoName", "", "name of the repository to delete the branch from")
	branchListChecks.Flags().StringVar(&GithubBranchName, "branchName", "", "the name of the branch to delete")
	branchListChecks.MarkFlagRequired("githubOrgName")
	branchListChecks.MarkFlagRequired("githubRepoName")
	branchListChecks.MarkFlagRequired("branchName")

	branchListChecks.Run = func(cmd *cobra.Command, args []string) {
		if err := ListBranchChecks(); err != nil {
			log.Fatalf("error when deleting branch: %v", err)
		}
	}

	branchList.Flags().StringVar(&GithubOrgName, "githubOrgName", "", "name of the organization the repos will be deleted from")
	branchList.Flags().StringVar(&GithubRepo, "githubRepoName", "", "name of the repository to delete the branch from")
	branchList.MarkFlagRequired("githubOrgName")
	branchList.MarkFlagRequired("githubRepoName")

	branchList.Run = func(cmd *cobra.Command, args []string) {
		if err := ListBranches(); err != nil {
			log.Fatalf("error when listing branches: %v", err)
		}
	}
}

func DeleteBranch() error {

	ctx := context.Background()
	var branchesToDelete []*github.Branch
	var page int

	if GithubBranchNameRegex != "" {
		for {
			branches, res, err := GithubClient.Repositories.ListBranches(ctx, GithubOrgName, GithubRepo, &github.BranchListOptions{ListOptions: github.ListOptions{PerPage: 100, Page: page}})
			if err != nil {
				return err
			}

			for _, b := range branches {

				match, err := regexp.MatchString(GithubBranchNameRegex, b.GetName())
				if err != nil {
					return fmt.Errorf("problem with regexp: %+v", err)
				}
				if match {
					branchesToDelete = append(branchesToDelete, b)
				}
			}

			if res.NextPage != 0 {
				page = res.NextPage
			} else {
				break
			}
		}

	} else if GithubBranchName != "" {
		branchesToDelete = append(branchesToDelete, &github.Branch{Name: &GithubBaseBranchName})
	} else {
		return fmt.Errorf("none of the parameters 'regex' or 'branchName' specified")
	}

	log.Printf("got %d branches to delete\n", len(branchesToDelete))

	var wg sync.WaitGroup
	for _, b := range branchesToDelete {

		wg.Add(1)
		b := b
		go func() {
			defer wg.Done()
			log.Printf("deleting github branch %s\n", b.GetName())
			_, err := GithubClient.Git.DeleteRef(ctx, GithubOrgName, GithubRepo, fmt.Sprintf("heads/%s", b.GetName()))
			if err != nil {
				log.Fatal(err)
			}
		}()

	}
	wg.Wait()

	return nil
}

func CreateBranch() error {

	ctx := context.Background()

	ref, _, err := GithubClient.Git.GetRef(ctx, GithubOrgName, GithubRepo, fmt.Sprintf("heads/%s", GithubBaseBranchName))
	if err != nil {
		return fmt.Errorf("error getting base branch %s: %+v", GithubBaseBranchName, err)
	}
	ref.Ref = github.String("refs/heads/" + GithubNewBranchName)
	if GithubBaseBranchSHA != "" {
		ref.Object.SHA = github.String(GithubBaseBranchSHA)
	}
	log.Printf("%+v", ref)
	//os.Exit(0)
	for i := 1; i <= 100; i++ {
		branchName := strconv.Itoa(rand.Int())
		ref.Ref = github.String("refs/heads/" + branchName)
		refer, res, err := GithubClient.Git.CreateRef(ctx, GithubOrgName, GithubRepo, ref)
		log.Printf("%+v", res.Response.Status)
		log.Printf("%+v", refer)
		if err != nil {
			return err
		}

		for {
			_, _, err = GithubClient.Git.GetRef(ctx, GithubOrgName, GithubRepo, fmt.Sprintf("heads/%s", branchName))
			if err != nil {
				log.Printf("error getting branch %s: %+v", branchName, err)
				time.Sleep(time.Second)
				continue
			}
			break
		}

		opts := &github.RepositoryContentFileOptions{
			Message: github.String("e2e test commit message"),
			Content: []byte("blablabladsfas adsfasd "),
			Branch:  github.String(branchName),
		}

		content, res, err := GithubClient.Repositories.CreateFile(context.Background(), GithubOrgName, GithubRepo, "test.yaml", opts)
		if err != nil {
			return fmt.Errorf("error when creating file contents: %v", err)
		}
		log.Printf("%+v", res.Response.Status)
		log.Printf("%+v", content)
	}
	return nil
}

func ListBranchChecks() error {

	ctx := context.Background()

	list, _, err := GithubClient.Checks.ListCheckRunsForRef(ctx, GithubOrgName, GithubRepo, fmt.Sprintf("heads/%s", GithubBranchName), &github.ListCheckRunsOptions{})
	if err != nil {
		return err
	}
	for _, check := range list.CheckRuns {
		log.Println("branch check:", check)
	}

	return nil
}

func ListBranches() error {

	ctx := context.Background()

	branches, _, err := GithubClient.Repositories.ListBranches(ctx, GithubOrgName, GithubRepo, &github.BranchListOptions{ListOptions: github.ListOptions{PerPage: 500}})
	if err != nil {
		return err
	}

	for _, b := range branches {
		log.Printf("%+v", b.GetName())
	}

	return nil
}
