package xp

import (
	"fmt"
	"strconv"
	"strings"
)

// Query searches for nodes matching the given path starting from root.
func Query(root Node, path string) ([]Node, error) {
	steps, err := parsePath(path)
	if err != nil {
		return nil, err
	}

	currentSet := []Node{root}

	for _, s := range steps {
		var nextSet []Node
		
		if s.Recursive {
			// Deep search (descendants)
			for _, node := range currentSet {
				matches := findAllDescendants(node, s)
				if s.Index > 0 {
					if s.Index <= len(matches) {
						matches = []Node{matches[s.Index-1]}
					} else {
						matches = []Node{}
					}
				}
				nextSet = append(nextSet, matches...)
			}
		} else {
			// Direct child
			for _, node := range currentSet {
				matches := findChildren(node, s)
				if s.Index > 0 {
					if s.Index <= len(matches) {
						matches = []Node{matches[s.Index-1]}
					} else {
						matches = []Node{}
					}
				}
				nextSet = append(nextSet, matches...)
			}
		}

		currentSet = nextSet
		if len(currentSet) == 0 {
			break
		}
	}

	return currentSet, nil
}

type step struct {
	Tag       string
	Recursive bool
	AttrKey   string
	AttrVal   string
	Index     int // 0 means no index
}

func parsePath(path string) ([]step, error) {
	parts := strings.Split(path, "/")
	var steps []step
	pendingRecursive := false

	for i, part := range parts {
		if part == "" {
			// If it's the first empty part and path starts with /, it's just the root anchor, ignore.
			// If it's a subsequent empty part (e.g. //), it means recursion.
			if i > 0 || strings.HasPrefix(path, "//") {
				// Wait, if path is "//div", parts are ["", "", "div"].
				// i=0 is "", i=1 is "".
				// If i=0 is empty, it usually just means absolute path.
				// If we see another empty, it means //
				if i > 0 && parts[i-1] == "" {
					// This is the second empty in a row, or //
					pendingRecursive = true
				} else if i > 0 {
					// /a//b -> ["", "a", "", "b"]
					// i=2 is empty. prev was "a".
					pendingRecursive = true
				} else {
					// i=0 is empty.
					// If path starts with //, then parts[1] will be empty too.
					// We can just ignore i=0 empty.
				}
			}
			continue
		}

		s := step{
			Recursive: pendingRecursive,
		}
		pendingRecursive = false

		// Parse tag and predicates
		// tag[predicate]
		bracketStart := strings.Index(part, "[")
		if bracketStart != -1 {
			if !strings.HasSuffix(part, "]") {
				return nil, fmt.Errorf("malformed predicate in %s", part)
			}
			s.Tag = part[:bracketStart]
			predicate := part[bracketStart+1 : len(part)-1]

			if strings.HasPrefix(predicate, "@") {
				// Attribute filter: @attr='val'
				eq := strings.Index(predicate, "=")
				if eq == -1 {
					return nil, fmt.Errorf("malformed attribute predicate %s", predicate)
				}
				s.AttrKey = predicate[1:eq]
				val := predicate[eq+1:]
				// Remove quotes
				if (strings.HasPrefix(val, "'") && strings.HasSuffix(val, "'")) ||
					(strings.HasPrefix(val, "\"") && strings.HasSuffix(val, "\"")) {
					s.AttrVal = val[1 : len(val)-1]
				} else {
					s.AttrVal = val
				}
			} else {
				// Try index
				idx, err := strconv.Atoi(predicate)
				if err == nil {
					s.Index = idx
				} else {
					// Unknown predicate, maybe ignore or error?
					// For "Good Enough", maybe just error
					return nil, fmt.Errorf("unsupported predicate %s", predicate)
				}
			}
		} else {
			s.Tag = part
		}

		steps = append(steps, s)
	}
	
	return steps, nil
}

func matchNode(n Node, s step) bool {
	// Check tag
	if s.Tag == "*" {
		if n.Tag() == "" {
			return false // Wildcard * matches any element, but not text nodes (which have empty tag)
		}
	} else if s.Tag != n.Tag() {
		return false
	}
	
	// Check attribute if specified
	if s.AttrKey != "" {
		if n.Attr(s.AttrKey) != s.AttrVal {
			return false
		}
	}
	
	return true
}

func findChildren(parent Node, s step) []Node {
	var results []Node
	for _, child := range parent.Children() {
		if matchNode(child, s) {
			results = append(results, child)
		}
	}
	return results
}

func findAllDescendants(root Node, s step) []Node {
	var results []Node
	
	// Depth-first traversal
	var traverse func(n Node)
	traverse = func(n Node) {
		for _, child := range n.Children() {
			if matchNode(child, s) {
				results = append(results, child)
			}
			traverse(child)
		}
	}
	
	traverse(root)
	return results
}
