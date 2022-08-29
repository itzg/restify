package main

import (
	"fmt"
	"github.com/itzg/restify"
	"github.com/yhat/scrape"
	"golang.org/x/net/html"
	"gopkg.in/alecthomas/kingpin.v2"
	"log"
	"os"
	"strings"
)

var (
	version = "dev"
	commit  = "master"
)

var (
	url = kingpin.Arg("url", "A URL to RESTify into JSON").
		Required().URL()
	byClass = kingpin.Flag("class", "If specified, first-level elements encountered with this class will be extracted.").
		String()
	byId = kingpin.Flag("id", "If specified, the element with this id will be extracted.").
		String()
	byTagName   = kingpin.Flag("tag", "If specified, the first-level element with this tag name will be extracted.").String()
	byAttribute = kingpin.Flag("attribute",
		"If specified, as key=value, the element with the given attribute name set to the given value is extracted.").
		String()
	showVersion = kingpin.Flag("version", "Print version and exit").
			Bool()
	debug = kingpin.Flag("debug", "Enable debugging output").
		Bool()
	userAgent = kingpin.Flag("user-agent", "user-agent header to provide with request").
			Default("restify/" + version).
			String()
	headers = kingpin.Flag("headers", "Additional headers to pass with request").
		StringMap()
)

func main() {

	kingpin.Parse()

	if *showVersion {
		log.Printf("Version: %s, Commit: %s\n", version, commit)
		os.Exit(0)
	}

	configs := make([]restify.RequestConfig, 0)
	if headers != nil {
		configs = append(configs, restify.WithHeaders(*headers))
	}
	root, err := restify.LoadContent(*url, *userAgent, configs...)
	if err != nil {
		log.Fatal("Failed to load content: ", err)
	}

	var subset []*html.Node
	if *byId != "" {
		elem, ok := restify.FindSubsetById(root, *byId)
		if !ok {
			log.Fatalf("Unable to find an element with the ID '%s'\n", *byId)
		}
		subset = append(subset, elem)
	} else if *byClass != "" {
		subset = restify.FindSubsetByClass(root, *byClass)
		if len(subset) == 0 {
			log.Fatalf("Unable to find an element with the class '%s'\n", *byClass)
		}
	} else if *byTagName != "" {
		subset = restify.FindSubsetByTagName(root, *byTagName)
		if len(subset) == 0 {
			log.Fatalf("Unable to find an element with the tag name '%s'\n", *byTagName)
		}
	} else if *byAttribute != "" {
		keyVal := strings.SplitN(*byAttribute, "=", 2)
		key := keyVal[0]
		if len(keyVal) == 1 {
			subset = restify.FindSubsetByAttributeName(root, key)
		} else {
			subset = restify.FindSubsetByAttributeNameValue(root, key, keyVal[1])
		}
		if len(subset) == 0 {
			log.Fatalf("Unable to find an element with attribute matcher %s", *byAttribute)
		}
	} else {
		subset = append(subset, root)
	}

	asJson, err := restify.ConvertHtmlToJson(subset)
	if err != nil {
		log.Fatal("Failed to parse HTML into JSON", err)
	}

	fmt.Print(string(asJson))
}

func matchByAttribute(key, value string) scrape.Matcher {
	return func(node *html.Node) bool {
		if node.Type == html.ElementNode {
			result := scrape.Attr(node, key)
			if *debug {
				fmt.Printf("Saw %s and result %s with attrs %+v\n", node.Data, result, node.Attr)
			}
			if result != "" && (value == "" || value == result) {
				return true
			}
		}
		return false
	}
}
