package html2md

import (
	"golang.org/x/net/html"
	"io"
	log "github.com/Sirupsen/logrus"
	"fmt"
	"strings"
	"bytes"
)

func getAttributeValue(attrs []html.Attribute, name string) (value string) {
	for _, attr := range attrs {
		if attr.Key == name {
			value = attr.Val
		}
	}
	return
}

func generate(initialNode *html.Node) (value string) {
	for node := initialNode; node != nil; node = node.NextSibling {
		switch node.Type {
		case html.TextNode:
			value += node.Data + generate(node.FirstChild)
		case html.ElementNode:
			tag := node.Data
			emchar := "*"
			if len(value) > 0 && value[len(value) - 1] == '*' {
				emchar = "_"
			}
			switch tag {
			case "b", "strong": {
				value += fmt.Sprintf("%s%s%s%s%s", emchar, emchar, generate(node.FirstChild), emchar, emchar)
			}
			case "i", "em": {
				value += fmt.Sprintf("%s%s%s", emchar, generate(node.FirstChild), emchar)
			}
			case "br": {
				value += "\n"
			}
			case "code": {
				value += fmt.Sprintf("`%s`", generate(node.FirstChild))
			}
			case "pre": {
				value += fmt.Sprintf("```%s```", generate(node.FirstChild))
			}
			case "hr": {
				value += "\n* * *\n"
			}
			case "h1", "h2", "h3", "h4", "h5", "h6":
				level := int(tag[1] - '0')
				value += fmt.Sprintf("\n%s %s\n", strings.Repeat("#", level), generate(node.FirstChild))
			case "a": {
				href := getAttributeValue(node.Attr, "href")
				value += fmt.Sprintf("[%s](%s)", generate(node.FirstChild), href)
			}
			case "ol", "ul": {
				value += "\n" + generate(node.FirstChild)
			}
			case "li": {
				if node.Parent != nil && node.Parent.Type == html.ElementNode {
					pTag := node.Parent.Data
					if pTag == "ul" {
						value += fmt.Sprint("* ", generate(node.FirstChild), "\n")
					} else if pTag == "ol" {
						value += fmt.Sprint("1. ", generate(node.FirstChild), "\n")
					} else {
						continue
					}
				}
			}
			default:
				// less breakage than simply omitting node
				value += generate(node.FirstChild)
			}
		default:
			log.Warn("Unsupported node type")
			value += generate(node.FirstChild)
		}
	}
	return
}

func HTML2MD(r io.Reader) (*string, error) {
	nodes, err := html.ParseFragment(r, &html.Node{Type: html.ElementNode})
	if err != nil {
		return nil, err
	}


	var v string
	for _, node := range nodes {
		v += generate(node)
	}

	return &v, nil
}

func main() {
	r := bytes.NewReader([]byte(`<b>test</b>
<h3>foo</h3>
<pre>
	abc
	def
</pre>
<a href="test"><strong>bcd</strong><em>rad</em></a>
<strong>c</strong>`))
	v, _:= HTML2MD(r)
	print(*v)
}