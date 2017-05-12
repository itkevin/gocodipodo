package main

import (
	"log"
	"os"
	"path"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/urfave/cli"
	"gopkg.in/headzoo/surf.v1"
)

const (
	LoginURL   = "https://kunde.comdirect.de/lp/wt/login"
	PostboxURL = "https://kunde.comdirect.de/itx/posteingangsuche"
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
		err := bow.Open(LoginURL)
		checkErr(err)

		fm, err := bow.Form("#login")
		checkErr(err)

		fm.Set("loginAction", "loginAction")
		fm.Input("param1", c.String("user"))
		fm.Input("param3", c.String("pass"))

		err = fm.Submit()
		checkErr(err)

		bow.Find(".error-message__text").Each(func(_ int, s *goquery.Selection) {
			log.Fatal("ERROR: ", s.Text())
		})

		err = bow.Open(PostboxURL)
		checkErr(err)

		fm1, err := bow.Form("#f1")
		checkErr(err)

		fm1.Set("f1-j_idt108-nextLink", "f1-j_idt108-nextLink")
		fm1.Set("f1", "f1")
		fm1.Input("f1-monatInput", c.String(""))
		fm1.Input("f1-jahrInput1", c.String(""))
		fm1.Input("f1-wkn1", c.String(""))
		fm1.Input("f1-docTyp1", c.String("0"))

		err = fm1.Submit()
		checkErr(err)

		//log.Printf("%s", bow)

		var links []string
		for _, link := range bow.Links() {
			if strings.Contains(link.URL.Path, "dokumentenabruf") {
				links = append(links, link.URL.String())
			}
		}

		subdir := "comdirect"
		os.MkdirAll(subdir, os.ModePerm)
		for _, link := range links {
			filename := subdir + "/" + path.Base(link)
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

func checkErr(err error) {
	if err != nil {
		log.Fatal("ERROR:", err)
	}
}
