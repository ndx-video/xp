package xp

import (
	"reflect"
	"testing"
)

type mockNode struct {
	tag      string
	attrs    map[string]string
	children []Node
	text     string
}

func (m *mockNode) Tag() string {
	return m.tag
}

func (m *mockNode) Attr(key string) string {
	if m.attrs == nil {
		return ""
	}
	return m.attrs[key]
}

func (m *mockNode) Children() []Node {
	return m.children
}

func (m *mockNode) Text() string {
	return m.text
}



func extractTags(nodes []Node) []string {
	var res []string
	for _, n := range nodes {
		res = append(res, n.Tag())
	}
	return res
}

func extractTexts(nodes []Node) []string {
	var res []string
	for _, n := range nodes {
		if n.Text() != "" {
			res = append(res, n.Text())
		} else {
			// Try children text?
			// For this test, leaf nodes have text.
			// If node is div, it has no text.
			res = append(res, "")
		}
	}
	return res
}

func extractAttrs(nodes []Node, key string) []string {
	var res []string
	for _, n := range nodes {
		res = append(res, n.Attr(key))
	}
	return res
}

func TestQueryFeatures(t *testing.T) {
	// Re-create structure for easy reference
	// root
	//   div (id=1)
	//     span (class=a) -> "Text A"
	//     span (class=b) -> "Text B"
	//   div (id=2)
	//     p
	//       span (class=a) -> "Text C"
	//   ul (id=u1)
	//     li (id=l1) -> "1"
	//     li (id=l2) -> "2"
	//   ul (id=u2)
	//     li (id=l3) -> "3"
	//     li (id=l4) -> "4"

	root := &mockNode{tag: "root", children: []Node{
		&mockNode{tag: "div", attrs: map[string]string{"id": "1"}, children: []Node{
			&mockNode{tag: "span", attrs: map[string]string{"class": "a"}, text: "Text A"},
			&mockNode{tag: "span", attrs: map[string]string{"class": "b"}, text: "Text B"},
		}},
		&mockNode{tag: "div", attrs: map[string]string{"id": "2"}, children: []Node{
			&mockNode{tag: "p", children: []Node{
				&mockNode{tag: "span", attrs: map[string]string{"class": "a"}, text: "Text C"},
			}},
		}},
		&mockNode{tag: "ul", attrs: map[string]string{"id": "u1"}, children: []Node{
			&mockNode{tag: "li", attrs: map[string]string{"id": "l1"}, text: "1"},
			&mockNode{tag: "li", attrs: map[string]string{"id": "l2"}, text: "2"},
		}},
		&mockNode{tag: "ul", attrs: map[string]string{"id": "u2"}, children: []Node{
			&mockNode{tag: "li", attrs: map[string]string{"id": "l3"}, text: "3"},
			&mockNode{tag: "li", attrs: map[string]string{"id": "l4"}, text: "4"},
		}},
	}}

	// 1. /tag: Direct Child
	t.Run("Direct Child", func(t *testing.T) {
		nodes, err := Query(root, "/div")
		if err != nil {
			t.Fatalf("Query failed: %v", err)
		}
		if len(nodes) != 2 {
			t.Errorf("Expected 2 divs, got %d", len(nodes))
		}
		ids := extractAttrs(nodes, "id")
		if ids[0] != "1" || ids[1] != "2" {
			t.Errorf("Expected ids 1, 2. Got %v", ids)
		}
	})

	// 2. //tag: Descendant (Deep)
	t.Run("Descendant", func(t *testing.T) {
		nodes, err := Query(root, "//span")
		if err != nil {
			t.Fatalf("Query failed: %v", err)
		}
		if len(nodes) != 3 {
			t.Errorf("Expected 3 spans, got %d", len(nodes))
		}
		texts := extractTexts(nodes)
		// Order depends on traversal. DFS: div1->spanA, div1->spanB, div2->p->spanC
		expected := []string{"Text A", "Text B", "Text C"}
		if !reflect.DeepEqual(texts, expected) {
			t.Errorf("Expected %v, got %v", expected, texts)
		}
	})

	// 3. *: Wildcard Tag
	t.Run("Wildcard", func(t *testing.T) {
		// /div/* -> should match spans of div1 and p of div2
		nodes, err := Query(root, "/div/*")
		if err != nil {
			t.Fatalf("Query failed: %v", err)
		}
		// div1 has 2 spans. div2 has 1 p. Total 3.
		if len(nodes) != 3 {
			t.Errorf("Expected 3 nodes, got %d", len(nodes))
		}
		tags := extractTags(nodes)
		expected := []string{"span", "span", "p"}
		if !reflect.DeepEqual(tags, expected) {
			t.Errorf("Expected %v, got %v", expected, tags)
		}
	})

	// 4. [@k='v']: Attribute Match
	t.Run("Attribute Match", func(t *testing.T) {
		nodes, err := Query(root, "//span[@class='a']")
		if err != nil {
			t.Fatalf("Query failed: %v", err)
		}
		if len(nodes) != 2 {
			t.Errorf("Expected 2 spans, got %d", len(nodes))
		}
		texts := extractTexts(nodes)
		expected := []string{"Text A", "Text C"}
		if !reflect.DeepEqual(texts, expected) {
			t.Errorf("Expected %v, got %v", expected, texts)
		}
	})

	// 5. [n]: Index (1-based)
	t.Run("Index", func(t *testing.T) {
		// //ul/li[1]
		// Should return first li of EACH ul.
		// ul1 -> li1
		// ul2 -> li3
		nodes, err := Query(root, "//ul/li[1]")
		if err != nil {
			t.Fatalf("Query failed: %v", err)
		}
		// Current implementation might return only li1 if index is applied globally.
		// We want to verify if it supports per-parent index or not.
		// If it returns 1 node, it's global index. If 2, it's per-parent.
		// The prompt asked to "make sure these features are implemented".
		// Usually [1] implies per-parent in XPath.
		
		ids := extractAttrs(nodes, "id")
		t.Logf("Got ids: %v", ids)
		
		// I will assert what I think is correct (per-parent), and if it fails, I fix the code.
		if len(nodes) != 2 {
			t.Errorf("Expected 2 nodes (one per ul), got %d", len(nodes))
		}
	})
}
