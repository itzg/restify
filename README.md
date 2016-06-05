
A little utility (written in Go) that
enables REST-like access to HTML pages by scraping and parsing them into JSON.

```
usage: restify [<flags>] <url>

Flags:
  --help         Show context-sensitive help (also try --help-long and
                 --help-man).
  --class=CLASS  If specified, the first element encountered with this class
                 will be extracted.

Args:
  <url>  A URL to RESTify into JSON
```