package node

import (
	"fmt"
	"strings"
)

type Class[T any] struct {
	Id       string     `json:"id"`
	Data     *T         `json:"data"`
	Children []Class[T] `json:"children"`
	Parent   *Class[T]  `json:"parent"`
}

func New[T any](
	id string,
	data *T,
	parent *Class[T],
) *Class[T] {
	return &Class[T]{
		Id:       id,
		Data:     data,
		Children: []Class[T]{},
		Parent:   parent,
	}
}

func (n *Class[T]) AppendChild(
	child Class[T],
) {
	child.Parent = n
	n.Children = append(
		n.Children,
		child,
	)
}

func (n *Class[T]) Search(
	path string,
	separator string,
) (*Class[T], error) {
	if len(path) == 0 {
		return nil, ErrBadPath
	}

	lineage := strings.Split(path, separator)
	if len(lineage) == 1 {
		if lineage[0] == n.Id {
			return n, nil
		} else {
			return nil, ErrNotFound
		}
	} else {
		subpath := strings.Join(lineage[1:], separator)
		for i := range n.Children {
			childNode, err := n.Children[i].Search(subpath, separator)
			if err == nil {
				return childNode, nil
			}
		}
		return nil, ErrNotFound
	}
}

func (n *Class[T]) setData(
	data T,
	replaceIfExists bool,
) (*Class[T], error) {
	if n.Data != nil {
		if replaceIfExists {
			n.Data = &data
			return n, nil
		} else {
			return nil, ErrAlreadyExists
		}
	} else {
		n.Data = &data
		return n, nil
	}
}

func (n *Class[T]) Insert(
	path string,
	separator string,
	data T,
	replaceIfExists bool,
) (*Class[T], error) {
	lineage := strings.Split(path, separator)
	if len(lineage) == 0 {
		return nil, ErrBadPath
	}

	if len(lineage) == 1 {
		if lineage[0] == n.Id {
			return n.setData(data, replaceIfExists)
		} else {
			return nil, ErrBadPath
		}
	}

	subpath := strings.Join(lineage[1:], separator)
	var selectedChild *Class[T] = nil
	for i := range n.Children {
		if lineage[1] == n.Children[i].Id {
			selectedChild = &n.Children[i]
			break
		}
	}
	if selectedChild == nil {
		n.AppendChild(*New(lineage[1], nil, n))
		selectedChild = &n.Children[len(n.Children)-1]
	}

	return selectedChild.Insert(
		subpath,
		separator,
		data,
		replaceIfExists,
	)
}

func (n *Class[T]) GetChildById(id string) *Class[T] {
	for i := range n.Children {
		if n.Children[i].Id == id {
			return &n.Children[i]
		}
	}
	return nil
}

func (n *Class[T]) Locate(
	separator string,
) string {
	currentNode := n
	path := n.Id
	for currentNode.Parent != nil {
		path = fmt.Sprintf("%s%s%s",
			currentNode.Parent.Id,
			separator,
			path)
		currentNode = currentNode.Parent
	}
	return path
}
