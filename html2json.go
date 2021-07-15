package restify

import (
	"bytes"
	"encoding/json"
	"golang.org/x/net/html"
	"log"
	"strings"
)

// ConvertHtmlToJson the given HTML nodes into JSON content where each
// HTML node is represented by the JsonNode structure.
func ConvertHtmlToJson(nodes []*html.Node) ([]byte, error) {
	rootJsonNodes := make([]JsonNode, len(nodes))

	for i, n := range nodes {
		rootJsonNodes[i].populateFrom(n)
	}

	return json.Marshal(rootJsonNodes)
}

// JsonNode is a JSON-ready representation of an HTML node.
type JsonNode struct {
	// Name is the name/tag of the element
	Name string `json:"name,omitempty"`
	// Attributes contains the attributs of the element other than id, class, and href
	Attributes map[string]string `json:"attributes,omitempty"`
	// Class contains the class attribute of the element
	Class string `json:"class,omitempty"`
	// Id contains the id attribute of the element
	Id string `json:"id,omitempty"`
	// Href contains the href attribute of the element
	Href string `json:"href,omitempty"`
	// Text contains the inner text of the element
	Text string `json:"text,omitempty"`
	// Elements contains the child elements of the element
	Elements []JsonNode `json:"elements,omitempty"`
}

func (n *JsonNode) populateFrom(htmlNode *html.Node) *JsonNode {
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
				n.Elements = make([]JsonNode, 0)
			}
			var jsonElemNode JsonNode
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
