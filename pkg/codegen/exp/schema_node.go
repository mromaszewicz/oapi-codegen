package exp

import (
	"strings"
)

type PathTreeNode struct {
	// PathElement is a segment of a path, eg, "schemas" in /components/schemas
	PathElement string
	// Children are the descendants of this node, keyed by name
	Children map[string]*PathTreeNode
	// Schema is a kin-openapi type, Schema or Parameter, etc.
	Schema any
	// Ref is set if this node is a reference to a schema.
	Ref string
	// The Parent node that has this one as a child.
	Parent *PathTreeNode
	// ElementName is used to generate type names when there are conflicts,
	// and corresponds to the location in the spec path, eg, something in "components"
	// would be a "Component".
	ElementName string
	// The FriendlyName is the shortest, easiest to use name for this node which
	// we can generate. It will be an alias to the full name, which is based
	// on the tree path. FriendlyNames can only be determined once all schemas
	// are loaded and collision tested.
	FriendlyName string
}

func (p *PathTreeNode) GetNodeByPathElements(path []string) *PathTreeNode {
	// The terminal case is that the path has one element.
	child, found := p.Children[path[0]]
	if !found {
		return nil
	}
	if len(path) == 1 {
		return child
	} else {
		return child.GetNodeByPathElements(path[1:])
	}
}

func (p *PathTreeNode) GetNodeByPath(path string) *PathTreeNode {
	// If there is a leading slash, remove it.
	if strings.HasPrefix(path, "/") {
		path = path[1:]
	}

	// If the string is empty when we get here, we are the node that they're
	// looking for.
	if path == "" {
		return p
	}

	// Now, we go through all the children and see if their string representation
	// matches the current path segment.
	for childName, child := range p.Children {
		if strings.HasPrefix(path, childName) {
			// we have a match.
			return child.GetNodeByPath(path[len(childName):])
		}
	}
	return nil
}

func (p *PathTreeNode) AddChild(child *PathTreeNode) {
	if p.Children == nil {
		p.Children = make(map[string]*PathTreeNode)
	}
	child.Parent = p
	p.Children[child.PathElement] = child
}
