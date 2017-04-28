package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/headzoo/surf/agent"
	"github.com/urfave/cli"
	"gopkg.in/headzoo/surf.v1"
)

func main() {
	app := cli.NewApp()
	app.Name = "gocodipodo - The comdirect Postbox downloader"
	app.Version = "0.1"
	app.Compiled = time.Now()
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "Kevin Lindecke",
			Email: "kevin@lindecke.co",
		},
	}

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "user, u",
			Usage: "comdirect `userid`",
		},
		cli.StringFlag{
			Name:  "pass, p",
			Usage: "comdirect `password`",
		},
	}

	app.Action = func(c *cli.Context) error {
		bow := surf.NewBrowser()
		bow.SetUserAgent(agent.Chrome())
		err := bow.Open("https://kunde.comdirect.de/lp/wt/login")
		if err != nil {
			panic(err)
		}

		fm, err := bow.Form("#login")
		if err != nil {
			panic(err)
		}

		fm.Set("loginAction", "loginAction")
		fm.Input("param1", c.String("user"))
		fm.Input("param3", c.String("pass"))

		err = fm.Submit()
		if err != nil {
			panic(err)
		}

		bow.Find(".error-message__text").Each(func(_ int, s *goquery.Selection) {
			fmt.Println(s.Text())
		})

		err = bow.Open("https://kunde.comdirect.de/itx/posteingangsuche")
		if err != nil {
			panic(err)
		}

		var links []string
		for _, link := range bow.Links() {
			if strings.Contains(link.URL.Path, "dokumentenabruf") {
				links = append(links, link.URL.String())
			}
		}

		for _, link := range links {
			filename := path.Base(link)
			fout, err := os.Create(filename)
			if err != nil {
				log.Printf(
					"Error creating file '%s'.", filename)
				continue
			}
			defer fout.Close()

			bow.Open(link)
			_, err = bow.Download(fout)
			if err != nil {
				log.Printf(
					"Error downloading file '%s'.", filename)
			} else {
				log.Printf(
					"Downloaded '%s'.", filename)
			}
		}

		return nil
	}

	app.Run(os.Args)
}
