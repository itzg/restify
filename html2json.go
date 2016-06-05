package main

import (
	"bytes"
	"golang.org/x/net/html"
	"encoding/json"
	"strings"
	"log"
)

func convertToJson(nodes []*html.Node) ([]byte, error) {
	rootJsonNodes := make([]jsonNode, len(nodes))

	for i,n := range nodes {
		rootJsonNodes[i].populateFrom(n)
	}

	return json.Marshal(rootJsonNodes)
}

type jsonNode struct {
	Name string `json:"name,omitempty"`
	Attributes map[string]string `json:"attributes,omitempty"`
	Class string `json:"class,omitempty"`
	Id string `json:"id,omitempty"`
	Href string `json:"href,omitempty"`
	Text       string `json:"text,omitempty"`
	Elements   []jsonNode `json:"elements,omitempty"`
}

func (n *jsonNode) populateFrom(htmlNode *html.Node) *jsonNode {
	switch htmlNode.Type {
	case html.ElementNode:
		n.Name = htmlNode.Data
		break

	case html.DocumentNode:
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
				n.Elements = make([]jsonNode,0)
			}
			var jsonElemNode jsonNode
			jsonElemNode.populateFrom(e)
			n.Elements = append(n.Elements, jsonElemNode)
		}

		e = e.NextSibling
	}

	if textBuffer.Len() > 0 {
		n.Text = textBuffer.String()
	}

	return n
}
