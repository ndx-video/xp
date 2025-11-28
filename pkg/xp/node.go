package xp

// Node represents a node in a tree structure that can be queried.
type Node interface {
	// Tag returns the element name (e.g., "div").
	Tag() string

	// Attr returns the value of the specified attribute, or an empty string if not present.
	Attr(key string) string

	// Children returns the child nodes of the current node.
	Children() []Node

	// Text returns the text content of the node (if applicable).
	Text() string
}
