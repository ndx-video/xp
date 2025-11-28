package main

import (
	"fmt"
	"os"

	"golang.org/x/net/html"
	"xp/pkg/xp"
)

type htmlNode struct {
	n *html.Node
}

func (h *htmlNode) Tag() string {
	if h.n.Type == html.ElementNode {
		return h.n.Data
	}
	return ""
}

func (h *htmlNode) Attr(key string) string {
	for _, a := range h.n.Attr {
		if a.Key == key {
			return a.Val
		}
	}
	return ""
}

func (h *htmlNode) Children() []xp.Node {
	var children []xp.Node
	for c := h.n.FirstChild; c != nil; c = c.NextSibling {
		// Include ElementNode and TextNode. 
		// CommentNode, DoctypeNode etc are ignored for now as per "Good Enough"
		if c.Type == html.ElementNode || c.Type == html.TextNode {
			children = append(children, &htmlNode{n: c})
		}
	}
	return children
}

func (h *htmlNode) Text() string {
	if h.n.Type == html.TextNode {
		return h.n.Data
	}
	return ""
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: xp <xpath>")
		os.Exit(1)
	}

	if os.Args[1] == "--version" || os.Args[1] == "-v" {
		fmt.Printf("xp version %s\n", xp.Version)
		os.Exit(0)
	}

	query := os.Args[1]

	doc, err := html.Parse(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing HTML: %v\n", err)
		os.Exit(1)
	}

	root := &htmlNode{n: doc}
	
	nodes, err := xp.Query(root, query)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error executing query: %v\n", err)
		os.Exit(1)
	}

	for _, node := range nodes {
		hn, ok := node.(*htmlNode)
		if ok {
			// If it's a text node, just print the text?
			// Or render it (which prints text).
			html.Render(os.Stdout, hn.n)
			// fmt.Println() // Render doesn't add newline, but maybe we want one?
			// If we print multiple elements, they might run together.
			// Let's add a newline for clarity, or maybe not if we want exact output.
			// The prompt says "Prints the HTML of the matching nodes".
			// If I match <li>a</li><li>b</li>, printing them concatenated is <li>a</li><li>b</li>.
			// Adding newline makes it <li>a</li>\n<li>b</li>.
			// I'll add a newline.
			fmt.Println()
		}
	}
}
