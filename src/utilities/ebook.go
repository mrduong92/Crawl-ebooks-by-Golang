package utilities

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/sync/errgroup"
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

func (ebooks *Ebooks) GetTotalPages(url string) error {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		return err
	}
	lastPageLink, _ := doc.Find("ul.pagination li:last-child a").Attr("href")
	if lastPageLink == "javascript:void();" {
		ebooks.TotalPages = 1
		return nil
	}
	split := strings.Split(lastPageLink, "?page=")
	totalPages, _ := strconv.Atoi(split[1])
	ebooks.TotalPages = totalPages
	return nil
}

func (ebooks *Ebooks) GetAllEbooks(currentUrl string) error {
	eg := errgroup.Group{}
	if ebooks.TotalPages > 0 {
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
	}
	return nil
}
