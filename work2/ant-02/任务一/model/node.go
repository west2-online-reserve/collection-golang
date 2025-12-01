package model

import "golang.org/x/net/html"

type Node struct {
	Type      html.NodeType
	Data      string
	ClassName string
}
