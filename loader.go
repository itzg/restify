package restify

import (
	"fmt"
	"github.com/yhat/scrape"
	"golang.org/x/net/html"
	"net/http"
	"net/url"
)

type RequestConfig func(*http.Request)

// WithHeaders configures additional headers in the request used in LoadContent
func WithHeaders(headers map[string]string) RequestConfig {
	return func(request *http.Request) {
		for k, v := range headers {
			request.Header.Set(k, v)
		}
	}
}

// LoadContent retrieves the HTML content from the given url.
// The userAgent is optional, but if provided should conform with https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/User-Agent
func LoadContent(url *url.URL, userAgent string, configs ...RequestConfig) (*html.Node, error) {
	request, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("Failed to request: %w", err)
	}

	request.Header.Set("accept", "*/*")
	if userAgent != "" {
		request.Header.Set("user-agent", userAgent)
	}
	for _, config := range configs {
		config(request)
	}

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("Failed to retrieve response: %w", err)
	}
	defer resp.Body.Close()

	root, err := html.Parse(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse response body: %w", err)
	}

	return root, nil
}

// FindSubsetById locates the HTML node within the given root that has an id attribute of given value.
// If the node is not found, then ok will be false.
func FindSubsetById(root *html.Node, id string) (n *html.Node, ok bool) {
	return scrape.Find(root, scrape.ById(id))
}

// FindSubsetByClass locates the HTML nodes with the given root that have the given className.
func FindSubsetByClass(root *html.Node, className string) []*html.Node {
	return scrape.FindAll(root, scrape.ByClass(className))
}

// FindSubsetByAttributeName retrieves the HTML nodes that have the requested
// attribute, regardless of their values.
func FindSubsetByAttributeName(root *html.Node, attribute string) []*html.Node {
	return FindSubsetByAttributeNameValue(root, attribute, "")
}

// FindSubsetByAttributeNameValue retrieves the HTML nodes that have the requested attribute with a specific value.
func FindSubsetByAttributeNameValue(root *html.Node, attribute string, value string) []*html.Node {
	return scrape.FindAll(root, matchByAttribute(attribute, value))
}

func matchByAttribute(key, value string) scrape.Matcher {
	return func(node *html.Node) bool {
		if node.Type == html.ElementNode {
			result := scrape.Attr(node, key)
			if result != "" && (value == "" || value == result) {
				return true
			}
		}
		return false
	}
}
