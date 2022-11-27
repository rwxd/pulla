package utils

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/google/go-github/v48/github"
)

func GetBackupPath(dir string, group string, name string) string {
	return filepath.Join(dir, group, name)
}

func GetUserRepos(ctx context.Context, client *github.Client, user string) ([]*github.Repository, error) {
	allRepos := []*github.Repository{}
	opts := github.RepositoryListOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}

	for {
		repos, resp, err := client.Repositories.List(ctx, user, &opts)
		if err != nil {
			return allRepos, err
		}

		allRepos = append(allRepos, repos...)
		if resp.NextPage == 0 {
			break
		}

		opts.Page = resp.NextPage
	}
	return allRepos, nil
}

func GetStarredRepositories(ctx context.Context, client *github.Client, user string) ([]*github.StarredRepository, error) {
	allRepos := []*github.StarredRepository{}
	opts := github.ActivityListStarredOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}

	for {
		repos, resp, err := client.Activity.ListStarred(ctx, user, &opts)
		if err != nil {
			return allRepos, err
		}

		allRepos = append(allRepos, repos...)
		if resp.NextPage == 0 {
			break
		}

		opts.Page = resp.NextPage
	}
	return allRepos, nil
}

func CloneRepo(cloneUrl string, dest string) error {
	fmt.Printf("Tracking %s\n", dest)
	cmd := exec.Command("git", "clone", cloneUrl, dest)
	out, err := cmd.Output()
	if err != nil {
		log.Println(string(out))
		return err
	}

	return nil
}

func CheckPathExists(dest string) bool {
	if _, err := os.Stat(dest); os.IsNotExist(err) {
		return false
	}

	return true
}

func UpdateRepo(path string) error {
	fmt.Printf("Updating %s\n", path)
	cmd := exec.Command("git", "-C", path, "pull", "--force")
	out, err := cmd.Output()
	if err != nil {
		log.Println(string(out))
		return err
	}

	return nil
}

func UpdateRemoteUrl(path string, url string) error {
	cmd := exec.Command("git", "-C", path, "remote", "set-url", "origin", url)
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func ReplaceTokenInCloneUrl(url string, token string) string {
	splitted := strings.Split(url, "://")
	return fmt.Sprintf("https://oauth2:%s@%s", token, splitted[1])
}
