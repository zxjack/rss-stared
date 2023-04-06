package main

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type Rss struct {
	XMLName xml.Name `xml:"rss"`
	Channel Channel  `xml:"channel"`
}

type Channel struct {
	XMLName xml.Name `xml:"channel"`
	Items   []Item   `xml:"item"`
}

type Item struct {
	XMLName     xml.Name `xml:"item"`
	Title       string   `xml:"title"`
	Link        string   `xml:"link"`
	Description string   `xml:"description"`
}

func main() {
	resp, err := http.Get("https://r.ifyes.online:6443/public.php?op=rss&id=-1&is_cat=0&q=&key=xt9z0r6325e20d77277")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var rss Rss
	err = xml.Unmarshal(body, &rss)
	if err != nil {
		panic(err)
	}

	for _, item := range rss.Channel.Items {
		fileName := strings.ReplaceAll(item.Title, " ", "_") + ".html"
		file, err := os.Create(fileName)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		_, err = file.WriteString(item.Description)
		if err != nil {
			panic(err)
		}
	}
}
