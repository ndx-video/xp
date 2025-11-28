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
