package router

import (
	"sort"

	"github.com/pkg/errors"
	"github.com/rwlist/gjrpc/internal/gen/protog"
)

func newTree(handlers []*handler, endpoints []endpoint) (*node, error) {
	root := newNode(nil)

	sort.Slice(handlers, func(i, j int) bool {
		return comparePaths(handlers[i].path, handlers[j].path)
	})
	for _, h := range handlers {
		nd := treeFindCreate(root, h.path)

		if nd.handler != nil {
			return nil, errors.Errorf("duplicate handler at path %s", protog.PathToString(h.path))
		}
		nd.handler = h
	}

	sort.Slice(endpoints, func(i, j int) bool {
		return comparePaths(endpoints[i].path, endpoints[j].path)
	})
	for _, e := range endpoints {
		nd := treeFindCreate(root, e.path)

		if nd.endpoint != nil {
			return nil, errors.Errorf("duplicate endpoint at path %s", protog.PathToString(e.path))
		}

		e := e
		nd.endpoint = &e
	}

	return root, nil
}

func comparePaths(a, b []string) bool {
	for i := 0; i < len(a) && i < len(b); i++ {
		if a[i] != b[i] {
			return a[i] < b[i]
		}
	}

	return len(a) < len(b)
}
