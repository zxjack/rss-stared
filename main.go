package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"text/template"
	"time"

	"github.com/google/go-github/v32/github"
	"golang.org/x/oauth2"
)

// RSS .
type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Channel Channel  `xml:"channel"`
}

// Channel .
type Channel struct {
	Title       string  `xml:"title"`
	Description string  `xml:"description"`
	Link        string  `xml:"link"`
	Items       []Item  `xml:"item"`
	LastUpdate  RssTime `xml:"lastBuildDate"`
}

// Item .
type Item struct {
	Title       string  `xml:"title"`
	Description CDATA   `xml:"description"`
	Link        string  `xml:"link"`
	PubDate     RssTime `xml:"pubDate"`
}

// CDATA .
type CDATA struct {
	Data string `xml:",cdata"`
}

// RssTime .
type RssTime struct {
	time.Time
}

// UnmarshalXML .
func (t *RssTime) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	const shortForm = "Mon, 02 Jan 2006 15:04:05 -0700"
	var v string
	d.DecodeElement(&v, &start)
	dt, err := time.Parse(shortForm, v)
	if err != nil {
		return err
	}
	*t = RssTime{dt}
	return nil
}

// readRSS .
func readRSS(url string) ([]Item, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var rss RSS
	err = xml.Unmarshal(body, &rss)
	if err != nil {
		return nil, err
	}

	return rss.Channel.Items, nil
}

func main() {
	rssItems, err := readRSS("https://r.ifyes.online:6443/public.php?op=rss&id=-1&is_cat=0&q=&key=xt9z0r6325e20d77277")
	if err != nil {
		log.Fatal(err)
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)

	client := github.NewClient(tc)
	repoOwner := "zxjack"
	repoName := "rss-stared"
	dir := "articles"
	author := &github.CommitAuthor{
		Name:  github.String("GitHub Actions"),
		Email: github.String("actions@github.com"),
	}

	for _, rssItem := range rssItems {
		fileName := fmt.Sprintf("%s.html", rssItem.Title)
		filePath := filepath.Join(dir, fileName)

		tmpl, err := template.New("article").Parse(`
<!doctype html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>{{ .Title }}</title>
</head>
<body>
    <h1>{{ .Title }}</h1>
    {{ .Description }}
</body>
</html>
`)
		if err != nil {
			log.Fatal(err)
		}

		htmlFile, err := os.Create(fileName)
		if err != nil {
			log.Fatal(err)
		}

		err = tmpl.Execute(htmlFile, rssItem)
		if err != nil {
			log.Fatal(err)
		}

		htmlFile.Close()

		fileContent, err := ioutil.ReadFile(filePath)
		if err != nil {
			log.Fatal(err)
		}

		commitMessage := fmt.Sprintf("Add article: %s", fileName)
		commit := &github.RepositoryContentFileOptions{
			Message: &commitMessage,
			Content: fileContent,
			Branch:  github.String("main"),
			Author:  author,
		}

		_, _, err = client.Repositories.CreateFile(context.Background(), repoOwner, repoName, filePath, commit)
		if err != nil {
			log.Fatal(err)
		}
	}
}