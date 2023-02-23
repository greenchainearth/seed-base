package abft

import (
	"github.com/greenchainearth/seed-base/hash"
	"github.com/greenchainearth/seed-base/inter/dag"
)

// EventSource is a callback for getting events from an external storage.
type EventSource interface {
	HasEvent(hash.Event) bool
	GetEvent(hash.Event) dag.Event
}
