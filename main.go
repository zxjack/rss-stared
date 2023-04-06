package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/go-github/v38/github"
	"golang.org/x/oauth2"
)

func main() {
	// Authenticate with GitHub using a personal access token
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		log.Fatal("GITHUB_TOKEN not set")
	}
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(context.Background(), ts)
	client := github.NewClient(tc)

	// Fetch the contents of the RSS feed
	feedURL := "https://r.ifyes.online:6443/public.php?op=rss&id=-1&is_cat=0&q=&key=xt9z0r6325e20d77277"
	resp, err := http.Get(feedURL)
	if err != nil {
		log.Fatal("Error fetching feed:", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading response body:", err)
	}

	// Create a new commit with the RSS feed file
	title := fmt.Sprintf("RSS feed %s", time.Now().Format("2006-01-02 15:04:05"))
	message := "Automatically generated commit"
	path := "rss-stared/rss.xml"
	content := &github.RepositoryContentFileOptions{
		Message: &message,
		Content: body,
	}
	branch := "main"

	_, _, err = client.Repositories.CreateFile(context.Background(), "zxjack", "rss-stared", path, content, &github.RepositoryContentFileOptions{
		Branch: &branch,
	})

	if err != nil {
		log.Fatal("Error creating file in repository:", err)
	}

	fmt.Println("Successfully created file in repository:", path)
}
