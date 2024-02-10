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

// Returns the filepath as string for a repository owned by someone
func GetBackupPath(dir string, owner string, name string) string {
	return filepath.Join(dir, owner, name)
}

// Get all owned repos of a given user
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

// Get all starred repositories of a given user
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

// Clone a repo to the file system
// using the git clone command
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

// Updates an existing git repository at the given path
func UpdateRepo(path string, force bool) error {
	fmt.Printf("Updating %s\n", path)
	cmdArgs := []string{"-C", path, "pull"}
	if force {
		cmdArgs = append(cmdArgs, "--force")
	}

	cmd := exec.Command("git", cmdArgs...)
	out, err := cmd.Output()
	if err != nil {
		log.Println(string(out))
		return err
	}

	return nil
}

// Updates the remote url of a repository
// with the given url
func UpdateRemoteUrl(path string, url string) error {
	cmd := exec.Command("git", "-C", path, "remote", "set-url", "origin", url)
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

// Adds oauth2 to a url by replacing parts of it
func ReplaceTokenInCloneUrl(url string, token string) string {
	splitted := strings.Split(url, "://")
	return fmt.Sprintf("https://oauth2:%s@%s", token, splitted[1])
}
