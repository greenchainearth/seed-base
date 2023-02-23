package abft

import (
	"github.com/greenchainearth/seed-base/abft/dagidx"
	"github.com/greenchainearth/seed-base/abft/election"
	"github.com/greenchainearth/seed-base/hash"
	"github.com/greenchainearth/seed-base/inter/idx"
	"github.com/greenchainearth/seed-base/inter/pos"
)

type OrdererCallbacks struct {
	ApplyAtropos func(decidedFrame idx.Frame, atropos hash.Event) (sealEpoch *pos.Validators)

	EpochDBLoaded func(idx.Epoch)
}

type OrdererDagIndex interface {
	dagidx.ForklessCause
}

// Unlike processes events to reach finality on their order.
// Unlike abft.Seed, this raw level of abstraction doesn't track cheaters detection
type Orderer struct {
	config Config
	crit   func(error)
	store  *Store
	input  EventSource

	election *election.Election
	dagIndex OrdererDagIndex

	callback OrdererCallbacks
}

// New creates Orderer instance.
// Unlike Seed, Orderer doesn't updates DAG indexes for events, and doesn't detect cheaters
// It has only one purpose - reaching consensus on events order.
func NewOrderer(store *Store, input EventSource, dagIndex OrdererDagIndex, crit func(error), config Config) *Orderer {
	p := &Orderer{
		config:   config,
		store:    store,
		input:    input,
		crit:     crit,
		dagIndex: dagIndex,
	}

	return p
}
