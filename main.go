package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/jackdanger/collectlinks"
	"golang.org/x/sync/errgroup"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

type Ebook struct {
	URL   string `json:"url"`
	Title string `json:"title"`
	Image string `json:"image"`
}

type Ebooks struct {
	List []Ebook `json:"ebooks"`
}

func getLinks(uri string) []string {
	var url []string
	resp, err := http.Get(uri)
	if err != nil {
		return url
	}
	defer resp.Body.Close()

	links := collectlinks.All(resp.Body)
	for _, link := range links {
		if !strings.Contains(link, uri) {
			continue
		}
		url = append(url, link)
	}
	return url
}

func (ebooks *Ebooks) getEbooksByUrl(url string) error {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		return err
	}

	doc.Find(".col-left ._2pin").Each(func(i int, s *goquery.Selection) {
		docTitle, exists := s.Find(".ellipsis a").Attr("title")
		if !exists {
			docTitle = ""
		}
		docLink, exists := s.Find(".ellipsis a").Attr("href")
		if !exists {
			docLink = "#"
		}
		docImg, exists := s.Find("a._3if7 img").Attr("src")
		if !exists {
			docImg = ""
		}
		Ebook := Ebook{
			URL:   docLink,
			Title: docTitle,
			Image: docImg,
		}
		ebooks.List = append(ebooks.List, Ebook)
	})
	return nil
}

func (ebooks *Ebooks) getAllEbooks(listURL []string) error {
	eg := errgroup.Group{}
	for _, url := range listURL {
		uri := url
		eg.Go(func() error {
			err := ebooks.getEbooksByUrl(uri)
			if err != nil {
				return err
			}
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		return err
	}
	return nil
}

func checkError(err error) {
	if err != nil {
		log.Println(err)
	}
}

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("Please specify start page")
		os.Exit(1)
	}
	links := getLinks(args[0])
	ebooks := Ebooks{}
	err := ebooks.getAllEbooks(links)
	checkError(err)
	ebooksJson, err := json.Marshal(ebooks)
	checkError(err)
	err = ioutil.WriteFile("output.json", ebooksJson, 0644)
	checkError(err)
}
