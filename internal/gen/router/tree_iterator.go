package router

func treeFindCreate(node *node, path []string) *node {
	for _, v := range path {
		next, ok := node.children[v]
		if !ok {
			next = newNode(node)
			node.children[v] = next
		}

		node = next
	}

	return node
}
