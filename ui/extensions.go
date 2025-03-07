// This file contains a few extensions to the "gomponents" library.
package ui

import (
	"maxwarden/security"

	. "maragu.dev/gomponents"
)

// Map a `map[T]U` to a [Group]
func MapMap[T comparable, U comparable](m map[T]U, cb func(U) Node) Group {
	var nodes []Node

	for k, _ := range m {
		nodes = append(nodes, cb(m[k]))
	}

	return nodes
}

// Map a `map[T]U` to a [Group]
// The callback provided must take the map key as an argument
func MapMapWithKey[T comparable, U comparable](m map[T]U, cb func(T, U) Node) Group {
	var nodes []Node

	for k, _ := range m {
		nodes = append(nodes, cb(k, m[k]))
	}

	return nodes
}

// Map a slice of anything to a [Group] (which is just a slice of [Node]-s).
// The callback must accept the index as an argument.
func MapWithIndex[T any](collection []T, callback func(index int, item T) Node) Group {
	nodes := make([]Node, 0, len(collection))
	for index, item := range collection {
		nodes = append(nodes, callback(index, item))
	}
	return nodes
}

func IfElse(condition bool, t Node, f Node) Node {
	if condition {
		return t
	} else {
		return f
	}
}

func IffElse(condition bool, t func() Node, f func() Node) Node {
	if condition {
		return t()
	} else {
		return f()
	}
}

// Sanitizes user input HTML
func SafeRaw(html string) Node {
	sanitized := security.SanitizationPolicy.Sanitize(html)
	return Raw(sanitized)
}

func CSSID(input string) string {
	return "#" + input
}

// For some reason this isn't included in the base distribution
func Template(children ...Node) Node {
	return El("template", children...)
}

func Open() Node {
	return Attr("open")
}