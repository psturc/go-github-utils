package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/google/go-github/v52/github"
	"github.com/spf13/cobra"
)

func init() {

	fileCreate.Flags().StringVar(&GithubOrgName, "githubOrgName", "", "name of the github organization")
	fileCreate.Flags().StringVar(&GithubRepo, "githubRepoName", "", "name of the repo where the hook should be set up")
	fileCreate.Flags().StringVar(&FilePath, "filePath", "", "path to the file that should be updated")
	fileCreate.Flags().StringVar(&FileContent, "fileContent", "", "new content of the file")
	fileCreate.Flags().StringVar(&GithubBranchName, "branchName", "", "the name of the branch to update")
	fileCreate.MarkFlagRequired("githubOrgName")
	fileCreate.MarkFlagRequired("githubRepoName")
	fileCreate.MarkFlagRequired("filePath")
	fileCreate.MarkFlagRequired("fileContent")

	fileCreate.Run = func(cmd *cobra.Command, args []string) {
		if err := UpdateFile(); err != nil {
			log.Fatalf("error when creating a file in a github repo: %v", err)
		}
	}

	fileUpdate.Flags().StringVar(&GithubOrgName, "githubOrgName", "", "name of the github organization")
	fileUpdate.Flags().StringVar(&GithubRepo, "githubRepoName", "", "name of the repo where the hook should be set up")
	fileUpdate.Flags().StringVar(&FilePath, "filePath", "", "path to the file that should be updated")
	fileUpdate.Flags().StringVar(&FileContent, "fileContent", "", "new content of the file")
	fileUpdate.Flags().StringVar(&GithubBranchName, "branchName", "", "the name of the branch to update")
	fileUpdate.MarkFlagRequired("githubOrgName")
	fileUpdate.MarkFlagRequired("githubRepoName")
	fileUpdate.MarkFlagRequired("filePath")
	fileUpdate.MarkFlagRequired("fileContent")

	fileUpdate.Run = func(cmd *cobra.Command, args []string) {
		if err := UpdateFile(); err != nil {
			log.Fatalf("error when committing to a github repo: %v", err)
		}
	}

	fileDelete.Flags().StringVar(&GithubOrgName, "githubOrgName", "", "name of the github organization")
	fileDelete.Flags().StringVar(&GithubRepo, "githubRepoName", "", "name of the repo where the hook should be set up")
	fileDelete.Flags().StringVar(&FilePath, "filePath", "", "path to the file that should be updated")
	fileDelete.Flags().StringVar(&GithubBranchName, "branchName", "", "the name of the branch to update")
	fileDelete.MarkFlagRequired("githubOrgName")
	fileDelete.MarkFlagRequired("githubRepoName")
	fileDelete.MarkFlagRequired("filePath")

	fileDelete.Run = func(cmd *cobra.Command, args []string) {
		if err := DeleteFile(); err != nil {
			log.Fatalf("error when deleting a file from github repo: %v", err)
		}
	}
}

func UpdateFile() error {

	ctx := context.Background()

	opts := &github.RepositoryContentGetOptions{}
	if GithubBranchName != "" {
		opts.Ref = fmt.Sprintf("heads/%s", GithubBranchName)
	}
	file, _, _, err := GithubClient.Repositories.GetContents(ctx, GithubOrgName, GithubRepo, "README.md", opts)
	if err != nil {
		return fmt.Errorf("error when listing file contents: %v", err)
	}

	fileSha := file.GetSHA()

	update := &github.RepositoryContentFileOptions{
		Message: github.String("test commit message"),
		SHA:     github.String(fileSha),
		Content: []byte(FileContent),
		Branch:  github.String(GithubBranchName),
	}
	_, _, err = GithubClient.Repositories.UpdateFile(ctx, GithubOrgName, GithubRepo, FilePath, update)
	if err != nil {
		return fmt.Errorf("error when updating a file on github: %v", err)
	}

	return nil
}

func DeleteFile() error {

	getOpts := &github.RepositoryContentGetOptions{}
	if GithubBranchName != "" {
		getOpts.Ref = fmt.Sprintf("heads/%s", GithubBranchName)
	}
	file, _, resp, err := GithubClient.Repositories.GetContents(context.Background(), GithubOrgName, GithubRepo, FilePath, getOpts)
	if err != nil {
		log.Printf("resp content: %+v", resp.StatusCode)
		return fmt.Errorf("error when listing file contents: %v", err)
	}

	deleteOpts := &github.RepositoryContentFileOptions{
		Message: github.String("test delete"),
		SHA:     github.String(file.GetSHA()),
	}
	if GithubBranchName != "" {
		deleteOpts.Branch = github.String(GithubBranchName)
	}
	contentResp, _, err := GithubClient.Repositories.DeleteFile(context.Background(), GithubOrgName, GithubRepo, FilePath, deleteOpts)
	if err != nil {
		return fmt.Errorf("error when deleting file on github: %v", err)
	}
	log.Printf("content resp: %+v", contentResp)
	return nil
}

func CreateFile() error {
	opts := &github.RepositoryContentFileOptions{
		Message: github.String("e2e test commit message"),
		Content: []byte(FileContent),
		Branch:  github.String("rhtap-demo-component-cwac"),
	}

	_, _, err := GithubClient.Repositories.CreateFile(context.Background(), GithubOrgName, GithubRepo, FilePath, opts)
	if err != nil {
		return fmt.Errorf("error when creating file contents: %v", err)
	}

	return nil
}
