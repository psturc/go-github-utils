package cmd

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/go-github/v52/github"
	"github.com/spf13/cobra"
)

func init() {

	webhookConfig.Flags().StringVar(&GithubOrgName, "githubOrgName", "", "name of the github organization")
	webhookConfig.Flags().StringVar(&GithubRepo, "githubRepo", "", "name of the repo where the hook should be set up")
	webhookConfig.MarkFlagRequired("githubOrgName")
	webhookConfig.MarkFlagRequired("githubRepo")

	webhookConfig.Run = func(cmd *cobra.Command, args []string) {
		if err := SetupWebhook(); err != nil {
			log.Fatalf("error when setting up webhook: %v", err)
		}
	}

	webhookList.Flags().StringVar(&GithubOrgName, "githubOrgName", "", "name of the github organization")
	webhookList.Flags().StringVar(&GithubRepo, "githubRepo", "", "name of the repo where the hook should be set up")
	webhookList.MarkFlagRequired("githubOrgName")
	webhookList.MarkFlagRequired("githubRepo")

	webhookList.Run = func(cmd *cobra.Command, args []string) {
		if err := ListWebhooks(); err != nil {
			log.Fatalf("error when listing webhooks: %v", err)
		}
	}
}

func ListWebhooks() error {
	hooks, _, err := GithubClient.Repositories.ListHooks(context.Background(), GithubOrgName, GithubRepo, &github.ListOptions{})
	if err != nil {
		return fmt.Errorf("error when listing webhooks: %+v", err)
	}

	url := hooks[0].Config["url"].(string)
	log.Println(url)
	return nil
}

func SetupWebhook() error {

	ctx := context.Background()

	hooks, _, err := GithubClient.Repositories.ListHooks(ctx, GithubOrgName, GithubRepo, &github.ListOptions{})
	if err != nil {
		return fmt.Errorf("error when listing webhooks: %v", err)
	}

	for _, hook := range hooks {
		createdAt := hook.GetCreatedAt()

		if createdAt.Before(time.Now().Add(time.Hour * -24)) {
			log.Printf("hook %s is older than a day, deleting...", hook.GetURL())
			log.Println(hook.Events)
			// TODO add delete function
			_, err := GithubClient.Repositories.DeleteHook(ctx, GithubOrgName, GithubRepo, *hook.ID)
			if err != nil {
				return fmt.Errorf("error when deleting webhook: %v", err)
			}
		}
	}

	newHookTemplate := &github.Hook{
		Active: github.Bool(true),
		Events: []string{"push"},
		Config: map[string]interface{}{
			"content_type": "json",
			"insecure_ssl": 0,
			"url":          "http://asdfasdfasjflakadlfkadjsflkjadasdfslkfjasdlkfjasdkfdjladsk.com",
		},
	}
	createdHook, _, err := GithubClient.Repositories.CreateHook(ctx, GithubOrgName, GithubRepo, newHookTemplate)
	if err != nil {
		return fmt.Errorf("error when creating webhook: %v", err)
	}

	log.Printf("webhook created: %s", *createdHook.URL)

	return nil
}
