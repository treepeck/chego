package chego

/*
TraversePreOrder traverses the tree in pre-order, starting from the specified node.
*/
func TraversePreOrder(n *Node, codes *[218]string, current string) {
	if n == nil {
		return
	}

	if n.Left == nil && n.Right == nil {
		(*codes)[n.Index] = current
		return
	}

	TraversePreOrder(n.Left, codes, current+"1")
	TraversePreOrder(n.Right, codes, current+"0")
}
