package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/google/go-github/v38/github"
	"golang.org/x/oauth2"
)

func main() {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		log.Fatal("GITHUB_TOKEN not found in environment")
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	// 读取文件
	file, err := os.Open("rss-stared/解决WIN11最让人不爽的2个问题-右键和任务栏.html")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	content, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	// 创建文件
	repo := "zxjack/rss-stared"
	owner := "zxjack"
	branch := "main"
	path := "rss-stared"

	opts := &github.RepositoryContentFileOptions{
		Message: github.String("Add file"),
		Content: content,
	}

	_, _, err = client.Repositories.CreateFile(ctx, owner, repo, path, opts, &github.RepositoryContentFileOptions{
		Branch: github.String(branch),
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("File added successfully")
}
