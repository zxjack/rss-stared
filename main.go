package main

import (
	"context"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/google/go-github/v37/github"
	"golang.org/x/oauth2"
)

type Item struct {
	XMLName     xml.Name `xml:"item"`
	Title       string   `xml:"title"`
	Link        string   `xml:"link"`
	Description string   `xml:"description"`
}

type Channel struct {
	XMLName     xml.Name `xml:"channel"`
	Title       string   `xml:"title"`
	Link        string   `xml:"link"`
	Description string   `xml:"description"`
	Items       []Item   `xml:"item"`
}

type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Channel Channel  `xml:"channel"`
}

type AtomEntry struct {
	Title   string `json:"title"`
	Link    string `json:"html_url"`
	Content string `json:"body"`
}

type AtomFeed struct {
	Title string       `json:"title"`
	Link  AtomLink     `json:"link"`
	Items []AtomEntry `json:"items"`
}

type AtomLink struct {
	Href string `json:"href"`
}

func main() {
	// Get starred repos from GitHub API
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	// Retrieve Atom feed
	opt := &github.ListOptions{PerPage: 100}
	var allEntries []AtomEntry
	for {
		repos, resp, err := client.Activity.ListStarred(ctx, "", opt)
		if err != nil {
			log.Fatal(err)
		}
		for _, repo := range repos {
			owner := repo.GetOwner().GetLogin()
			name := repo.GetName()
			entries, _, err := client.Repositories.ListCommits(ctx, owner, name, nil)
			if err != nil {
				log.Fatal(err)
			}
			for _, entry := range entries {
				entryURL := entry.GetHTMLURL()
				sha := entry.GetSHA()
				commit, _, err := client.Git.GetCommit(ctx, owner, name, sha)
				if err != nil {
					log.Fatal(err)
				}
				message := commit.GetMessage()
				if strings.Contains(message, "rss") {
					content := commit.Files[0].GetContent()
					if strings.Contains(content, "<rss") {
						var rss RSS
						xml.Unmarshal([]byte(content), &rss)
						for _, item := range rss.Channel.Items {
							allEntries = append(allEntries, AtomEntry{
								Title:   item.Title,
								Link:    item.Link,
								Content: item.Description,
							})
						}
					} else if strings.Contains(content, "<feed xmlns") {
						var atomFeed AtomFeed
						xml.Unmarshal([]byte(content), &atomFeed)
						allEntries = append(allEntries, atomFeed.Items...)
					}
				}
			}
		}
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	// Generate HTML file from Atom feed
	html := "<!DOCTYPE html>\n<html>\n<head>\n<meta charset=\"UTF-8\">\n<title>Starred RSS Articles</title>\n</head>\n<body>\n"
	for _, entry := range allEntries {
		html += fmt.Sprintf("<h2><a href=\"%s\">
