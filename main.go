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
	"strconv"
	"strings"
)

type Ebook struct {
	URL   string `json:"url"`
	Title string `json:"title"`
	Image string `json:"image"`
}

type Ebooks struct {
	TotalPages  int     `json:"total_pages"`
	TotalEbooks int     `json:"total_ebooks"`
	List        []Ebook `json:"ebooks"`
}

func NewEbooks() *Ebooks {
	return &Ebooks{}
}

func getPagesUrl(uri string) []string {
	var listUrl []string
	resp, err := http.Get(uri)
	if err != nil {
		return listUrl
	}
	defer resp.Body.Close()

	fmt.Println(resp.Body)

	links := collectlinks.All(resp.Body)
	for _, link := range links {
		if !strings.Contains(link, uri) {
			continue
		}
		listUrl = append(listUrl, link)
	}
	return listUrl
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
		ebooks.TotalEbooks++
		ebooks.List = append(ebooks.List, Ebook)
	})
	return nil
}

func (ebooks *Ebooks) getTotalPages(url string) error {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		return err
	}
	lastPageLink, _ := doc.Find("ul.pagination li:last-child a").Attr("href")
	split := strings.Split(lastPageLink, "?page=")
	totalPages, _ := strconv.Atoi(split[1])
	ebooks.TotalPages = totalPages
	return nil
}

func (ebooks *Ebooks) getAllEbooks(currentUrl string) error {
	eg := errgroup.Group{}
	for i := 1; i <= ebooks.TotalPages; i++ {
		uri := fmt.Sprintf("%v?page=%v", currentUrl, i)
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
	currentUrl := args[0]
	ebooks := NewEbooks()
	err := ebooks.getTotalPages(currentUrl)
	checkError(err)
	err = ebooks.getAllEbooks(currentUrl)
	checkError(err)
	ebooksJson, err := json.Marshal(ebooks)
	checkError(err)
	err = ioutil.WriteFile("output.json", ebooksJson, 0644)
	checkError(err)
}
