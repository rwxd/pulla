package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/google/go-github/v48/github"
	"github.com/rwxd/pulla/utils"
	"golang.org/x/oauth2"
)

var (
	userFlag           = flag.String("user", "", "User to backup")
	tokenFlag          = flag.String("token", "", "User token")
	destinationFlag    = flag.String("dest", "", "Destination directory")
	workerFlag         = flag.Int("worker", 5, "Number of concurrent pulls")
	daemonFlag         = flag.Bool("daemon", false, "Run as daemon")
	daemonIntervalFlag = flag.Int("interval", 12, "Interval in hours")
)

func main() {
	flag.Parse()

	if len(*destinationFlag) == 0 {
		fmt.Println("No --dest provided")
		flag.CommandLine.Usage()
		os.Exit(1)
	}

	if len(*tokenFlag) == 0 {
		fmt.Println("No --token provided")
		flag.CommandLine.Usage()
		os.Exit(1)
	}

	if *daemonFlag {
		fmt.Printf("Running as daemon with interval of %d hours\n", *daemonIntervalFlag)
		for {
			fmt.Println("Starting backup")
			backup(*tokenFlag, *destinationFlag)
			fmt.Println("Backup finished")
			fmt.Printf("Sleeping for %d hours\n", *daemonIntervalFlag)
			time.Sleep(time.Duration(*daemonIntervalFlag) * time.Hour)
		}
	} else {
		backup(*tokenFlag, *destinationFlag)
	}
}

// Starts the backup
func backup(token string, destination string) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	userRepos, err := utils.GetUserRepos(ctx, client, *userFlag)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Found %d user repositories\n", len(userRepos))

	starredRepos, err := utils.GetStarredRepositories(ctx, client, *userFlag)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Found %d starred repositories\n", len(starredRepos))

	for _, repo := range starredRepos {
		userRepos = append(userRepos, repo.Repository)
	}

	handleRepositories(token, destination, userRepos, *workerFlag)
}

// Clones or updates multiple repositories to the given path
func handleRepositories(token string, destPath string, repos []*github.Repository, worker int) {
	var wg sync.WaitGroup
	guard := make(chan int, worker)

	for _, repo := range repos {
		guard <- 1 // would block if guard channel is already filled
		wg.Add(1)

		clonePath := utils.GetBackupPath(destPath, *repo.Owner.Login, *repo.Name)
		url := utils.ReplaceTokenInCloneUrl(*repo.CloneURL, token)

		go func(url string, clonePath string) {
			defer wg.Done()

			if utils.CheckPathExists(clonePath) {
				err := utils.UpdateRemoteUrl(clonePath, url)
				if err != nil {
					fmt.Println(err)
				}

				err = utils.UpdateRepo(clonePath, false)
				if err != nil {
					fmt.Println(err)
				}
				<-guard // release guard
				return
			}

			err := utils.CloneRepo(url, clonePath)
			if err != nil {
				fmt.Println(err)
			}
			<-guard // release guard
		}(url, clonePath)
	}
	wg.Wait()
}
