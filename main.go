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
	"time"

	"github.com/google/go-github/v39/github"
	"golang.org/x/oauth2"
)

// RSS 代表 XML 中的 <rss> 标签
type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Channel Channel  `xml:"channel"`
}

// Channel 代表 XML 中的 <channel> 标签
type Channel struct {
	Title       string `xml:"title"`
	Description string `xml:"description"`
	Link        string `xml:"link"`
	Items       []Item `xml:"item"`
}

// Item 代表 XML 中的 <item> 标签
type Item struct {
	Title   string    `xml:"title"`
	Link    string    `xml:"link"`
	PubDate time.Time `xml:"pubDate"`
}

func main() {
	// 设置 GitHub API 的 token 和 owner/repo
	token := os.Getenv("GITHUB_TOKEN")
	ownerRepo := strings.Split(os.Getenv("GITHUB_REPOSITORY"), "/")

	// 创建 GitHub 客户端
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	// 抓取 RSS 内容
	resp, err := http.Get("https://r.ifyes.online:6443/public.php?op=rss&id=-1&is_cat=0&q=&key=xt9z0r6325e20d77277")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// 解析 RSS 内容
	rssBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	var rss RSS
	err = xml.Unmarshal(rssBytes, &rss)
	if err != nil {
		log.Fatal(err)
	}

	// 创建文件
	for _, item := range rss.Channel.Items {
		// 文件名以 title 命名
		filename := fmt.Sprintf("%s.html", item.Title)
		fileContent := fmt.Sprintf("<html><head><title>%s</title></head><body><h1>%s</h1><p><a href=\"%s\">%s</a></p></body></html>", item.Title, item.Title, item.Link, item.Link)

		// 创建文件并保存到 GitHub 仓库
		opt := &github.RepositoryContentFileOptions{
			Message: github.String(fmt.Sprintf("Add file %s", filename)),
			Content: []byte(fileContent),
		}
		_, _, err = client.Repositories.CreateFile(ctx, ownerRepo[0], ownerRepo[1], fmt.Sprintf("rss-stared/%s", filename), opt)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("Created file %s\n", filename)
	}
}
