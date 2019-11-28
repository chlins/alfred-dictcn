package main

import (
	"net/http"

	"github.com/PuerkitoBio/goquery"
	aw "github.com/deanishe/awgo"
	"golang.org/x/net/html"
)

var (
	wf *aw.Workflow
)

func init() {
	wf = aw.New()
}

func run() {
	var query string
	if args := wf.Args(); len(args) > 0 {
		query = args[0]
	}

	doRequest(query)
	wf.SendFeedback()
}

func doRequest(word string) {
	url := "http://apii.dict.cn/mini.php?q=" + word
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		wf.FatalError(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		wf.FatalError(err)
	}

	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		wf.FatalError(err)
	}

	var res []string
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.TextNode {
			res = append(res, n.Data)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}

	doc.Find("div#e").Each(func(i int, selection *goquery.Selection) {
		for _, node := range selection.Nodes {
			f(node)
		}
	})

	for _, e := range res {
		wf.NewItem(e)
	}
}

func main() {
	wf.Run(run)
}
