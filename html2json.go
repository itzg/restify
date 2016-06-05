package main

import (
	"bytes"
	"golang.org/x/net/html"
	"encoding/json"
	"strings"
	"log"
)

func convertToJson(node *html.Node) ([]byte, error) {
	var rootJsonNode jsonNode

	rootJsonNode.populateFrom(node)

	return json.Marshal(rootJsonNode)
}

type jsonNode struct {
	Attributes map[string]string `json:"attributes,omitempty"`
	Class string `json:"class,omitempty"`
	Id string `json:"id,omitempty"`
	Href string `json:"href,omitempty"`
	Text       string `json:"text,omitempty"`
	Elements   map[string]*jsonNode `json:"elements,omitempty"`
}

func (n *jsonNode) populateFrom(htmlNode *html.Node) {
	switch htmlNode.Type {
	case html.ElementNode, html.DocumentNode:
		break

	default:
		log.Fatal("Given node needs to be an element or document")
	}

	var textBuffer bytes.Buffer

	if len(htmlNode.Attr) > 0 {
		n.Attributes = make(map[string]string)
		var a html.Attribute
		for _, a = range htmlNode.Attr {
			switch a.Key {
			case "class":
				n.Class = a.Val

			case "id":
				n.Id = a.Val

			case "href":
				n.Href = a.Val

			default:
				n.Attributes[a.Key] = a.Val
			}
		}
	}

	e := htmlNode.FirstChild
	for e != nil {
		switch e.Type {
		case html.TextNode:
			trimmed := strings.TrimSpace(e.Data)
			if len(trimmed) > 0 {
				// mimic HTML text normalizing
				if textBuffer.Len() > 0 {
					textBuffer.WriteString(" ")
				}
				textBuffer.WriteString(trimmed)
			}

		case html.ElementNode:
			if n.Elements == nil {
				n.Elements = make(map[string]*jsonNode)
			}
			n.Elements[e.Data] = new(jsonNode)
			n.Elements[e.Data].populateFrom(e)
		}

		e = e.NextSibling
	}

	if textBuffer.Len() > 0 {
		n.Text = textBuffer.String()
	}
}
