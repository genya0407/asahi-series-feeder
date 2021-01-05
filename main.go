package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/labstack/echo/v4"
)

func handleRSS(c echo.Context) error {
	id := c.Param("id")
	crawler := AsahiSeriesCrawler{SeriesID: id}
	feed, err := crawler.Crawl()
	if err != nil {
		panic(err)
	}
	r, err := feed.ToRSSReader()
	if err != nil {
		panic(err)
	}
	s, err := ioutil.ReadAll(r)
	if err != nil {
		panic(err)
	}

	return c.String(http.StatusOK, string(s))
}

func handleAtom(c echo.Context) error {
	id := c.Param("id")
	crawler := AsahiSeriesCrawler{SeriesID: id}
	feed, err := crawler.Crawl()
	if err != nil {
		panic(err)
	}
	r, err := feed.ToAtomReader()
	if err != nil {
		panic(err)
	}
	s, err := ioutil.ReadAll(r)
	if err != nil {
		panic(err)
	}

	return c.String(http.StatusOK, string(s))
}

func main() {
	port := flag.String("port", "8080", "listening port number")
	flag.Parse()

	e := echo.New()
	e.GET("/series/rss/:id", handleRSS)
	e.GET("/series/atom/:id", handleAtom)
	e.Logger.Fatal(e.Start(fmt.Sprintf("localhost:%s", *port)))
}
