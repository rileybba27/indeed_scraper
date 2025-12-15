package main

import (
	"fmt"
	"strings"

	"golang.org/x/net/html"
)

type HTMLPath struct {
	Tag   string
	Index int
}

func FindNodeBasedOnPath(parent *html.Node, path []HTMLPath) (*html.Node, error) {
	var node *html.Node = nil
	var cursor *html.Node = parent
	index := 0
	pathLevel := 0
	child := cursor.FirstChild

	// log.Println("Starting request for node based on path", path)
	// log.Println("Top node is of data", parent.Data)
	// log.Println("Top child is of data", child.Data)

	for child != nil {
		if child.Type != html.ElementNode {
			child = child.NextSibling
			continue
		}

		last := child.NextSibling == nil
		pathPart := path[pathLevel]
		if (index == pathPart.Index || (last && pathPart.Index < 0)) && child.Data == pathPart.Tag {
			index = 0
			cursor = child
			child = cursor.FirstChild
			if pathLevel == len(path)-1 {
				break
			}

			pathLevel += 1
			continue
		}

		index += 1
		child = child.NextSibling
	}

	if pathLevel == len(path)-1 {
		node = cursor
	}

	if node == nil {
		return nil, fmt.Errorf("Failed to find Node based on path")
	}

	return node, nil
}

func FindNodesOfClass(parent *html.Node, class string) []*html.Node {
	var nodes = make([]*html.Node, 0)

	var crawler func(*html.Node)
	crawler = func(n *html.Node) {
		for _, attr := range n.Attr {
			if attr.Key != "class" {
				continue
			}

			val := strings.TrimSpace(attr.Val)
			if !strings.Contains(val, class) {
				continue
			}

			if strings.Contains(val, " ") && !strings.Contains(val, class+" ") && !strings.HasSuffix(val, " "+class) {
				continue
			}

			nodes = append(nodes, n)
			break
		}

		for child := n.FirstChild; child != nil; child = child.NextSibling {
			crawler(child)
		}
	}

	crawler(parent)
	return nodes
}

func FindNodeByAttr(doc *html.Node, key string, value string) (*html.Node, error) {
	var body *html.Node = nil
	var crawler func(*html.Node)
	crawler = func(node *html.Node) {
		if body != nil {
			return
		}

		for _, attr := range node.Attr {
			if attr.Key == key && attr.Val == value {
				body = node
				return
			}
		}

		for child := node.FirstChild; child != nil; child = child.NextSibling {
			crawler(child)
		}
	}

	crawler(doc)
	if body != nil {
		return body, nil
	}

	return nil, fmt.Errorf("Missing node with attribute %s of value %s", key, value)
}

func GetNodeAttr(doc *html.Node, key string) (string, error) {
	for _, attribute := range doc.Attr {
		if attribute.Key == key {
			return attribute.Val, nil
		}
	}

	return "", fmt.Errorf("Couldn't find node attribute %s", key)
}

func CollectNodeText(doc *html.Node) string {
	var text string
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.TextNode {
			text += n.Data + " "
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	text = strings.TrimSpace(text)
	return text
}
