package router

import "sort"

type node struct {
	parent   *node
	children map[string]*node

	handler  *handler
	endpoint *endpoint
}

func newNode(parent *node) *node {
	return &node{
		parent:   parent,
		children: map[string]*node{},
	}
}

type nodeChild struct {
	name string
	next *node
}

func (n *node) sortedChildren() []nodeChild {
	var str []string
	for name := range n.children {
		str = append(str, name)
	}
	sort.Strings(str)

	var children []nodeChild
	for _, s := range str {
		children = append(children, nodeChild{
			name: s,
			next: n.children[s],
		})
	}
	return children
}
