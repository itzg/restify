
A little utility/library (written in Go) that enables REST-like access to HTML pages by scraping and parsing them into JSON.

[![CircleCI](https://circleci.com/gh/itzg/restify/tree/master.svg?style=svg)](https://circleci.com/gh/itzg/restify/tree/master)

```
usage: restify [<flags>] <url>

Flags:
  --help                      Show context-sensitive help (also try --help-long and --help-man).
  --class=CLASS               If specified, first-level elements encountered with this class will be extracted.
  --id=ID                     If specified, the element with this id will be extracted.
  --tag=TAGNAME               If specified, the first-level element with this tag name will be extracted.
  --attribute=ATTRIBUTE       If specified, as key=value, the element with the given attribute name set to the given value is extracted.
  --version                   Print version and exit
  --debug                     Enable debugging output
  --user-agent="restify/1.4.0"  user-agent header to provide with request

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

## Examples

Locate the latest Minecraft Bedrock server version by picking off the `<a>`'s with `data-platform` set:

```bash
restify --attribute=data-platform https://www.minecraft.net/en-us/download/server/bedrock/
```

which produces:
```json
[
  {"name":"a","attributes":{"data-platform":"serverBedrockWindows","role":"button"},"class":"btn btn-disabled-outline mt-4 downloadlink","href":"https://minecraft.azureedge.net/bin-win/bedrock-server-1.12.0.28.zip","text":"Download"},
  {"name":"a","attributes":{"data-platform":"serverBedrockLinux","role":"button"},"class":"btn btn-disabled-outline mt-4 downloadlink","href":"https://minecraft.azureedge.net/bin-linux/bedrock-server-1.12.0.28.zip","text":"Download"}
]
```

or to grab just the Linux instance:

```bash
restify --attribute=data-platform=serverBedrockLinux https://www.minecraft.net/en-us/download/server/bedrock/
```

## Using as a library

The package `github.com/itzg/restify` provides the library functions used by the command-line utility.
