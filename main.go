package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type Feed struct {
	XMLName xml.Name `xml:"feed"`
	Entries []Entry  `xml:"entry"`
}

type Entry struct {
	Title   string `xml:"title"`
	Link    Link   `xml:"link"`
	Content string `xml:"content"`
}

type Link struct {
	Href string `xml:"href,attr"`
}

func main() {
	// Get GitHub token from environment variable
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		log.Fatal("GITHUB_TOKEN environment variable is not set")
	}

	// Authenticate with GitHub API
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	// Get RSS feed
	resp, err := client.Get(context.Background(), "https://r.ifyes.online:6443/public.php?op=rss&id=-1&is_cat=0&q=&key=xt9z0r6325e20d77277", nil)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Parse RSS feed
	var feed Feed
	err = xml.Unmarshal(body, &feed)
	if err != nil {
		log.Fatal(err)
	}

	// Loop through entries and create HTML file for each one
	for _, entry := range feed.Entries {
		title := entry.Title
		content := entry.Content

		// Remove invalid characters from title for file name
		title = strings.ReplaceAll(title, " ", "-")
		title = strings.ReplaceAll(title, "/", "-")
		title = strings.ReplaceAll(title, "\\", "-")
		title = strings.ReplaceAll(title, ":", "-")
		title = strings.ReplaceAll(title, "*", "-")
		title = strings.ReplaceAll(title, "?", "-")
		title = strings.ReplaceAll(title, "\"", "-")
		title = strings.ReplaceAll(title, "<", "-")
		title = strings.ReplaceAll(title, ">", "-")
		title = strings.ReplaceAll(title, "|", "-")

		// Create HTML file
		fileName := title + ".html"
		file, err := os.Create(fileName)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		// Write content to file
		_, err = file.WriteString(content)
		if err != nil {
			log.Fatal(err)
		}

		// Commit file to GitHub repository
		commitMessage := "Add " + fileName
		opts := &github.RepositoryContentFileOptions{
			Message: &commitMessage,
			Content: []byte(content),
		}
		_, _, err = client.Repositories.CreateFile(ctx, "zxjack", "rss-stared", fileName, opts)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Created %s\n", fileName)
	}
}
