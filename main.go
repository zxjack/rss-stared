package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"unicode"

	"github.com/google/go-github/v37/github"
	"golang.org/x/oauth2"
)

const (
	baseURL        = "https://r.ifyes.online:6443"
	feedPath       = "/public.php?op=rss&id=-1&is_cat=0&q=&key=xt9z0r6325e20d77277"
	repoOwner      = "zxjack"
	repoName       = "rss-stared"
	defaultContent = `<!doctype html><html><head><meta charset="utf-8"><title>%s</title></head><body>%s</body></html>`
)

type AtomFeed struct {
	XMLName xml.Name `xml:"http://www.w3.org/2005/Atom feed"`
	Title   string   `xml:"title"`
	Entries []struct {
		Title   string `xml:"title"`
		Content string `xml:"content"`
	} `xml:"entry"`
}

func main() {
	// Create a GitHub client using a personal access token
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		log.Fatal("GitHub token not found in environment variables")
	}
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	// Get the Atom feed from the source URL
	resp, err := http.Get(baseURL + feedPath)
	if err != nil {
		log.Fatalf("Failed to fetch feed: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response body: %v", err)
	}

	// Parse the Atom feed XML
	var atomFeed AtomFeed
	err = xml.Unmarshal(body, &atomFeed)
	if err != nil {
		log.Fatalf("Failed to unmarshal Atom feed XML: %v", err)
	}

	// Iterate over the feed entries and create HTML files for each one
	for _, entry := range atomFeed.Entries {
		// Clean up the title and use it as the file name
		// fileName := strings.ReplaceAll(entry.Title, " ", "_") + ".html"
		// Clean up the title and use it as the file name
		cleanedTitle := strings.Map(func(r rune) rune {
			if unicode.IsLetter(r) || unicode.IsNumber(r) || r == '-' || r == '_' {
				return r
			}
			return '_'
		}, entry.Title)

		fileName := cleanedTitle + ".html"
		fileContent := fmt.Sprintf(defaultContent, entry.Title, entry.Content)

		// Check if the file already exists in the repository
		// Use the file's SHA to determine if it's a duplicate
		contents, _, _, err := client.Repositories.GetContents(ctx, repoOwner, repoName, fileName, &github.RepositoryContentGetOptions{})
		if err == nil {
			if contents != nil && contents.SHA != nil {
				fmt.Printf("File already exists: %s (SHA: %s)\n", fileName, *contents.SHA)
				continue
			}
		}

		// Create a new file in the GitHub repository
		fileOpts := &github.RepositoryContentFileOptions{
			Message: github.String("Add new file: " + fileName),
			Content: []byte(fileContent),
			Branch:  github.String("main"),
		}
		_, _, err = client.Repositories.CreateFile(ctx, repoOwner, repoName, fileName, fileOpts)
		if err != nil {
			log.Fatalf("Failed to create file: %v", err)
		}
		fmt.Printf("Created file: %s\n", fileName)
	}
}
