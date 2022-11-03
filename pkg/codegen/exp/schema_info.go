package exp

import (
	"strings"
)

type SchemaInfo struct {
}

type PathTreeNodesByPath map[string]*PathTreeNode

func GenerateSchemaInfosFromTree(root *PathTreeNode) PathTreeNodesByPath {
	nodesByPath := make(PathTreeNodesByPath)
	traverseCollectNodes(nil, root, nodesByPath)
	return nodesByPath
}

func traverseCollectNodes(path []string, node *PathTreeNode, nodesByPath PathTreeNodesByPath) {
	if node.Ref != "" {
		return
	}
	if node.Schema != nil {
		nodesByPath[strings.Join(path, "/")] = node
	}
	for _, childName := range SortedMapKeys(node.Children) {
		childNode := node.Children[childName]
		traverseCollectNodes(append(path, childName), childNode, nodesByPath)
	}
}
