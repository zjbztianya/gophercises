package link

import (
	"golang.org/x/net/html"
	"io"
	"strings"
)

type Link struct {
	Href string
	Text string
}

func getNodes(root *html.Node) []*html.Node {
	if root.Type == html.ElementNode && root.Data == "a" {
		return []*html.Node{root}
	}

	var nodes []*html.Node
	for child := root.FirstChild; child != nil; child = child.NextSibling {
		nodes = append(nodes, getNodes(child)...)
	}
	return nodes
}

func getText(node *html.Node) string {
	if node.Type == html.TextNode {
		return node.Data
	}

	if node.Type != html.ElementNode {
		return ""
	}

	var str string
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		str += getText(child)
	}

	return strings.Join(strings.Fields(str), " ")
}

func buildLink(node *html.Node) (link Link) {
	for _, attr := range node.Attr {
		if attr.Key == "href" {
			link = Link{Href: attr.Val}
			link.Text = getText(node)
			return
		}
	}
	return
}

func Parse(r io.Reader) ([]Link, error) {
	root, err := html.Parse(r)
	if err != nil {
		return nil, err
	}
	nodes := getNodes(root)
	var links []Link
	for _, node := range nodes {
		links = append(links, buildLink(node))
	}
	return links, nil
}
