package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Feed struct {
	XMLName xml.Name `xml:"feed"`
	Entries []Entry  `xml:"entry"`
}

type Entry struct {
	Title     string    `xml:"title"`
	Link      Link      `xml:"link"`
	Published time.Time `xml:"published"`
	Content   Content   `xml:"content"`
}

type Link struct {
	Href string `xml:"href,attr"`
}

type Content struct {
	Type  string `xml:"type,attr"`
	Value string `xml:",chardata"`
}

func main() {
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

	var feed Feed
	err = xml.Unmarshal(body, &feed)
	if err != nil {
		log.Fatal("Error unmarshalling XML:", err)
	}

	for _, entry := range feed.Entries {
		title := strings.ReplaceAll(entry.Title, " ", "_")
		filename := fmt.Sprintf("%s.html", title)

		f, err := os.Create(filename)
		if err != nil {
			log.Fatal("Error creating file:", err)
		}
		defer f.Close()

		_, err = f.WriteString(entry.Content.Value)
		if err != nil {
			log.Fatal("Error writing to file:", err)
		}

		err = os.Rename(filename, filepath.Join("rss-stared", filename))
		if err != nil {
			log.Fatal("Error moving file to rss-stared directory:", err)
		}

		fmt.Printf("Created file: %s\n", filename)
	}
}
