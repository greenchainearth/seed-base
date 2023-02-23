package tdag

import (
	"github.com/greenchainearth/seed-base/hash"
	"github.com/greenchainearth/seed-base/inter/dag"
)

type TestEvent struct {
	dag.MutableBaseEvent
	Name string
}

func (e *TestEvent) AddParent(id hash.Event) {
	parents := e.Parents()
	parents.Add(id)
	e.SetParents(parents)
}
