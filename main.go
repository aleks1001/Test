package main

import (
	"net/http"
	"net/url"
	"errors"
	"log"
	"fmt"
	"golang.org/x/net/html"
)

type Document struct {
	*Selection
	Url       *url.URL
	rootNode  *html.Node
	tokenizer *html.Tokenizer
}

type Selection struct {
	Nodes    []*html.Node
	document *Document
	prevSel  *Selection
}

func main() {
//	http.HandleFunc("/test", someFunc)
//	http.ListenAndServe(":8080", nil)

	doc, err := NewDocument("http://willfe.com/2013/10/2/great-x3-albion-prelude-add-list/")
	if err != nil {
		log.Fatal(err)
	}

	doc.ParseToken()

	parseNode(doc.rootNode)
}

/*

func someFunc(w http.ResponseWriter, r *http.Request) {


	res, e := http.Get("http://golang-examples.tumblr.com/post/47426518779/parse-html")
	if e != nil {
		log.Fatalln(e)
	}

	tkz := html.NewTokenizer(res.Body)

	for {
		tokenType := tkz.Next()
		if tokenType == html.ErrorToken {
			w.Write([]byte(tokenType.String()))
			log.Fatalln(tokenType.String())
		}
		token := tkz.Token()
		switch tokenType {
		case html.StartTagToken:
			for _, a := range token.Attr {
				if a.Key == "href" && token.Data != "" {
					log.Println(a.Val)
					break
				}
			}
		case html.TextToken:
			for _, a := range token.Attr {
				log.Println("Text Tocken: " + a.Key + " " + a.Val)
			}
		}
	}
	//	var buf bytes.Buffer
	//	buf.WriteString(doc.tkz.Text())
	//	tt := doc.tkz.Next()
	//	w.Write([]byte(tt.String()))
}
*/

func (d *Document) ParseToken() {
	fmt.Println(d)
	for {
		tokenType := d.tokenizer.Next()
		if tokenType == html.ErrorToken {
			log.Fatalln("Something wrong!")
		}
		token := d.tokenizer.Token()
		switch tokenType {
		case html.StartTagToken:
			for _, a := range token.Attr {
				log.Println(a.Key + " " + a.Val + " " + a.Namespace)
				if a.Key == "trackindex" {
					log.Println(a.Val)
					break
				}
			}
		case html.TextToken:
			log.Println(token.Data)
		case html.EndTagToken:
			log.Println(token.Data)
		case html.SelfClosingTagToken:
			log.Println(token.Data)
		}
	}
}

func parseNode(n *html.Node) {
	if n.Type == html.ElementNode && n.Data == "td" {
		for _, a := range n.Attr {
			log.Println(a.Key + " " + a.Val + " " + a.Namespace)
			if a.Key == "trackindex" {
				//log.Println(a.Val)
				break
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		parseNode(c)
	}
}

func NewDocument(url string) (*Document, error) {
	res, e := http.Get(url)
	if e != nil {
		return nil, e
	}
	return NewDocumentFromResponse(res)
}

func NewDocumentFromResponse(res *http.Response) (*Document, error) {
	if res == nil {
		return nil, errors.New("Response is nil pointer")
	}
	// Parse the HTML into nodes
	root, e := html.Parse(res.Body)
	tkz := html.NewTokenizer(res.Body)
	if e != nil {
		return nil, e
	}
	defer res.Body.Close()
	return newDocument(root, tkz, res.Request.URL), nil
}

func newDocument(root *html.Node, tokenizer *html.Tokenizer, url *url.URL) *Document {
	d := &Document{nil, url, root, tokenizer}
	d.Selection = newSingleSelection(root, d)
	return d
}

func newSingleSelection(node *html.Node, doc *Document) *Selection {
	return &Selection{[]*html.Node{node}, doc, nil}
}
