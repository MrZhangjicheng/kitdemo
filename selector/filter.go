package selector

import "context"

type NodeFilter func(context.Context, []Node) []Node
