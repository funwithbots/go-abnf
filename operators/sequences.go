package operators

// Concat defines a simple, ordered string of values (i.e., a concatenation of contiguous characters) by listing a
// sequence of rule names.
func Concat(key string, rules ...Operator) Operator {
	return func(s []byte) Alternatives {
		var concat func(l int, rules []Operator) Alternatives
		concat = func(l int, rules []Operator) Alternatives {
			if len(rules) == 0 {
				return Alternatives{
					{
						Key:   key,
						Value: s[:l],
					},
				}
			}

			var nodes Alternatives
			subNodes := rules[0](s[l:])
			for _, node := range subNodes {
				// add node as child of next nodes
				for _, n := range concat(l+len(node.Value), rules[1:]) {
					n.Children = append([]*Node{node}, n.Children...)
					nodes = append(nodes, n)
				}
			}
			if len(nodes) > 0 {
				nodes[0] = nodes.Best()
				for i := 1; i < len(nodes[1:]); i++ {
					nodes[i] = nil
				}
				nodes = nodes[:1]
			}
			return nodes
		}
		return concat(0, rules)
	}
}

// Alts defines a sequence of alternative elements that are separated by a forward slash ("/").
// Therefore, "foo / bar" will accept <foo> or <bar>.
func Alts(key string, rules ...Operator) Operator {
	return func(s []byte) Alternatives {
		var nodes Alternatives
		for _, rule := range rules {
			subNodes := rule(s)
			for _, node := range subNodes {
				nodes = append(nodes, &Node{
					Key:      key,
					Value:    node.Value,
					Children: Children{node},
				})
			}
		}
		return nodes
	}
}
