
A little utility (written in Go) that
enables REST-like access to HTML pages by scraping and parsing them into JSON.

```
usage: restify [<flags>] <url>

Flags:
  --help         Show context-sensitive help (also try --help-long and
                 --help-man).
  --class=CLASS  If specified, first-level elements encountered with this class
                 will be extracted.

Args:
  <url>  A URL to RESTify into JSON
```

## Output Structure

When a successful URL retrieval and match occurs, the utility will output
the JSON conversion to stdout.

The top-level structure is an array of `jsonNode`, where each `jsonNode` is
structured as:

```
{
  "name":   "...element name...",
  "class":  "...class attribute, if present...",
  "id":     "...id attribute, if present...",
  "href":   "...href attribute, if present...",
  "text":   "...element text content, if present...",
  "elements": [
    ...jsonNodes, if present...
  ]
}
```