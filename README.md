# xp

> **The "Good Enough" XML/HTML Path Query Engine for Go.**
> *Lightweight. Zero-Dependency Core. Pragmatic.*

`xp` is a path query library designed for scenarios where full XPath 1.0 support is overkill, but you still need to find specific nodes in a DOM tree efficiently.

It was built to power the component resolution system of the **AsciiDoc CMS** and to provide a query interface for **asciidoc-xml**.

## 🚀 Features

* **Zero Dependencies:** The core library (`pkg/xp`) uses only the Go standard library.
* **Interface Driven:** It works on *any* node structure (standard `html.Node`, custom DOMs, JSON trees) as long as they implement the `xp.Node` interface.
* **"Good Enough" Syntax:** Supports the subset of XPath you actually use 99% of the time.

## 📦 Installation

```bash
go get [github.com/ndx-video/xp](https://github.com/ndx-video/xp)

## Usage

package main

import (
    "fmt"
    "[github.com/ndx-video/xp](https://github.com/ndx-video/xp)"
)

func main() {
    // Assuming 'root' is your DOM node
    // Find all 'div' elements with class 'hero'
    nodes, _ := xp.Query(root, "//div[@class='hero']")

    for _, n := range nodes {
        fmt.Println("Found:", n.Tag())
    }
}
```

## Supported Syntax
```
Syntax,Description,Example
/tag,Direct Child,/html/body
//tag,Descendant (Deep),//div
*,Wildcard Tag,/div/*
[@k='v'],Attribute Match,//section[@id='main']
[n],Index (1-based),//ul/li[1]
```

## Install CLI tool
```bash
go install [github.com/ndx-video/xp/cmd/xp@latest](https://github.com/ndx-video/xp/cmd/xp@latest)
```

## CLI usage
```bash
# Extract the title from a webpage
curl -s [https://example.com](https://example.com) | xp "//title"

# Find specific content
cat index.html | xp "//div[@class='content']/p[1]"
```

## 🏗 Integration with AsciiDoc CMS
xp serves as the query engine for the CMS templating system. It allows theme developers to inject content using simple path expressions:

```html
<!-- Inside a CMS Layout Template -->
<div class="hero-text">
    {{ .Doc.Query "//cms-component[@name='hero']" }}
</div>
```

## 🤝 Contributing
We love contributions! If you want to add support for [contains()] or other operators, feel free to fork and PR.

Note: We use vi for editing.

## 📄 License
MIT License.